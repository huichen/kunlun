package rest

// Web 服务 API 基础字段
type ResponseBase struct {
	// Code == 0 代表无错误，否则返回 Message 错误信息
	Code int `json:"code"`

	// 错误信息
	Message string `json:"message"`
}
