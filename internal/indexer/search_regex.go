package indexer

import (
	"errors"
	"regexp"
	"sort"

	"kunlun/pkg/types"
)

// 搜索正则表达式
type SearchRegexRequest struct {
	// 搜索关键词
	Regex string

	// 为 true 时只选择不匹配的文档
	Negate bool

	// 是否区分大小写
	CaseSensitive bool

	// 表达式的字串
	Tokens []string

	// 每个文件最多返回多少结果，设为 0 或者 -1 全部返回
	MaxResultsPerFile int

	// 是否只检查符号
	IsSymbol bool

	// 当不为空时，仅从下面的文件列表中搜索
	// 如果 CandidateDocsNegate == false 则用作黑名单
	CandidateDocs       *[]uint64
	CandidateDocsNegate bool

	DocFilter func(docID uint64) bool
}

type SearchRegexResponse struct {
	Documents []types.DocumentWithSections
	Negate    bool

	RegexSearchTimes int
}

// 在索引中查找包含关键词的文档
func (indexer *Indexer) SearchRegex(request SearchRegexRequest) (*SearchRegexResponse, error) {
	if request.CandidateDocsNegate && request.Negate {
		return nil, errors.New("negate 正则表达式不能和 docsID 黑名单共存")
	}

	// 仅当没有白名单，且正则表达式为 negate 时，返回 negate
	negate := false
	if (request.CandidateDocs == nil || len(*request.CandidateDocs) == 0) && request.Negate {
		negate = true
		request.Negate = false
	}

	matchedDocs, times, err := indexer.internalSearchRegex(request)
	if err != nil {
		return nil, err
	}
	return &SearchRegexResponse{
		Documents:        matchedDocs,
		Negate:           negate,
		RegexSearchTimes: times,
	}, nil
}

func (indexer *Indexer) internalSearchRegex(request SearchRegexRequest) ([]types.DocumentWithSections, int, error) {
	maxResults := request.MaxResultsPerFile
	if maxResults <= 0 {
		maxResults = -1
	}

	// 得到 docID 黑名单
	excludeDocIDs := make(map[uint64]bool)
	if request.CandidateDocsNegate && request.CandidateDocs != nil && len(*request.CandidateDocs) != 0 {
		for _, doc := range *request.CandidateDocs {
			excludeDocIDs[doc] = true
		}
	}

	if len(request.Tokens) == 0 {
		// 没有候选 token 的情况，尽可能缩小文档搜索范围
		var includeDocIDs *[]uint64
		if !request.CandidateDocsNegate && request.CandidateDocs != nil && len(*request.CandidateDocs) != 0 {
			includeDocIDs = request.CandidateDocs
		} else {
			includeDocIDs = &indexer.documentIDs
		}

		// 对候选文档做正则表达式搜索，返回
		return indexer.searchRegexInDocs(
			includeDocIDs, nil, excludeDocIDs,
			request.Regex, request.Negate, request.CaseSensitive,
			maxResults, request.DocFilter)
	}

	// 通过检索局部关键词，得到一个较小的候选集
	docIDs := []uint64{}
	runSearch := false
	var matchedDocs []DocumentWithLines
	if request.CandidateDocs == nil || !request.Negate {
		var err error
		matchedDocs, err = indexer.searchMultiTokens(
			request.Tokens, request.CaseSensitive, request.DocFilter)
		if err == nil {
			runSearch = true
		} else {
			// 如果出错的话就不用
			matchedDocs = nil
		}
	}
	if request.CandidateDocs != nil && len(*request.CandidateDocs) != 0 && !request.CandidateDocsNegate {
		if runSearch {
			// 如果做了检索，则取交集
			matchedDocs = andMerge(*request.CandidateDocs, matchedDocs)
		} else {
			// 否则返回白名单
			docIDs = *request.CandidateDocs
		}
	}

	// 前面已经使用 docFilterFunc 做了过滤，这里就不需要了
	return indexer.searchRegexInDocs(
		&docIDs, matchedDocs, excludeDocIDs,
		request.Regex, request.Negate, request.CaseSensitive,
		maxResults, nil)
}

// 逐个文档做正则表达式匹配
// 这个操作很昂贵，请尽可能缩小 docIDs 的范围
func (indexer *Indexer) searchRegexInDocs(
	includeDocIDs *[]uint64,
	matchedDocs []DocumentWithLines,
	excludeDocIDs map[uint64]bool,
	regex string,
	negate bool,
	caseSensitive bool,
	maxResultsPerFile int,
	shouldDocBeRecalled func(uint64) bool,
) ([]types.DocumentWithSections, int, error) {
	// 是否大小写敏感
	if !caseSensitive {
		regex = "(?i)" + regex
	}
	re := regexp.MustCompile(regex)

	// 启动 worker
	numRegexSearchWorkers := indexer.numIndexerShards
	infoChan := make(chan regexSearchInfo, numRegexSearchWorkers*2)
	returnChan := make(chan regexSearchReturn, numRegexSearchWorkers*2)
	for i := 0; i < numRegexSearchWorkers; i++ {
		go indexer.regexSearcher(infoChan, returnChan, re, maxResultsPerFile, negate)
	}

	// 发送任务到 workers
	// 使用 docIDs
	if matchedDocs == nil {
		for _, docID := range *includeDocIDs {
			// 如果在排除名单里，不召回
			if _, ok := excludeDocIDs[docID]; ok {
				continue
			}

			// 然后用外部函数判断是否应该被召回
			if shouldDocBeRecalled != nil && !shouldDocBeRecalled(docID) {
				continue
			}

			infoChan <- regexSearchInfo{DocumentID: docID}
		}
	} else {
		for _, doc := range matchedDocs {
			// 如果在排除名单里，不召回
			if _, ok := excludeDocIDs[doc.DocumentID]; ok {
				continue
			}

			// 然后用外部函数判断是否应该被召回
			if shouldDocBeRecalled != nil && !shouldDocBeRecalled(doc.DocumentID) {
				continue
			}

			infoChan <- regexSearchInfo{
				DocumentID: doc.DocumentID,
				Lines:      doc.Lines,
			}
		}
	}

	// 发送终止信号
	for i := 0; i < numRegexSearchWorkers; i++ {
		infoChan <- regexSearchInfo{Exit: true}
	}

	// 收集结果
	searchTimes := 0
	retDocs := []types.DocumentWithSections{}
	for i := 0; i < numRegexSearchWorkers; i++ {
		ret := <-returnChan
		retDocs = append(retDocs, ret.Docs...)
		searchTimes += ret.SearchedTimes
	}

	// 排序
	sort.Slice(retDocs, func(i, j int) bool {
		return retDocs[i].DocumentID < retDocs[j].DocumentID
	})

	return retDocs, searchTimes, nil
}

type regexSearchInfo struct {
	DocumentID uint64
	Lines      []uint32
	Exit       bool
}
type regexSearchReturn struct {
	Docs          []types.DocumentWithSections
	SearchedTimes int
}

// 正则表达式匹配工作线程
func (indexer *Indexer) regexSearcher(
	infoChan chan regexSearchInfo,
	returnChan chan regexSearchReturn,
	re *regexp.Regexp,
	maxResultsPerFile int,
	negate bool,
) {
	searchTimes := 0

	retDocs := []types.DocumentWithSections{}
	for {
		info := <-infoChan
		if info.Exit {
			break
		}

		indexer.indexerLock.RLock()
		docContent, ok := indexer.documentIDToContentMap[info.DocumentID]
		indexer.indexerLock.RUnlock()
		if !ok {
			continue
		}

		if len(info.Lines) == 0 {
			// 没有候选行的时候，搜索全文
			// 由于我们的引擎不会更改文档内容，所以这行不需要上锁
			intIndex := re.FindAllIndex(*docContent, maxResultsPerFile)

			searchTimes++
			if !negate {
				if intIndex == nil {
					continue
				}
				index := []types.Section{}
				for _, idx := range intIndex {
					if len(idx) > 0 && idx[1] > idx[0] {
						index = append(index,
							types.Section{
								Start: uint32(idx[0]),
								End:   uint32(idx[1]),
							})
					}
				}
				retDocs = append(retDocs, types.DocumentWithSections{
					DocumentID: info.DocumentID,
					Sections:   index,
				})
			} else {
				if intIndex != nil {
					continue
				}
				retDocs = append(retDocs, types.DocumentWithSections{
					DocumentID: info.DocumentID,
				})
			}
		} else {
			sections := []types.Section{}
			for _, line := range info.Lines {
				lineContent, lineStart := indexer.GetLineContent(info.DocumentID, line)

				// 搜行
				intIndex := re.FindAllIndex(lineContent, -1)

				searchTimes++
				if intIndex == nil {
					continue
				}
				for _, idx := range intIndex {
					if len(idx) > 0 && idx[1] > idx[0] {
						sections = append(sections,
							types.Section{
								Start: uint32(idx[0]) + lineStart,
								End:   uint32(idx[1]) + lineStart,
							})
					}
				}
			}
			if len(sections) == 0 {
				continue
			}
			retDocs = append(retDocs, types.DocumentWithSections{
				DocumentID: info.DocumentID,
				Sections:   sections,
			})
		}
	}

	returnChan <- regexSearchReturn{
		Docs:          retDocs,
		SearchedTimes: searchTimes,
	}
}
