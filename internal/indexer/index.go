package indexer

import (
	"errors"
	"fmt"

	"kunlun/pkg/log"
	"kunlun/pkg/types"
)

var (
	logger = log.GetLogger()
)

// 将 path 指定的文件添加到索引
// 如果该文件不存在于索引中，会给该文件分配一个自增的文档 ID
// 如果文件已经存在，则返回错误
//
// 该函数协程安全，请尽可能并发调用
func (indexer *Indexer) IndexFile(content []byte, info types.IndexFileInfo) error {
	if indexer.finished {
		return errors.New("indexer 已经完成索引")
	}

	path := info.Path
	lang := info.Language
	repoPath := info.RepoLocalPath
	repoRemoteURL := info.RepoRemoteURL
	pathInRepo := info.PathInRepo

	indexer.indexerLock.RLock()
	if _, ok := indexer.documentPathToIDMap[path]; ok {
		indexer.indexerLock.RUnlock()
		return fmt.Errorf("文件 %s 已经存在，请勿重复索引", path)
	}
	indexer.indexerLock.RUnlock()

	// 解析行起始位置
	lines := []uint32{}
	pre := byte('\n')
	for i, c := range content {
		if pre == '\n' {
			lines = append(lines, uint32(i))
		}
		pre = c
	}

	// 新生成一个 docID，并更新计数器等
	indexer.indexerLock.Lock()
	indexer.documentCounter++
	docID := indexer.documentCounter
	indexer.documentIDs = append(indexer.documentIDs, docID)
	indexer.documentPathToIDMap[path] = docID
	indexer.documentIDToPathMap[docID] = path
	indexer.documentIDToContentMap[docID] = &content
	indexer.totalContentSize += uint64(len(content))
	indexer.totalDocumentCount += 1

	// 添加语言
	var language *Language
	var ok bool
	if lang != "" {
		language, ok = indexer.langNameToIDMap[lang]
		if !ok {
			indexer.langCounter++
			language = &Language{
				ID:   indexer.langCounter,
				Name: lang,
			}
			indexer.langNameToIDMap[lang] = language
		}
	}

	// 试图找到代码仓库
	repo, ok := indexer.localPathToRepoMap[repoPath]
	if !ok {
		repo, ok = indexer.remoteURLToRepoMap[repoRemoteURL]
	}

	// 添加文件元信息
	indexer.documentIDToMetaMap[docID] = &DocumentMeta{
		DocumentID:         docID,
		LocalPath:          path,
		Size:               uint64(len(content)),
		LineStartLocations: lines,
		Language:           language,
		PathInRepo:         pathInRepo,
		Repo:               repo,
	}

	// 索引文件名
	if pathInRepo != "" {
		filename := []byte(pathInRepo)
		indexer.documentIDToFilenameMap[docID] = &filename
	}

	indexer.indexerLock.Unlock()

	// 索引内容
	indexer.contentIndexInfoChan <- IndexInfo{
		DocID:        docID,
		Content:      content,
		CTagsEntries: info.CTagsEntries,
	}

	return nil
}

func (indexer *Indexer) contentIndexWorker(shard int) {
	numDocsProcessed := 0
	for {
		info := <-indexer.contentIndexInfoChan
		if info.Exit {
			return
		}

		err := indexer.contentNgramIndices[shard].IndexDocument(info.DocID, info.Content, info.CTagsEntries)

		// 更新计数
		indexer.indexerLock.Lock()
		if err != nil {
			indexer.failedDocs++
			logger.Error(err)
		}
		indexer.documentIndexed++
		indexer.indexerLock.Unlock()

		numDocsProcessed++
		if indexer.maxDocsPerShard != 0 && numDocsProcessed > indexer.maxDocsPerShard {
			indexer.IncreaseShard()
			break
		}
	}
}
