package parser

import (
	"bufio"
	"bytes"
	"strings"
	"time"
)

type HostsParser struct{}

func (h *HostsParser) Format() string { return "hosts" }

func (h *HostsParser) Parse(data []byte) ([]ParsedEntry, error) {
	s := bufio.NewScanner(bytes.NewReader(data))
	var out []ParsedEntry
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		ip := f[0]
		domain := strings.ToLower(strings.TrimSuffix(f[1], "."))
		if ip != "0.0.0.0" && ip != "127.0.0.1" && ip != "::1" {
			continue
		}
		out = append(out, ParsedEntry{
			Domain:  domain,
			Fetched: time.Now(),
		})
	}
	return out, s.Err()
}

func init() {
	Register(&HostsParser{})
}
