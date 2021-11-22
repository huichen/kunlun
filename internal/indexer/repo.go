package indexer

import (
	"errors"
	"fmt"

	"github.com/huichen/kunlun/internal/common_types"
)

type CodeRepository struct {
	// 自增 ID
	ID uint64

	// 仓库在操作系统的路径，不包括结尾的 /
	LocalPath string

	// 仓库的远程路径
	RemoteURL string
}

// 将代码仓库添加到索引
func (indexer *Indexer) IndexRepo(info common_types.IndexRepoInfo) error {
	if indexer.finished {
		logger.Fatal("indexer 已经完成索引，请勿再添加")
	}

	indexer.indexerLock.Lock()
	defer indexer.indexerLock.Unlock()

	localPath := info.RepoLocalPath
	remoteURL := info.RepoRemoteURL

	if localPath == "" && remoteURL == "" {
		return errors.New("代码仓库 LocalPath 和 RemoteURL 不能都为空")
	}

	if _, ok := indexer.localPathToRepoMap[localPath]; ok {
		return fmt.Errorf("代码仓库 %s 已经存在，请勿重复添加", localPath)
	}

	// 更新 repo 计数
	indexer.repoCounter++

	// 当外部传入 repoID 时，使用外部 ID
	repoID := indexer.repoCounter
	if info.RepoID != 0 {
		if _, ok := indexer.idToRepoMap[info.RepoID]; ok {
			return fmt.Errorf("代码仓库 ID %d 已经存在，请勿重复添加", info.RepoID)
		}
		repoID = info.RepoID
		if indexer.repoCounter <= repoID {
			indexer.repoCounter = repoID + 1
		}
	}

	repo := &CodeRepository{
		ID:        repoID,
		LocalPath: localPath,
		RemoteURL: remoteURL,
	}

	if localPath != "" {
		indexer.localPathToRepoMap[localPath] = repo
	}
	if remoteURL != "" {
		indexer.remoteURLToRepoMap[remoteURL] = repo
	}
	indexer.idToRepoMap[indexer.repoCounter] = repo

	return nil
}

// 先返回 remote URL，如果为空则返回 local path
func (indexer *Indexer) GetRepoNameFromID(repoID uint64) string {
	indexer.indexerLock.RLock()
	defer indexer.indexerLock.RUnlock()

	if repo, ok := indexer.idToRepoMap[repoID]; ok {
		if repo.RemoteURL != "" {
			return repo.RemoteURL
		}
		return repo.LocalPath
	}

	return ""
}
