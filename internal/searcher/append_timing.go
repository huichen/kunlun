package searcher

import (
	"time"

	"kunlun/pkg/types"
)

// 搜索结果中添加延时信息
func appendTimingInfo(context *Context, response *types.SearchResponse) {
	response.RecallDurationInMicroSeconds = context.recallEndTime.Sub(*context.searchStartTime).Microseconds()
	response.SearchDurationInMicroSeconds = time.Since(*context.searchStartTime).Microseconds()
}
