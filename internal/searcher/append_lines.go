package searcher

import (
	"errors"

	"github.com/huichen/kunlun/internal/common_types"
	"github.com/huichen/kunlun/internal/indexer"
	"github.com/huichen/kunlun/pkg/types"
)

// 将搜索结果中文档的 section 信息转化为行
func appendLinesToResponse(
	context *Context,
	idxr *indexer.Indexer,
	response *types.SearchResponse,
) error {
	infoChan := make(chan appendLinesToDocumentInfo, context.searcherOptions.AnnotatorProcessors)
	retChan := make(chan error, context.searcherOptions.AnnotatorProcessors)
	for i := 0; i < context.searcherOptions.AnnotatorProcessors; i++ {
		go appendLinesToDocumentWorker(context, idxr, infoChan, retChan)
	}

	totalDocs := 0
	for _, repo := range response.Repos {
		totalDocs += len(repo.Documents)
	}

	go func() {
		for _, repo := range response.Repos {
			for id := range repo.Documents {
				infoChan <- appendLinesToDocumentInfo{
					searchedDocument: &repo.Documents[id],
				}
			}
		}
	}()

	var err error
	for i := 0; i < totalDocs; i++ {
		workerErr := <-retChan
		if workerErr != nil {
			err = workerErr
		}
	}
	if err != nil {
		return err
	}
	for i := 0; i < context.searcherOptions.AnnotatorProcessors; i++ {
		infoChan <- appendLinesToDocumentInfo{exit: true}
	}

	totalLines := 0
	for _, repo := range response.Repos {
		numLinesInRepo := 0
		for id := range repo.Documents {
			numLinesInRepo += repo.Documents[id].NumLinesInDocument
		}

		repo.NumLinesInRepo = numLinesInRepo
		totalLines += numLinesInRepo
	}
	response.NumLines = totalLines
	return nil
}

type appendLinesToDocumentInfo struct {
	searchedDocument *types.SearchedDocument
	exit             bool
}

func appendLinesToDocumentWorker(
	context *Context,
	idxr *indexer.Indexer,
	infoChan chan appendLinesToDocumentInfo,
	retChan chan error) {

	for {
		info := <-infoChan
		if info.exit {
			return
		}
		err := appendLinesToDocument(context, idxr, info.searchedDocument)
		retChan <- err
	}
}

func appendLinesToDocument(
	context *Context,
	idxr *indexer.Indexer,
	searchedDocument *types.SearchedDocument,
) error {
	contextLines := context.request.NumContextLines
	if contextLines < 0 {
		contextLines = 0
	}

	var doc *common_types.DocumentWithSections
	var ok bool
	if doc, ok = context.docIDToDocumentWithSectionsMap[searchedDocument.DocumentID]; !ok {
		return errors.New("doc id 没找到")
	}

	lineOutputs, numResults, err := idxr.GetLinesFromSections(
		doc.DocumentID, doc.Sections,
		contextLines,
		context.request.MaxLinesPerDocument)
	if err != nil {
		return err
	}
	if err := context.checkTimeout(); err != nil {
		return err
	}

	searchedDocument.Lines = lineOutputs
	searchedDocument.NumLinesInDocument = numResults

	return nil
}
