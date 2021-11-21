package kls

import (
	"kunlun/internal/searcher"
	"kunlun/pkg/types"
)

const (
	contextLines        = 2
	pageSize            = 1000
	maxDocumentsPerRepo = 100
)

func (kls *KLS) search() {
	if kls.inSearching {
		kls.printSuggestion("检索中，请勿重复搜索", "")
		return
	}
	kls.inSearching = true

	inputText := kls.inputField.GetText()

	// 清空内容
	kls.fileList.Clear()
	kls.fileContent.Clear()

	kls.printSuggestion("%s", "[yellow]搜索中，请稍等 ...")

	var err error
	kls.response, err = kls.kgn.Search(types.SearchRequest{
		Query:               inputText,
		ReturnLineContent:   true,
		NumContextLines:     contextLines,
		MaxLinesPerDocument: 30,
		PageSize:            pageSize,
		PageNum:             0,
		TimeoutInMs:         2000, // 超时 1 秒
		MaxDocumentsPerRepo: maxDocumentsPerRepo,
	})
	kls.inSearching = false

	if err != nil {
		kls.printSuggestion("[red]搜索错误：%s", err)
		return
	}

	kls.currentRepoID = -1
	kls.currentFileID = -1

	if kls.response == nil {
		return
	}

	kls.fileContent.Clear()

	q, err := searcher.ParseQuery(inputText)
	if err != nil {
		kls.printSuggestion("[red]query解析错误：%s", err)
	} else if len(kls.response.Repos) == 0 {
		kls.printSuggestion("query=%s, 搜索耗时%.3fms, 没有检索到任何结果",
			q.TrimmedQuery, float32(kls.response.RecallDurationInMicroSeconds)/1000.0)
	} else if len(kls.response.Repos) == pageSize {
		kls.printSuggestion("query=%s, 搜索耗时%.3fms, 返回前%d个仓库的%d行结果",
			q.TrimmedQuery, float32(kls.response.RecallDurationInMicroSeconds)/1000.0,
			len(kls.response.Repos), kls.response.NumLines)
	} else {
		kls.printSuggestion("query=%s, 搜索耗时%.3fms, 检索到%d个仓库的%d行结果",
			q.TrimmedQuery, float32(kls.response.RecallDurationInMicroSeconds)/1000.0,
			len(kls.response.Repos), kls.response.NumLines)
	}

	kls.redrawRepoList()
	kls.redrawFileList()
	kls.redrawFileContent()
	kls.app.Draw()
}
