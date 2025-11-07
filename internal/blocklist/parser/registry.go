package parser

type Parser interface {
	Parse([]byte) ([]ParsedEntry, error)
	Format() string
}

var registry = map[string]Parser{}

func Register(p Parser) {
	registry[p.Format()] = p
}

func Get(format string) (Parser, bool) {
	p, ok := registry[format]
	return p, ok
}
