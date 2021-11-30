package indexer

import (
	"errors"

	"github.com/huichen/kunlun/internal/common_types"
	"github.com/huichen/kunlun/pkg/types"
)

type documentWithLines struct {
	DocumentID uint64

	Lines []uint32
}

// 从多个 token 中找到匹配的文档和对应行
// 这里有几个限定条件
// 1、这些 token 出现在同一行
// 2、同时遵从 caseSensitve 和 shouldDocBeRecalled 的要求
func (indexer *Indexer) searchMultiTokens(
	tokens []string,
	caseSensitive bool,
	shouldDocBeRecalled func(uint64) bool,
) ([]documentWithLines, error) {
	if len(tokens) == 0 {
		// 这和返回空数组（没有匹配的文档）不是一个概念，因此报错
		return nil, errors.New("tokens 不能为空")
	}

	// 首先得到所有 token 的搜索结果
	tokenSearchResults := []*[]common_types.DocumentWithSections{}
	var err error
	for _, token := range tokens {
		request := SearchTokenRequest{
			Token:         token,
			IsSymbol:      false, // 不在符号里搜索正则表达式
			CaseSensitive: caseSensitive,
			DocFilter:     shouldDocBeRecalled,
		}
		resp, searchErr := indexer.SearchToken(request)
		if searchErr == nil {
			tokenSearchResults = append(tokenSearchResults, &resp)
		} else {
			err = searchErr
		}
	}
	if len(tokenSearchResults) == 0 && err != nil {
		// 当所有的 token 都返回错误时，返回错误，上游会忽略多 token 的搜索结果
		return nil, err
	}

	return indexer.mergeResults(tokenSearchResults)
}

// 多个搜索结果取行交集
func (indexer *Indexer) mergeResults(results []*[]common_types.DocumentWithSections) ([]documentWithLines, error) {
	if len(results) == 0 {
		return nil, nil
	}

	pointers := make([]int, len(results))
	preDocIDs := make([]uint64, len(results))
	ret := []documentWithLines{}

	for {
		// 第一步：找到各个数组头部的最小值
		min := uint64(0)     // 保存最小值
		minAssigned := false // 最小值是否被赋值了（用于第一次给 min 赋值用）
		numMinValues := 0
		for i := 0; i < len(results); i++ {
			docs := results[i]
			if pointers[i] >= len(*docs) {
				goto breakPoint
			}

			v := (*docs)[pointers[i]].DocumentID
			if !minAssigned {
				// 第一次给 min 赋值
				min = v
				minAssigned = true
				numMinValues = 1
			} else if min > v {
				min = v
				numMinValues = 1
			} else if min == v {
				numMinValues++
			}

			// 合法性校验：头元素大于前一个元素
			if pointers[i] > 0 && preDocIDs[i] > v {
				return nil, errors.New("数组中元素不是严格递增的")
			}
		}

		if !minAssigned {
			// 所有数组已经穷尽，退出
			break
		}

		// 添加到返回结果，当且仅当头部全等于 min
		if numMinValues == len(results) {
			sectionsToMerge := []*[]types.Section{}
			for i := 0; i < len(results); i++ {
				docs := results[i]
				sectionsToMerge = append(sectionsToMerge, &(*docs)[pointers[i]].Sections)
			}

			// 额外校验，保证输出结果按照 docID 严格递增
			if len(ret) > 0 && ret[len(ret)-1].DocumentID >= min {
				return nil, errors.New("DocumentID 必须严格递增")
			}

			// 做行融合（并集）之后添加到输出
			lines, err := indexer.mergeLinesWithSections(min, sectionsToMerge)
			if err != nil {
				return nil, err
			}

			if len(ret) == 0 || ret[len(ret)-1].DocumentID != min {
				if len(lines) > 0 {
					ret = append(ret, documentWithLines{
						DocumentID: min,
						Lines:      lines,
					})
				}
			}
		}

		// 更新各个数组的头部指针
		for i := 0; i < len(results); i++ {
			docs := results[i]
			v := (*docs)[pointers[i]].DocumentID
			if min == v {
				preDocIDs[i] = v // 用于做严格递增校验
				pointers[i] = pointers[i] + 1
			}
		}
	}

breakPoint:

	return ret, nil
}

// 按行内交集
func (indexer *Indexer) mergeLinesWithSections(documentID uint64, sections []*[]types.Section) ([]uint32, error) {
	pointers := make([]int, len(sections))
	headLineNumbers := make([]uint32, len(sections))
	for i := 0; i < len(sections); i++ {
		v, _, _, err := indexer.getLine(documentID, (*sections[i])[0].Start)
		if err != nil {
			return nil, err
		}
		headLineNumbers[i] = v
	}

	ret := []uint32{}
	for {
		// 找到各个数组头部的最小值
		min := uint32(0)
		minAssigned := false
		for i := 0; i < len(sections); i++ {
			if pointers[i] >= len(*sections[i]) {
				break
			}
			v := headLineNumbers[i]

			if !minAssigned || min > v {
				min = v
				minAssigned = true
			}
		}

		if !minAssigned {
			// 所有数组已经穷尽，退出
			break
		}

		// 更新那些头等于 min 的数组的头指针
		numMinValues := 0
		for i := 0; i < len(sections); i++ {
			v := headLineNumbers[i]
			if min == v {
				pointers[i] = pointers[i] + 1
				if pointers[i] < len(*sections[i]) {
					v, _, _, err := indexer.getLine(documentID, (*sections[i])[pointers[i]].Start)
					if err != nil {
						return nil, err
					}
					headLineNumbers[i] = v
				}
				numMinValues++
			}
		}

		if numMinValues == len(sections) {
			if len(ret) == 0 || ret[len(ret)-1] != min {
				ret = append(ret, min)
			}
		}
	}

	return ret, nil
}
