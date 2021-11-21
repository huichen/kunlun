package searcher

import (
	"errors"

	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/internal/query"
	"github.com/huichen/kunlun/pkg/types"
)

// 使用 OR 或者 AND 逻辑归并多个 query 的文档结果（保存在 context.query.QueryResults）
// 1、qs 中可以是 negate 或者 non negate queries，但不能全都是 negate queries
// 2、不支持 OR 和 negate queries 连用
// 3、query 的文档结果中的 DocumentID 和 Lines 必须都是严格递增的（不能有相同元素）
func mergeQueries(context *Context, qs []*query.Query, or bool) ([]types.DocumentWithSections, error) {
	if len(qs) == 0 {
		return nil, nil
	}

	pointers := make([]int, len(qs))
	preDocIDs := make([]uint64, len(qs))
	ret := []types.DocumentWithSections{}

	// 首先计算有多少个非 negate 的 query
	totalNonNegateQueries := 0
	totalNegateQueries := 0
	for _, q := range qs {
		if q == nil {
			return nil, errors.New("query 不能为 nil")
		}
		if context.query.QueryResults[q.ID] == nil {
			return nil, errors.New("query 的文档结果不能为 nil")
		}
		if !q.Negate {
			totalNonNegateQueries++
		} else {
			totalNegateQueries++
		}
	}

	// 合法性校验
	if totalNonNegateQueries == 0 {
		return nil, errors.New("不能全是 negate queries")
	}
	if or && totalNegateQueries != 0 {
		return nil, errors.New("OR 操作不能和 - 联用")
	}

	for {
		// 第一步：找到各个数组头部的最小值
		min := uint64(0)      // 保存最小值大小
		minHasNegate := false // 是否有 negate query 的头部等于最小值
		minAssigned := false  // 最小值是否被赋值了（用于第一次给 min 赋值用）
		numNonNegateMins := 0 // 等于最小值的非 negate query 个数
		for i := 0; i < len(qs); i++ {
			docs := context.query.QueryResults[qs[i].ID]
			if pointers[i] >= len(*docs) {
				// 跳过已经到达尽头的数组
				if !or && !qs[i].Negate {
					goto breakPoint
				}

				continue
			}

			v := (*docs)[pointers[i]].DocumentID
			if !minAssigned {
				// 第一次给 min 赋值
				min = v
				minAssigned = true
				minHasNegate = qs[i].Negate
				if !qs[i].Negate {
					numNonNegateMins = 1
				}
			} else if min >= v {
				if min == v {
					// 和已有的 min 相同，更新 negate 和非 negate query 计数
					if qs[i].Negate {
						minHasNegate = true
					} else {
						numNonNegateMins++
					}
				} else {
					// 需要更新 min 的情况
					minHasNegate = qs[i].Negate
					min = v
					if !qs[i].Negate {
						numNonNegateMins = 1
					} else {
						numNonNegateMins = 0
					}
				}
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

		// 添加到返回结果，当且仅当：
		// 1、OR 操作，或者
		// 2、AND 操作，且仅当命中全部 non negate queries
		if or || (!minHasNegate && numNonNegateMins == totalNonNegateQueries) {
			sectionsToMerge := []*[]types.Section{}
			for i := 0; i < len(qs); i++ {
				docs := context.query.QueryResults[qs[i].ID]
				if pointers[i] >= len(*docs) {
					// 跳过已经到达尽头的数组
					continue
				}
				if min == (*docs)[pointers[i]].DocumentID {
					sectionsToMerge = append(sectionsToMerge, &(*docs)[pointers[i]].Sections)
				}
			}

			// 额外校验，保证输出结果按照 docID 严格递增
			if len(ret) > 0 && ret[len(ret)-1].DocumentID >= min {
				return nil, errors.New("DocumentID 必须严格递增")
			}

			// 做行融合（并集）之后添加到输出
			sections, err := mergeSections(sectionsToMerge)
			if err != nil {
				return nil, err
			}
			ret = append(ret, types.DocumentWithSections{
				DocumentID: min,
				Sections:   sections,
			})

		}

		// 更新各个数组的头部指针
		for i := 0; i < len(qs); i++ {
			docs := context.query.QueryResults[qs[i].ID]
			if pointers[i] >= len(*docs) {
				// 跳过已经到达尽头的数组
				continue
			}

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

// 求并集
func mergeSections(lines []*[]types.Section) ([]types.Section, error) {
	if len(lines) == 0 {
		return nil, nil
	}

	pointers := make([]int, len(lines))
	preV := make([]uint32, len(lines))
	ret := []types.Section{}
	for {
		// 找到各个数组头部的最小值
		min := uint32(0)
		minAssigned := false
		for i := 0; i < len(lines); i++ {
			if pointers[i] >= len(*lines[i]) {
				// 跳过已经到达尽头的数组
				continue
			}
			v := (*lines[i])[pointers[i]].Start

			// 合法性校验
			end := (*lines[i])[pointers[i]].End
			if end <= v {
				return nil, errors.New("end 必须大于 start")
			}

			if !minAssigned || min > v {
				min = v
				minAssigned = true
			}

			// 合法性校验：头元素大于前一个元素
			if pointers[i] > 0 && preV[i] > v {
				return nil, errors.New("数组中元素不是严格递增的")
			}
		}

		if !minAssigned {
			// 所有数组已经穷尽，退出
			break
		}

		// 更新那些头等于 min 的数组的头指针
		minSections := []types.Section{}
		for i := 0; i < len(lines); i++ {
			if pointers[i] >= len(*lines[i]) {
				// 跳过已经到达尽头的数组
				continue
			}

			v := (*lines[i])[pointers[i]].Start
			if min == v {
				minSections = append(minSections, (*lines[i])[pointers[i]])
				pointers[i] = pointers[i] + 1
				preV[i] = v
			}
		}
		indexer.SortAndDedup(&minSections, func(i, j int) bool {
			if minSections[i].Start != minSections[j].Start {
				return minSections[i].Start < minSections[j].Start
			}
			return minSections[i].End < minSections[j].End
		})

		ret = append(ret, minSections...)
	}

	return ret, nil
}
