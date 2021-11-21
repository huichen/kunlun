package searcher

import (
	"sync"

	"kunlun/internal/indexer"
)

// 文档过滤器
// 通过文档的编程语言、所属仓库和文件名等判断某个文档是否应该被召回
// 过滤器中保存了判断的状态，以避免同一种语言和仓库多次判断
// 判断时支持复杂的检索条件（Query 中读取）
type DocFilter struct {
	Indexer *indexer.Indexer

	Query *SearchQuery

	langFilter      map[uint64]bool
	repoFilter      map[uint64]bool
	fileFilter      map[uint64]bool
	docIDsWhitelist map[uint64]bool

	shouldRecallRepo func(uint64) bool

	lock sync.RWMutex
}

// 创建过滤器
// shouldRecallRepo 从外部传入一个判断仓库是否应该召回的方法，用于权限控制
// 注意：shouldRecallRepo 必须保证协程安全
func NewDocFilter(q *SearchQuery, idxr *indexer.Indexer, shouldRecallRepo func(uint64) bool, docIDs []uint64) *DocFilter {
	wl := make(map[uint64]bool)
	if len(docIDs) > 0 {
		for _, id := range docIDs {
			wl[id] = true
		}
	}
	if len(wl) == 0 {
		wl = nil
	}

	return &DocFilter{
		Indexer: idxr,

		Query: q,

		langFilter: make(map[uint64]bool),
		repoFilter: make(map[uint64]bool),
		fileFilter: make(map[uint64]bool),

		docIDsWhitelist: wl,

		shouldRecallRepo: shouldRecallRepo,
	}
}

func (df *DocFilter) ShouldRecallDocument(docID uint64) bool {
	if df.docIDsWhitelist != nil {
		if _, ok := df.docIDsWhitelist[docID]; !ok {
			return false
		}
	}

	// 获取文档的仓库名和语言
	meta := df.Indexer.GetMeta(docID)
	lang := ""
	langID := uint64(0)
	repoName := ""
	repoID := uint64(0)
	fileName := ""
	var hasLangCache, hasRepoCache, hasFileCache bool
	var match bool
	if meta != nil {
		if meta.Language != nil {
			lang = meta.Language.Name
			langID = meta.Language.ID
		}
		if meta.Repo != nil {
			if meta.Repo.RemoteURL != "" {
				repoName = meta.Repo.RemoteURL
			} else {
				repoName = meta.Repo.LocalPath
			}
			repoID = meta.Repo.ID
		}
		if meta.PathInRepo != "" {
			fileName = meta.PathInRepo
		} else {
			fileName = meta.LocalPath
		}
	}

	// 先过一遍缓存
	if df.Query.LanguageQuery != nil {
		df.lock.RLock()
		match, hasLangCache = df.langFilter[langID]
		df.lock.RUnlock()
		if hasLangCache && !match {
			return false
		}
	}
	if df.Query.FileQuery != nil {
		df.lock.RLock()
		match, hasFileCache = df.fileFilter[docID]
		df.lock.RUnlock()
		if hasFileCache && !match {
			return false
		}
	}
	if df.Query.RepoQuery != nil || df.shouldRecallRepo != nil {
		df.lock.RLock()
		match, hasRepoCache = df.repoFilter[repoID]
		df.lock.RUnlock()
		if hasRepoCache && !match {
			return false
		}
	}

	// 如果没有缓存但有 lang query，过一遍
	if !hasLangCache && df.Query.LanguageQuery != nil {
		match := matchRegexpQueries(lang, df.Query.LanguageQuery, &df.Query.LangRe, df.Query.LanguageNames)
		// 更新缓存
		df.lock.Lock()
		df.langFilter[langID] = match
		df.lock.Unlock()
		if !match {
			return false
		}
	}

	// 调用外部函数对 repo 做检测
	if !hasRepoCache && df.shouldRecallRepo != nil {
		match = df.shouldRecallRepo(repoID)
		if !match {
			df.lock.Lock()
			df.repoFilter[repoID] = false
			df.lock.Unlock()
			return false
		}
	}

	// 通过 query 对 repo 检测
	if !hasRepoCache && df.Query.RepoQuery != nil {
		match = matchRegexpQueries(repoName, df.Query.RepoQuery, &df.Query.RepoRe, df.Query.RepoNames)
		if !match {
			df.lock.Lock()
			df.repoFilter[repoID] = false
			df.lock.Unlock()
			return false
		}
	}

	// repo 的两个条件都通过了
	if !hasRepoCache && (df.Query.RepoQuery != nil || df.shouldRecallRepo != nil) {
		df.lock.Lock()
		df.repoFilter[repoID] = true
		df.lock.Unlock()
	}

	// 对文件名做检测
	if df.Query.FileQuery != nil {
		match := matchRegexpQueries(fileName, df.Query.FileQuery, &df.Query.FileRe, nil)
		df.lock.Lock()
		df.fileFilter[docID] = match
		df.lock.Unlock()
		if !match {
			return false
		}
	}

	return true
}

func (df *DocFilter) ShouldRecallRepo(repoID uint64) bool {
	// 首先获取文档的仓库名和语言
	var hasRepoCache bool
	var match bool

	if df.Query.RepoQuery != nil || df.shouldRecallRepo != nil {
		df.lock.RLock()
		match, hasRepoCache = df.repoFilter[repoID]
		df.lock.RUnlock()
		if hasRepoCache && !match {
			return false
		}
	}

	// 调用外部函数对 repo 做检测
	if !hasRepoCache && df.shouldRecallRepo != nil {
		match = df.shouldRecallRepo(repoID)
		logger.Info(match)
		if !match {
			df.lock.Lock()
			df.repoFilter[repoID] = false
			df.lock.Unlock()
			return false
		}
	}

	// 通过 query 对 repo 检测
	if !hasRepoCache && df.Query.RepoQuery != nil {
		repoName := df.Indexer.GetRepoNameFromID(repoID)
		match = matchRegexpQueries(repoName, df.Query.RepoQuery, &df.Query.RepoRe, df.Query.RepoNames)
		if !match {
			df.lock.Lock()
			df.repoFilter[repoID] = false
			df.lock.Unlock()
			return false
		}
	}

	// repo 的两个条件都通过了
	if !hasRepoCache && (df.Query.RepoQuery != nil || df.shouldRecallRepo != nil) {
		df.lock.Lock()
		df.repoFilter[repoID] = true
		df.lock.Unlock()
	}

	return true
}
