package engine

import (
	"sync"

	"kunlun/internal/indexer"
	"kunlun/internal/searcher"
	"kunlun/internal/walker"
	"kunlun/pkg/log"
	"kunlun/pkg/types"
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
			engine.indexer.IndexRepo(
				types.IndexRepoInfo{
					RepoLocalPath: file.RepoLocalPath,
					RepoRemoteURL: file.RepoRemoteURL,
				})
		}
	}
}
