// SPDX-License-Identifier: GPL-3.0-or-later
package dnsengine

import (
	"sync/atomic"
	"time"

	"github.com/lopster568/phantomDNS/internal/config"
	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/metrics"
	"github.com/lopster568/phantomDNS/internal/policy"
	"github.com/lopster568/phantomDNS/internal/storage/repositories"
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
	m.SetRcode(r, dns.RcodeRefused)
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
	if err := w.WriteMsg(resp); err != nil {
		logger.Log.Error("Failed to write DNS response: " + err.Error())
	}
}

// ProcessDNSQuery processes the DNS query and returns a response
func (e *Engine) ProcessDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	start := time.Now()
	success := false

	defer func() {
		elapsed := time.Since(start)
		e.metrics.Record(elapsed, success)
	}()

	if r == nil || len(r.Question) == 0 {
		logger.Log.Warn("Received empty DNS query")
		return
	}
	domainName := r.Question[0].Name
	// --- Step 1: Check blocklist first ---
	if e.blocklist != nil {
		blocked, err := e.blocklist.IsBlocked(domainName)
		if err != nil {
			logger.Log.Error("Blocklist check failed: " + err.Error())
		} else if blocked {
			logger.Log.Infof("Blocked by blocklist: %s", domainName)
			e.respondBlocked(w, r, domainName, "blocklist")
			success = true
			return
		}
	}

	// --- Step 2: Evaluate policy ---
	decision, err := e.policyEngine.Evaluate(domainName)
	if err != nil {
		logger.Log.Error("Failed to evaluate policy: " + err.Error())
		return
	}

	switch decision.Action {
	case policy.ActionDeny:
		logger.Log.Infof("Blocking via policy %s", decision.PolicyID)
		e.respondBlocked(w, r, domainName, decision.PolicyID)
		success = true

	case policy.ActionRedirect:
		e.respondRedirect(w, r, domainName, decision.RedirectIP)
		success = true

	default: // policy.ActionAllow
		e.forwardUpstream(w, r, domainName)
		success = true
	}
}
