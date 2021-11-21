package engine

import (
	"kunlun/internal/walker"
	"kunlun/pkg/log"
)

// 将 dir 下各级目录的 git 仓库 pull 更新
func (engine *KunlunEngine) PullRepos(dir string) error {
	repoPullWalker, err := walker.NewIndexWalker(engine.options.WalkerOptions, true)
	if err != nil {
		logger.Error(err)
		return err
	}

	go repoPullWalker.WalkDir(dir)

	fileChan := repoPullWalker.GetFileChan()

	for file := range fileChan {
		// 本批次索引完成
		if file.WalkingDone {
			break
		}

		// 出错的话继续
		if file.Error != nil {
			err = file.Error
		}

		if file.IsRepo {
			path := file.RepoRemoteURL
			if path == "" {
				path = file.RepoLocalPath
			}
			log.GetLogger().Infof("索引仓库 %s", path)
		}
	}

	if err != nil {
		logger.Error(err)
	}
	return err
}
