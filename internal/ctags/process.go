package ctags

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// 从命令行读取输出
func (cp *CTagsParser) read(rep *reply) error {
	// 读取返回
	if !cp.out.Scan() {
		err := cp.out.Err()
		cp.Close()
		return err
	}

	if bytes.Equal([]byte("(null)"), cp.out.Bytes()) {
		return nil
	}

	// 解析输出
	err := json.Unmarshal(cp.out.Bytes(), rep)
	if err != nil {
		return fmt.Errorf("JSON 反序列化失败")
	}
	return nil
}

// 发送文件内容到命令行
func (cp *CTagsParser) post(req *request, content []byte) error {
	// 发送请求头
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	body = append(body, '\n')
	if _, err = cp.in.Write(body); err != nil {
		return err
	}

	// 发送文件内容
	_, err = cp.in.Write(content)
	return err
}

type request struct {
	Command  string `json:"command"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

type reply struct {
	Typ     string `json:"_type"`
	Name    string `json:"name"`
	Version string `json:"version"`

	Command string `json:"command"`

	Path      string `json:"path"`
	Language  string `json:"language"`
	Line      int    `json:"line"`
	Kind      string `json:"kind"`
	End       int    `json:"end"`
	Scope     string `json:"scope"`
	ScopeKind string `json:"scopeKind"`
	Access    string `json:"access"`
	Signature string `json:"signature"`
}
