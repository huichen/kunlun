package indexer

import (
	"strings"
	"time"
)

// 在开始检索前必须先调用这个函数
func (indexer *Indexer) Finish() {
	if indexer.finished {
		logger.Fatal("Finish 函数不能两次调用")
	}

	// 等待索引完成
	for indexer.documentCounter > indexer.documentIndexed {
		time.Sleep(time.Millisecond)
	}

	indexer.indexerLock.Lock()
	defer indexer.indexerLock.Unlock()

	if indexer.finished {
		logger.Fatal("Finish 函数不能两次调用")
	}

	// 尝试补充文件的 repo 信息
	// 如果正常遍历的话通常并不需要，只是为了防止遍历的时候会有遗漏
	for _, meta := range indexer.documentIDToMetaMap {
		// 如果添加过，不重复添加
		if meta.Repo != nil {
			continue
		}
		for _, repo := range indexer.localPathToRepoMap {
			if strings.HasPrefix(meta.LocalPath, repo.LocalPath+"/") {
				meta.Repo = repo
				meta.PathInRepo = strings.TrimPrefix(meta.LocalPath, repo.LocalPath+"/")
				break
			}
		}
	}

	// 退出 worker
	for i := 0; i < indexer.numIndexerShards; i++ {
		indexer.contentIndexInfoChan <- IndexInfo{Exit: true}
	}

	indexer.finished = true
}
