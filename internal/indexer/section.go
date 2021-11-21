package indexer

import (
	"errors"

	"kunlun/pkg/types"
)

// 将多个 sections 做区间合并，比如
// [2,4) + [3,5) + [7,9) + [9,11] -> [2,5) + [7,11)
func mergeSections(sections []types.Section) ([]types.Section, error) {
	ret := []types.Section{}

	idx := 0
	for idx < len(sections) {
		start := sections[idx].Start
		end := sections[idx].End
		if start >= end {
			return nil, errors.New("Section 中 End 必须大于 Start")
		}

		preStart := start
		preEnd := end
		idx++
		for idx < len(sections) && sections[idx].Start <= end {
			newStart := sections[idx].Start
			newEnd := sections[idx].End

			// 合法性校验
			if newStart >= newEnd {
				return nil, errors.New("Section 中 End 必须大于 Start")
			}
			if newStart < preStart {
				return nil, errors.New("Section 数组中 Start 必须严格递增")
			}
			if newStart == preStart && newEnd < preEnd {
				return nil, errors.New("Section 数组中 End 必须严格递增")
			}
			end = max(end, newEnd)
			idx++

			preStart = newStart
			preEnd = newEnd
		}

		ret = append(ret, types.Section{
			Start: start,
			End:   end,
		})
	}

	return ret, nil
}
