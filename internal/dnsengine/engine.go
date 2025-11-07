// SPDX-License-Identifier: GPL-3.0-or-later
package dnsengine

import (
	"github.com/lopster568/phantomDNS/internal/config"
	"github.com/lopster568/phantomDNS/internal/logger"
	"github.com/lopster568/phantomDNS/internal/policy"
	"github.com/lopster568/phantomDNS/internal/storage/models"
	"github.com/lopster568/phantomDNS/internal/storage/repositories"
	"github.com/lopster568/phantomDNS/internal/utils"
	"github.com/miekg/dns"
)

type BlocklistChecker interface {
	IsBlocked(domain string) (bool, error)
}

type Engine struct {
	upstreamManager *UpstreamManager
	repos           *repositories.Store
	policyEngine    *policy.Engine
	blocklist       BlocklistChecker
}

func (e *Engine) AttachBlocklistChecker(b BlocklistChecker) {
	e.blocklist = b
}

func NewDNSEngine(cfg config.DataPlaneConfig, repos *repositories.Store, pE *policy.Engine) (*Engine, error) {
	mgr, err := NewUpstreamManager(cfg.UpstreamResolvers, 4)
	if err != nil {
		return nil, err
	}
	return &Engine{
		upstreamManager: mgr,
		repos:           repos,
		policyEngine:    pE,
	}, nil
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
	go e.logQuery(r.Id, domain, w.RemoteAddr().String(), "block ("+reason+")")
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
	go e.logQuery(r.Id, domain, w.RemoteAddr().String(), "redirect")
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
	go e.logQuery(resp.Id, domain, w.RemoteAddr().String(), "allow")
}

func (e *Engine) logQuery(id uint16, domain, client, action string) {
	if e.repos == nil || e.repos.QueryLogs == nil {
		return
	}
	dnslog := &models.DNSQuery{
		ID:       uint(id),
		Domain:   domain,
		ClientIP: utils.AnonymizeIP(client),
		Action:   action,
	}
	if err := e.repos.QueryLogs.Save(dnslog); err != nil {
		logger.Log.Error("Failed to log DNS query: " + err.Error())
	}
}

// ProcessDNSQuery processes the DNS query and returns a response
func (e *Engine) ProcessDNSQuery(w dns.ResponseWriter, r *dns.Msg) {
	if r == nil || len(r.Question) == 0 {
		logger.Log.Warn("Received empty DNS query")
		return
	}
	domainName := r.Question[0].Name
	logger.Log.Infof("Received DNS query for %s", domainName)

	// --- Step 1: Check blocklist first ---
	if e.blocklist != nil {
		blocked, err := e.blocklist.IsBlocked(domainName)
		if err != nil {
			logger.Log.Error("Blocklist check failed: " + err.Error())
		} else if blocked {
			logger.Log.Infof("Blocked by blocklist: %s", domainName)
			e.respondBlocked(w, r, domainName, "blocklist")
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

	case policy.ActionRedirect:
		e.respondRedirect(w, r, domainName, decision.RedirectIP)

	default: // policy.ActionAllow
		e.forwardUpstream(w, r, domainName)
	}
}
