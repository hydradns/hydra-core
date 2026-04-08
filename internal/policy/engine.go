package policy

import (
	"strings"
	"sync/atomic"
)

type Engine struct {
	snapshot atomic.Value // holds *Snapshot
}

func NewPolicyEngine() *Engine {
	e := &Engine{}
	e.snapshot.Store(buildSnapshot([]Policy{}))
	return e
}

// LoadPolicies replaces the full policy set atomically
func (e *Engine) LoadPolicies(policies []Policy) error {
	newSnap := buildSnapshot(policies)
	e.snapshot.Store(newSnap)
	return nil
}

// Evaluate returns the decision for given domain and context.
// Current implementation:
//   - normalize domain
//   - bloom negative => ALLOW
//   - exact match => return highest-priority policy decision
//   - otherwise => ALLOW (TODO: wildcard/regex)
func (e *Engine) Evaluate(domain string) (Decision, error) {
	d := normalizeDomain(domain)
	snap := e.snapshot.Load().(*PolicySnapshot)

	// Check exact match first, then walk up parent domains for subdomain matching.
	// e.g., www.godaddy.com → check www.godaddy.com, then godaddy.com
	parts := strings.Split(d, ".")
	for i := 0; i < len(parts)-1; i++ {
		candidate := strings.Join(parts[i:], ".")
		if snap.Bloom != nil && !snap.Bloom.TestString(candidate) {
			continue
		}
		if pols, ok := snap.Exact[candidate]; ok && len(pols) > 0 {
			best := pickHighestPriority(pols)
			return policyDecision(best), nil
		}
	}

	return Decision{Action: ActionAllow}, nil
}

// pickHighestPriority returns the matching policy with highest priority (deterministic).
func pickHighestPriority(candidates []*Policy) *Policy {
	var best *Policy
	for _, p := range candidates {
		if best == nil || p.Priority > best.Priority {
			best = p
		} else if p.Priority == best.Priority {
			// deterministic tie-breaker: lexicographic ID
			if p.ID < best.ID {
				best = p
			}
		}
	}
	return best
}

func policyDecision(p *Policy) Decision {
	if p == nil {
		return Decision{Action: ActionAllow}
	}
	var act Action
	switch strings.ToUpper(p.Action) {
	case "BLOCK":
		act = ActionDeny
	case "REDIRECT":
		act = ActionRedirect
	case "ALLOW":
		act = ActionAllow
	default:
		act = ActionAllow
	}
	return Decision{
		Action:     act,
		PolicyID:   p.ID,
		Category:   p.Category,
		RedirectIP: p.Redirect,
	}
}
