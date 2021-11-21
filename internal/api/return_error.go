package api

import (
	"encoding/json"
	"io"
	"net/http"

	"kunlun/api/rest"
	"kunlun/pkg/log"
)

var (
	logger = log.GetLogger()
)

// 封装的错误返回函数，同时也打印一条日志
func ReturnError(w http.ResponseWriter, req *http.Request, err rest.ResponseBase) {
	logger.Infof("API错误，code = %d，msg = %s", err.Code, err.Message)
	resp, _ := json.Marshal(err)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(resp))
}
