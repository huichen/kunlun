package engine

import (
	"sync"

	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/internal/searcher"
	"github.com/huichen/kunlun/internal/walker"
	"github.com/huichen/kunlun/pkg/log"
	"github.com/huichen/kunlun/pkg/types"
)

var (
	logger = log.GetLogger()
)

type KunlunEngine struct {
	options *types.EngineOptions

	searcher *searcher.Searcher

	indexer *indexer.Indexer

	walker *walker.IndexWalker

	walkerWaitGroup sync.WaitGroup

	finished bool
}

func NewKunlunEngine(options *types.EngineOptions) (*KunlunEngine, error) {
	if options == nil {
		options = types.NewEngineOptions()
	}

	engine := &KunlunEngine{}
	engine.options = options

	// 初始化搜索器
	engine.searcher = searcher.NewSearcher(options.SearcherOptions)

	// 初始化索引器
	engine.indexer = indexer.NewIndexer(options.IndexerOptions)

	// 初始化遍历器
	var err error
	engine.walker, err = walker.NewIndexWalker(options.WalkerOptions, false)
	if err != nil {
		return nil, err
	}

	// 监听 walker 返回并加入索引
	go engine.receiveWalkerOutput()

	return engine, nil
}

func (engine *KunlunEngine) receiveWalkerOutput() {
	fileChan := engine.walker.GetFileChan()
	for file := range fileChan {
		// 本批次索引完成
		if file.WalkingDone {
			engine.walkerWaitGroup.Done()
			continue
		}

		// 出错的话继续
		if file.Error != nil {
			continue
		}

		if !file.IsRepo {
			// 索引文件
			engine.indexer.IndexFile(file.Content,
				types.IndexFileInfo{
					Path:          file.AbsPath,
					Language:      file.Language,
					PathInRepo:    file.PathInRepo,
					RepoLocalPath: file.RepoLocalPath,
					RepoRemoteURL: file.RepoRemoteURL,
					CTagsEntries:  file.CTagsEntries,
				})
		} else {
			// 索引代码仓库
			path := file.RepoRemoteURL
			if path == "" {
				path = file.RepoLocalPath
			}
			log.GetLogger().Infof("索引仓库 %s", path)

			// 查找是否存在外部 ID
			var repoID uint64
			if id, ok := engine.options.RepoRemoteURLToIDMap[file.RepoRemoteURL]; ok {
				if id == 0 {
					log.GetLogger().Fatal("engine.options.RepoRemoteURLToIDMap 不能包含 0 ID")
				}
				repoID = id
			}

			engine.indexer.IndexRepo(
				types.IndexRepoInfo{
					RepoID:        repoID,
					RepoLocalPath: file.RepoLocalPath,
					RepoRemoteURL: file.RepoRemoteURL,
				})
		}
	}
}
