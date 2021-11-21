package indexer

import (
	"sync"

	"github.com/huichen/kunlun/internal/ngram_index"
	"github.com/huichen/kunlun/pkg/types"
)

type Indexer struct {
	// 读写锁保证线程安全
	indexerLock sync.RWMutex

	// 处于性能考虑，当 Finish 函数结束时我们将这个值设为 true
	// 此时只允许有读操作，从而避免使用锁
	finished bool

	// ngram -> 文档内容 的反向索引，分 shards
	contentNgramIndices  []*ngram_index.NgramIndex
	contentIndexInfoChan chan IndexInfo
	numIndexerShards     int
	maxDocsPerShard      int

	// 自增计数器，用于给文档标记 ID
	documentCounter uint64
	documentIDs     []uint64 // 自增数组
	documentIndexed uint64

	// 文档路径和 ID 之间的相互映射
	documentPathToIDMap map[string]uint64
	documentIDToPathMap map[uint64]string

	// 文档 ID 到内容和文件名原始字节的映射
	documentIDToContentMap  map[uint64]*[]byte
	documentIDToFilenameMap map[uint64]*[]byte

	// 文档 ID 到文档元信息的映射
	documentIDToMetaMap map[uint64]*DocumentMeta

	// 自增计数器，用于给代码仓库标记 ID
	repoCounter uint64

	// 存储仓库相关映射
	localPathToRepoMap map[string]*CodeRepository // 仓库在操作系统路径到仓库的映射
	remoteURLToRepoMap map[string]*CodeRepository // 仓库远程路径到仓库的映射
	idToRepoMap        map[uint64]*CodeRepository // 仓库 ID 到仓库的映射

	// 编程语言自增计数器和映射
	langCounter     uint64
	langNameToIDMap map[string]*Language

	// 统计信息
	totalContentSize   uint64
	totalDocumentCount uint64
	failedDocs         uint64
}

type IndexInfo struct {
	DocID        uint64
	Content      []byte
	CTagsEntries []*types.CTagsEntry

	// 退出信号
	Exit bool
}

func NewIndexer(options *types.IndexerOptions) *Indexer {
	if options == nil {
		options = types.NewIndexerOptions()
	}

	contentIndices := []*ngram_index.NgramIndex{}
	for i := 0; i < options.NumIndexerShards; i++ {
		contentIndices = append(contentIndices, ngram_index.NewNgramIndex())
	}

	indexer := Indexer{
		contentNgramIndices:  contentIndices,
		numIndexerShards:     options.NumIndexerShards,
		maxDocsPerShard:      options.MaxDocsPerShard,
		contentIndexInfoChan: make(chan IndexInfo, options.NumIndexerShards*2),

		documentCounter:         0,
		documentPathToIDMap:     make(map[string]uint64),
		documentIDToPathMap:     make(map[uint64]string),
		documentIDToMetaMap:     make(map[uint64]*DocumentMeta),
		documentIDToContentMap:  make(map[uint64]*[]byte),
		documentIDToFilenameMap: make(map[uint64]*[]byte),

		repoCounter:        0,
		localPathToRepoMap: make(map[string]*CodeRepository),
		remoteURLToRepoMap: make(map[string]*CodeRepository),
		idToRepoMap:        make(map[uint64]*CodeRepository),

		langCounter:     0,
		langNameToIDMap: make(map[string]*Language),
	}

	// 启动 content index worker
	for i := 0; i < options.NumIndexerShards; i++ {
		go indexer.contentIndexWorker(i)
	}

	return &indexer
}

func (indexer *Indexer) IncreaseShard() {
	indexer.indexerLock.Lock()
	shard := indexer.numIndexerShards
	indexer.numIndexerShards++
	indexer.indexerLock.Unlock()

	indexer.contentNgramIndices = append(indexer.contentNgramIndices, ngram_index.NewNgramIndex())

	go indexer.contentIndexWorker(shard)
}
