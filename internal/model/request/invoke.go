package request

type Invoke struct {
	Method string `json:"method" validate:"required,oneof=GET POST"`
	Path   string `json:"path" validate:"required"`
}
