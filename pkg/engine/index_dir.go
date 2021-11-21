package engine

import (
	"errors"
)

// 将 dir 指定的目录里的文件加入引擎索引
// 在开始搜索前务必调用 Finish 函数
func (engine *KunlunEngine) IndexDir(dir string) error {
	if engine.walker == nil {
		return errors.New("engine 创建时未指定 walker 选项")
	}

	// 开始遍历
	engine.walkerWaitGroup.Add(1)
	engine.walker.WalkDir(dir)

	return nil
}
