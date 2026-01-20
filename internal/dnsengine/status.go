package dnsengine

type Status struct {
	Running          bool
	AcceptingQueries bool
	PolicyEnabled    bool
	LastError        string
}

func (e *Engine) Status() Status {
	var lastErr string
	if v := e.state.lastError.Load(); v != nil {
		lastErr = v.(string)
	}

	return Status{
		Running:          true,
		AcceptingQueries: e.state.acceptQueries.Load(),
		// PolicyEnabled:    false,
		LastError: lastErr,
	}
}
