package engine

import (
	"errors"
)

// 将 git 远程仓库里的文件加入引擎索引
// 在开始搜索前务必调用 Finish 函数
func (engine *KunlunEngine) IndexGitRemote(gitURL string) error {
	if engine.walker == nil {
		return errors.New("engine 创建时未指定 walker 选项")
	}

	// 开始遍历
	engine.walkerWaitGroup.Add(1)
	engine.walker.WalkGitRemote(gitURL)

	return nil
}
