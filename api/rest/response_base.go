package rest

type ResponseBase struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
