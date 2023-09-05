package helper

type meta struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type response struct {
	Meta meta `json:"meta"`
	Data any  `json:"data"`
}

func ResponseAPI(success bool, code int, message string, data any) response {
	return response{Meta: meta{
		Success: success,
		Code:    code,
		Message: message,
	}, Data: data}
}
