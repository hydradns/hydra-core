package threat

import (
	"testing"
)

func TestDetector_NormalDomains(t *testing.T) {
	d := NewDetector()

	normal := []string{
		"google.com",
		"example.com",
		"stackoverflow.com",
		"en.wikipedia.org",
		"mail.google.com",
		"api.github.com",
	}

	for _, domain := range normal {
		r := d.Analyze(domain)
		if r.IsSuspicious {
			t.Errorf("normal domain %q flagged as suspicious: score=%.2f method=%s reason=%s",
				domain, r.ThreatScore, r.DetectionMethod, r.Reason)
		}
	}
}

func TestDetector_HexDGA(t *testing.T) {
	d := NewDetector()

	hex := []string{
		"a1b2c3d4e5f6a7b8.evil.com",
		"deadbeefcafebabe.malware.net",
	}

	for _, domain := range hex {
		r := d.Analyze(domain)
		if !r.IsSuspicious {
			t.Errorf("hex DGA domain %q not flagged", domain)
		}
		if r.DetectionMethod != "dga_hex" {
			t.Errorf("hex DGA domain %q: expected method dga_hex, got %s", domain, r.DetectionMethod)
		}
	}
}

func TestDetector_HighEntropy(t *testing.T) {
	d := NewDetector()

	// Random-looking domains
	suspicious := []string{
		"xkq7mz9plw2vb8nt.com",
		"r4nd0m5tr1ngd0m41n.net",
	}

	flagged := 0
	for _, domain := range suspicious {
		r := d.Analyze(domain)
		if r.IsSuspicious {
			flagged++
		}
	}
	if flagged == 0 {
		t.Error("no high-entropy domains were flagged")
	}
}

func TestDetector_LongDomain(t *testing.T) {
	d := NewDetector()

	long := "aaaaaaaaaaabbbbbbbbbbccccccccccddddddddddeeeeeeeeee.tunneling.com"
	r := d.Analyze(long)
	if !r.IsSuspicious {
		t.Errorf("long domain not flagged: %s", long)
	}
}

func TestDetector_DeepSubdomains(t *testing.T) {
	d := NewDetector()

	deep := "a.b.c.d.e.f.example.com"
	r := d.Analyze(deep)
	if !r.IsSuspicious {
		t.Error("deep subdomain not flagged")
	}
	if r.DetectionMethod != "subdomain_depth" {
		t.Errorf("expected subdomain_depth, got %s", r.DetectionMethod)
	}
}

func TestShannonEntropy(t *testing.T) {
	// "aaaa" has 0 entropy
	e := shannonEntropy("aaaa")
	if e != 0 {
		t.Errorf("expected 0 entropy for 'aaaa', got %.2f", e)
	}

	// "abcd" has 2.0 entropy (4 equally frequent chars)
	e = shannonEntropy("abcd")
	if e < 1.9 || e > 2.1 {
		t.Errorf("expected ~2.0 entropy for 'abcd', got %.2f", e)
	}
}
