package indexer

import (
	"errors"

	"kunlun/internal/ngram_index"
)

// 在返回的二键匹配文档中，用 keyword 精准匹配过滤：
// 从 docs 的 StartLocations-offset 里查找所有匹配 keyword 的文档和起始位置
func (indexer *Indexer) filterDocumentsWithFullMatch(
	contentMap *map[uint64]*[]byte,
	docs []ngram_index.DocumentWithLocations,
	offset uint32,
	keyword []byte,
	caseSensitive bool,
) ([]ngram_index.DocumentWithLocations, error) {
	indexer.indexerLock.RLock()
	defer indexer.indexerLock.RUnlock()
	newDocIndex := 0
	for _, doc := range docs {
		docID := doc.DocumentID

		docContent, ok := (*contentMap)[docID]
		if !ok {
			return nil, errors.New("文档没找到")
		}

		newIndex := 0
		for _, loc := range doc.StartLocations {
			newLoc := int(loc) - int(offset)
			for matchIndex, c := range keyword {
				contentIndex := matchIndex + newLoc
				if contentIndex < 0 || contentIndex >= len(*docContent) {
					// 检查是否越界
					goto breakPoint
				}
				if !matchWithCase(c, (*docContent)[contentIndex], caseSensitive) {
					goto breakPoint
				}
			}
			doc.StartLocations[newIndex] = loc - offset
			newIndex++
		breakPoint:
		}
		doc.StartLocations = doc.StartLocations[:newIndex]

		if len(doc.StartLocations) > 0 {
			docs[newDocIndex].DocumentID = docID
			docs[newDocIndex].StartLocations = doc.StartLocations
			newDocIndex++
		}
	}

	docs = docs[:newDocIndex]

	return docs, nil
}

func matchWithCase(b1 byte, b2 byte, caseSensitive bool) bool {
	if !caseSensitive && b1 >= 'A' && b1 <= 'Z' {
		b1 = b1 + Atoa
	}

	if !caseSensitive && b2 >= 'A' && b2 <= 'Z' {
		b2 = b2 + Atoa
	}

	return b1 == b2
}
