package request

type Invoke struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
