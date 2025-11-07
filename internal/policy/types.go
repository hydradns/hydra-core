package policy

import "time"

type Action int

const (
	ActionAllow Action = iota
	ActionDeny
	ActionRedirect
)

func (a Action) String() string {
	switch a {
	case ActionAllow:
		return "allow"
	case ActionDeny:
		return "block"
	case ActionRedirect:
		return "redirect"
	default:
		return "unknown"
	}
}

type Policy struct {
	ID          string   `json:"id" gorm:"primaryKey;size:64"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Redirect    string   `json:"redirect_ip"` // optional redirect IP
	Source      string   `json:"source"`      // "local"|"external"|"system"
	Enabled     bool     `json:"enabled"`
	Priority    int      `json:"priority"`
	Action      string   `json:"action"`           // BLOCK|ALLOW|LOG_ONLY|REDIRECT
	Domains     []string `gorm:"-" json:"domains"` // stored in separate table or JSON
	Regexes     []string `gorm:"-" json:"regexes"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Decision is the result of evaluating a query.
type Decision struct {
	Action     Action
	PolicyID   string
	Category   string
	RedirectIP string
}
