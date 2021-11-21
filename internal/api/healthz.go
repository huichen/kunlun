package api

import (
	"io"
	"kunlun/internal/resource/engine"
	"net/http"
)

func Healthz(w http.ResponseWriter, req *http.Request) {
	if engine.GetEngine().IsFinished() {
		io.WriteString(w, "ok")
	} else {
		io.WriteString(w, "not_ok")
	}
}
