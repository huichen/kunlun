package types

import (
	"runtime"
)

// 搜索器创建参数
type SearcherOptions struct {
	// 每次请求可以发起多少线程用于做注解
	AnnotatorProcessors int
}

func NewSearcherOptions() *SearcherOptions {
	return &SearcherOptions{
		AnnotatorProcessors: runtime.NumCPU(),
	}
}

func (options *SearcherOptions) SetAnnotatorProcessors(num int) *SearcherOptions {
	if num > 0 {
		options.AnnotatorProcessors = num
	}
	return options
}
