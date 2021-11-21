package searcher

import (
	"github.com/huichen/kunlun/internal/indexer"

	"github.com/huichen/kunlun/pkg/types"
)

// 给搜索结果的文档添加文本内容
func appendContentToResponse(context *Context, idxr *indexer.Indexer, response *types.SearchResponse) {
	for _, repo := range response.Repos {
		for _, doc := range repo.Documents {
			appendContentToDocument(context, idxr, &doc)
		}
	}
}

func appendContentToDocument(context *Context, idxr *indexer.Indexer, doc *types.SearchedDocument) {
	for id, l := range doc.Lines {
		doc.Lines[id].Content, _ = idxr.GetLineContent(doc.DocumentID, uint32(l.LineNumber))
	}
	// 添加 tags
	if context.request.HightlightStartTag != "" && context.request.HightlightEndTag != "" {
		doc.Lines = addHighlightTags(doc.Lines, context.request.HightlightStartTag, context.request.HightlightEndTag)
	}
}

func addHighlightTags(lines []types.Line, start string, end string) []types.Line {
	for id := range lines {
		content := lines[id].Content
		highlights := lines[id].Highlights
		if len(highlights) == 0 {
			continue
		}

		var newContent []byte
		hindex := 0
		for index, c := range content {
			if hindex >= len(highlights) {
				newContent = append(newContent, c)
				continue
			}
			if index < int(highlights[hindex].Start) {
				newContent = append(newContent, c)
			} else if index == int(highlights[hindex].Start) {
				newContent = append(newContent, []byte(start)...)
				newContent = append(newContent, c)
				if highlights[hindex].Start+1 == highlights[hindex].End {
					newContent = append(newContent, []byte(end)...)
					hindex++
				}
			} else if index > int(highlights[hindex].Start) && index < int(highlights[hindex].End-1) {
				newContent = append(newContent, c)
			} else if index == int(highlights[hindex].End-1) {
				newContent = append(newContent, c)
				newContent = append(newContent, []byte(end)...)
				hindex++
			} else if index > int(highlights[hindex].End-1) {
				newContent = append(newContent, c)
			}
		}

		lines[id].Content = newContent
	}

	return lines
}
