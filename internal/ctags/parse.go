package ctags

import (
	"github.com/huichen/kunlun/internal/common_types"
)

func (cp *CTagsParser) Parse(path string, content []byte) ([]*common_types.CTagsEntry, error) {
	// 发送请求
	req := request{
		Command:  "generate-tags",
		Size:     len(content),
		Filename: path,
	}
	if err := cp.post(&req, content); err != nil {
		return nil, err
	}

	// 循环读取解析的每条记录
	var es []*common_types.CTagsEntry
	for {
		var rep reply
		if err := cp.read(&rep); err != nil {
			return nil, err
		}
		if rep.Typ == "completed" {
			break
		}

		e := common_types.CTagsEntry{
			Sym:      rep.Name,
			Path:     rep.Path,
			Line:     rep.Line,
			Kind:     rep.Kind,
			Language: rep.Language,
		}

		es = append(es, &e)
	}

	return es, nil
}
