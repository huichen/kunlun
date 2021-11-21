package api

import (
	"io"
	"net/http"

	"github.com/huichen/kunlun/internal/resource/engine"
)

// 健康检查接口，用于返回搜索服务是否可用
// 引擎在索引构建阶段会返回 not_ok
func Healthz(w http.ResponseWriter, req *http.Request) {
	if engine.GetEngine().IsFinished() {
		io.WriteString(w, "ok")
	} else {
		io.WriteString(w, "not_ok")
	}
}
