// SPDX-License-Identifier: GPL-3.0-or-later
package dnsengine

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/lopster568/phantomDNS/internal/config"
	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/metrics"
	"github.com/lopster568/phantomDNS/internal/policy"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"github.com/lopster568/phantomDNS/internal/storage/repositories"
	"github.com/lopster568/phantomDNS/internal/threat"
	"github.com/miekg/dns"
)

type BlocklistChecker interface {
	IsBlocked(domain string) (bool, error)
}

type RuntimeState struct {
	acceptQueries atomic.Bool
	// policyEnabled atomic.Bool
	lastError atomic.Value
}

type Engine struct {
	upstreamManager *UpstreamManager
	policyEngine    *policy.Engine
	blocklist       BlocklistChecker
	state           *RuntimeState
	metrics         *metrics.QueryMetrics
	queryLog        repositories.QueryLogRepository
	statistics      repositories.StatisticsRepository
	threatDetector  *threat.Detector
}

func (e *Engine) AttachBlocklistChecker(b BlocklistChecker) {
	e.blocklist = b
}

func NewDNSEngine(cfg config.DataPlaneConfig, repos *repositories.Store, pE *policy.Engine) (*Engine, error) {
	mgr, err := NewUpstreamManager(cfg.UpstreamResolvers, 4)
	state := &RuntimeState{}
	state.acceptQueries.Store(false)
	qm := metrics.NewQueryMetrics()

	if err != nil {
		return nil, err
	}
	return &Engine{
		upstreamManager: mgr,
		policyEngine:    pE,
		state:           state,
		metrics:         qm,
		queryLog:        repos.QueryLogs,
		statistics:      repos.Statistics,
		threatDetector:  threat.NewDetector(),
	}, nil
}

func (e *Engine) SetAcceptQueries(enabled bool) {
	e.state.acceptQueries.Store(enabled)
}

// Cleanup the resources used by the Engine
func (e *Engine) Shutdown() {
	if e.upstreamManager != nil {
		e.upstreamManager.Close()
	}
}

func (e *Engine) respondBlocked(w dns.ResponseWriter, r *dns.Msg, domain, reason string) {
	m := new(dns.Msg)
	m.SetReply(r)
	// Return 0.0.0.0 instead of REFUSED — browsers treat REFUSED as "try another DNS"
	// but 0.0.0.0 causes an immediate connection failure (ERR_CONNECTION_REFUSED)
	rr, err := dns.NewRR(r.Question[0].Name + " 60 IN A 0.0.0.0")
	if err == nil {
		m.Answer = append(m.Answer, rr)
	}
	if err := w.WriteMsg(m); err != nil {
		logger.Log.Error("Failed to write DNS block response: " + err.Error())
	}
}

func (e *Engine) respondRedirect(w dns.ResponseWriter, r *dns.Msg, domain, ip string) {
	m := new(dns.Msg)
	m.SetReply(r)
	rr, err := dns.NewRR(domain + " 60 IN A " + ip)
	if err != nil {
		logger.Log.Error("Failed to create redirect RR: " + err.Error())
		m.SetRcode(r, dns.RcodeServerFailure)
		_ = w.WriteMsg(m)
		return
	}
	m.Answer = append(m.Answer, rr)
	if err := w.WriteMsg(m); err != nil {
		logger.Log.Error("Failed to write DNS redirect response: " + err.Error())
	}
}

func (e *Engine) forwardUpstream(w dns.ResponseWriter, r *dns.Msg, domain string) {
	resp, err := e.upstreamManager.Exchange(r, 5, 2)
	if err != nil {
		logger.Log.Error("Upstream query failed: " + err.Error())
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
		_ = w.WriteMsg(m)
		return
	}
	if resp == nil {
		logger.Log.Error("Upstream returned nil response for: " + domain)
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
		_ = w.WriteMsg(m)
		return
	}
	if err := w.WriteMsg(resp); err != nil {
		logger.Log.Error("Failed to write DNS response: " + err.Error())
	}
}

// normalizeDomain lowercases and strips the trailing dot from a DNS FQDN.
func normalizeDomain(d string) string {
	return strings.TrimSuffix(strings.ToLower(d), ".")
}

// ProcessDNSQuery processes the DNS query and returns a response
func (e *Engine) ProcessDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	if r == nil || len(r.Question) == 0 {
		return
	}

	if !e.state.acceptQueries.Load() {
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeRefused)
		_ = w.WriteMsg(m)
		return
	}

	start := time.Now()
	success := false

	defer func() {
		elapsed := time.Since(start)
		e.metrics.Record(elapsed, success)
	}()

	domainName := normalizeDomain(r.Question[0].Name)
	clientIP := ""
	if w.RemoteAddr() != nil {
		clientIP = w.RemoteAddr().String()
	}

	// Run threat detection on every query (non-blocking, just scoring)
	var threatResult threat.Result
	if e.threatDetector != nil {
		threatResult = e.threatDetector.Analyze(domainName)
	}

	// --- Step 1: Check blocklist first ---
	if e.blocklist != nil {
		blocked, err := e.blocklist.IsBlocked(domainName)
		if err != nil {
			logger.Log.Error("Blocklist check failed: " + err.Error())
		} else if blocked {
			logger.Log.Infof("Blocked by blocklist: %s", domainName)
			e.logQuery(domainName, clientIP, "block", threatResult)
			e.respondBlocked(w, r, domainName, "blocklist")
			success = true
			return
		}
	}

	// --- Step 2: Evaluate policy ---
	decision, err := e.policyEngine.Evaluate(domainName)
	if err != nil {
		logger.Log.Error("Failed to evaluate policy: " + err.Error())
		e.logQuery(domainName, clientIP, "error", threatResult)
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
		_ = w.WriteMsg(m)
		return
	}

	switch decision.Action {
	case policy.ActionDeny:
		logger.Log.Infof("Blocking via policy %s", decision.PolicyID)
		e.logQuery(domainName, clientIP, "block", threatResult)
		e.respondBlocked(w, r, domainName, decision.PolicyID)
		success = true

	case policy.ActionRedirect:
		e.logQuery(domainName, clientIP, "redirect", threatResult)
		e.respondRedirect(w, r, domainName, decision.RedirectIP)
		success = true

	default: // policy.ActionAllow
		action := "allow"
		if threatResult.IsSuspicious {
			action = "flagged"
			logger.Log.Warnf("Suspicious domain allowed: %s (score=%.2f, method=%s)", domainName, threatResult.ThreatScore, threatResult.DetectionMethod)
		}
		e.logQuery(domainName, clientIP, action, threatResult)
		e.forwardUpstream(w, r, domainName)
		success = true
	}
}

func (e *Engine) logQuery(domain, clientIP, action string, tr threat.Result) {
	if e.queryLog == nil {
		return
	}
	q := &models.DNSQuery{
		Domain:          domain,
		ClientIP:        clientIP,
		Action:          action,
		IsSuspicious:    tr.IsSuspicious,
		ThreatScore:     tr.ThreatScore,
		DetectionMethod: tr.DetectionMethod,
		ThreatReason:    tr.Reason,
	}
	// Map "flagged" to "allow" for statistics (flagged domains are still forwarded)
	statsAction := action
	if statsAction == "flagged" {
		statsAction = "allow"
	}

	go func() {
		if err := e.queryLog.Save(q); err != nil {
			logger.Log.Errorf("Failed to log query: %v", err)
		}
		if e.statistics != nil {
			if err := e.statistics.IncrementCounter(statsAction); err != nil {
				logger.Log.Errorf("Failed to increment stats: %v", err)
			}
		}
	}()
}

func (e *Engine) Metrics() *metrics.QueryMetrics {
	return e.metrics
}
