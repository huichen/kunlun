package engine

import (
	"time"
)

func (engine *KunlunEngine) Finish() {
	// 等待所有遍历请求完成
	engine.walkerWaitGroup.Wait()

	// 让子弹飞一会
	time.Sleep(time.Millisecond * 10)

	// 先终止索引器，再终止遍历器
	engine.indexer.Finish()
	engine.walker.Finish()

	engine.finished = true
}

func (engine *KunlunEngine) IsFinished() bool {
	return engine.finished
}
