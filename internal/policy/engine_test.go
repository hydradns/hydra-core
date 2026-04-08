package policy

import (
	"testing"
)

func TestNormalizeDomain(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"EXAMPLE.COM", "example.com"},
		{"example.com.", "example.com"},
		{"  Example.COM.  ", "example.com"},
		{"already.normal", "already.normal"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := normalizeDomain(tt.input); got != tt.want {
			t.Errorf("normalizeDomain(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestEvaluate_AllowByDefault(t *testing.T) {
	e := NewPolicyEngine()
	d, err := e.Evaluate("unknown.com")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionAllow {
		t.Errorf("expected ActionAllow, got %v", d.Action)
	}
}

func TestEvaluate_BlockExactDomain(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "block-ads", Action: "BLOCK", Priority: 100, Domains: []string{"ads.example.com"}},
	})

	d, err := e.Evaluate("ads.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionDeny {
		t.Errorf("expected ActionDeny, got %v", d.Action)
	}
	if d.PolicyID != "block-ads" {
		t.Errorf("expected policy ID block-ads, got %q", d.PolicyID)
	}
}

func TestEvaluate_CaseInsensitive(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "p1", Action: "BLOCK", Priority: 100, Domains: []string{"ADS.Example.COM"}},
	})

	d, err := e.Evaluate("ads.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionDeny {
		t.Errorf("expected ActionDeny for case-insensitive match")
	}
}

func TestEvaluate_TrailingDotNormalized(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "p1", Action: "BLOCK", Priority: 100, Domains: []string{"blocked.com"}},
	})

	d, err := e.Evaluate("blocked.com.")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionDeny {
		t.Errorf("expected ActionDeny for domain with trailing dot")
	}
}

func TestEvaluate_Redirect(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "redir", Action: "REDIRECT", Priority: 100, Redirect: "1.2.3.4", Domains: []string{"redirect.me"}},
	})

	d, err := e.Evaluate("redirect.me")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionRedirect {
		t.Errorf("expected ActionRedirect, got %v", d.Action)
	}
	if d.RedirectIP != "1.2.3.4" {
		t.Errorf("expected redirect IP 1.2.3.4, got %q", d.RedirectIP)
	}
}

func TestEvaluate_HighestPriorityWins(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "low", Action: "ALLOW", Priority: 10, Domains: []string{"test.com"}},
		{ID: "high", Action: "BLOCK", Priority: 200, Domains: []string{"test.com"}},
	})

	d, err := e.Evaluate("test.com")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionDeny {
		t.Errorf("expected highest priority (BLOCK) to win, got %v", d.Action)
	}
	if d.PolicyID != "high" {
		t.Errorf("expected policy ID 'high', got %q", d.PolicyID)
	}
}

func TestEvaluate_TieBreakByID(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "zzz", Action: "ALLOW", Priority: 100, Domains: []string{"tie.com"}},
		{ID: "aaa", Action: "BLOCK", Priority: 100, Domains: []string{"tie.com"}},
	})

	d, err := e.Evaluate("tie.com")
	if err != nil {
		t.Fatal(err)
	}
	if d.PolicyID != "aaa" {
		t.Errorf("expected lexicographically first ID 'aaa' on tie, got %q", d.PolicyID)
	}
}

func TestEvaluate_BloomNegativeShortCircuit(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "p1", Action: "BLOCK", Priority: 100, Domains: []string{"blocked.com"}},
	})

	// Domain not in bloom filter should be allowed without hitting exact map
	d, err := e.Evaluate("notblocked.com")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionAllow {
		t.Errorf("expected ActionAllow for domain not in bloom filter")
	}
}

func TestEvaluate_MultipleDomainsSamePolicy(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{
		{ID: "multi", Action: "BLOCK", Priority: 100, Domains: []string{"a.com", "b.com", "c.com"}},
	})

	for _, domain := range []string{"a.com", "b.com", "c.com"} {
		d, err := e.Evaluate(domain)
		if err != nil {
			t.Fatal(err)
		}
		if d.Action != ActionDeny {
			t.Errorf("expected %s to be blocked", domain)
		}
	}
}

func TestEvaluate_EmptyPolicies(t *testing.T) {
	e := NewPolicyEngine()
	e.LoadPolicies([]Policy{})

	d, err := e.Evaluate("anything.com")
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionAllow {
		t.Errorf("expected ActionAllow with empty policies")
	}
}

func TestPolicyDecision_UnknownAction(t *testing.T) {
	d := policyDecision(&Policy{ID: "x", Action: "UNKNOWN_ACTION"})
	if d.Action != ActionAllow {
		t.Errorf("expected ActionAllow for unknown action, got %v", d.Action)
	}
}

func TestPolicyDecision_Nil(t *testing.T) {
	d := policyDecision(nil)
	if d.Action != ActionAllow {
		t.Errorf("expected ActionAllow for nil policy, got %v", d.Action)
	}
}

func TestActionString(t *testing.T) {
	tests := []struct {
		a    Action
		want string
	}{
		{ActionAllow, "allow"},
		{ActionDeny, "block"},
		{ActionRedirect, "redirect"},
		{Action(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.a.String(); got != tt.want {
			t.Errorf("Action(%d).String() = %q, want %q", tt.a, got, tt.want)
		}
	}
}
