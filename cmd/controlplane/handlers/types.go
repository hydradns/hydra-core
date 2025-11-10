package handlers

// ResponseGeneric represents a generic API response
type ResponseGeneric struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  *string     `json:"error"`
}
