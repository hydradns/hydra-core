package parser

// import (
// 	"regexp"
// 	"strings"
// )

// type Parser interface {
// 	Parse(data []byte) ([]string, error)
// 	Format() string
// }

// type PlainTextParser struct{}

// type HostsParser struct {
// 	lineRegex *regexp.Regexp
// }

// func NewHostsParser() *HostsParser {
// 	return &HostsParser{
// 		lineRegex: regexp.MustCompile(`(?i)^(?:0\.0\.0\.0|127\.0\.0\.1)?\s*([a-z0-9.-]+)$`),
// 	}
// }

// func (p *PlainTextParser) Parse(data []byte) ([]string, error) {
// 	lines := strings.Split(string(data), "\n")
// 	var domains []string
// 	for _, line := range lines {
// 		domain := normalizeDomain(line)
// 		if domain != "" {
// 			domains = append(domains, domain)
// 		}
// 	}
// 	return domains, nil
// }

// func (p *PlainTextParser) Format() string {
// 	return "plaintext"
// }

// // normalizeDomain trims spaces, removes trailing dots, and lowercases
// func normalizeDomain(d string) string {
// 	d = strings.TrimSpace(d)
// 	d = strings.TrimSuffix(d, ".")
// 	d = strings.ToLower(d)
// 	if d == "" || strings.ContainsAny(d, " /\\") {
// 		return ""
// 	}
// 	return d
// }
