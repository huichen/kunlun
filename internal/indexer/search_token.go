package indexer

import "github.com/huichen/kunlun/internal/common_types"

const (
	Atoa = 32
)

// 通过单串（token）搜索文档
// 单串指的是非正则表达式的简单字符串，可以是文本中单串也可以是符号（变量、方法名等）
type SearchTokenRequest struct {
	// 搜索关键词
	Token string

	// 是否只检查符号
	IsSymbol bool

	// 是否区分大小写
	CaseSensitive bool

	// 文档过滤器
	DocFilter func(docID uint64) bool
}

// 在索引中查找包含关键词的文档
// 如果没有返回结果，返回空数组（非 nil）
func (indexer *Indexer) SearchToken(request SearchTokenRequest) ([]common_types.DocumentWithSections, error) {
	matchedDocs, err := indexer.internalSearchToken(request)
	if err != nil {
		return nil, err
	}

	return matchedDocs, nil
}

func (indexer *Indexer) internalSearchToken(request SearchTokenRequest) ([]common_types.DocumentWithSections, error) {
	// 获取 keyword 对应的 index keys，distance 等信息
	offset, key1, key2, distance := indexer.getTwoKeysFromToken(request.Token)
	if key1 == 0 {
		return nil, nil
	}

	// 搜索文档内容
	keyword := []byte(request.Token)
	matchedDocs, err := indexer.searchContent(
		keyword, offset, key1, key2, distance, request.CaseSensitive, request.IsSymbol, request.DocFilter)

	return matchedDocs, err
}
