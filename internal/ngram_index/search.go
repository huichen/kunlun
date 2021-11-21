package ngram_index

import "sort"

type DocumentWithLocations struct {
	DocumentID     uint64
	StartLocations []uint32
}

func (index *NgramIndex) SearchTwoKeys(
	key1 IndexKey, key2 IndexKey, distance uint32,
	shouldDocBeRecalled func(uint64) bool,
	isSymbol bool,
) ([]DocumentWithLocations, error) {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	var documents1, documents2 *SortedKeyedDocuments
	var ok bool

	// 对 key1 或者 key2 找不到文档的特殊情况做处理
	if isSymbol {
		if documents1, ok = index.symbolIndexMap[key1]; !ok {
			return nil, nil
		}
		if documents2, ok = index.symbolIndexMap[key2]; !ok {
			return nil, nil
		}
	} else {
		if documents1, ok = index.indexMap[key1]; !ok {
			return nil, nil
		}
		if documents2, ok = index.indexMap[key2]; !ok {
			return nil, nil
		}
	}
	for len(*documents1) == 0 || len(*documents2) == 0 {
		return nil, nil
	}

	// 两个指针递增，找到相同的 documentID
	retDocuments := []DocumentWithLocations{}
	var doc1Index, doc2Index int
	var docID1, docID2 uint64
	for doc1Index < len(*documents1) && doc2Index < len(*documents2) {
		docID1 = (*documents1)[doc1Index].DocumentID
		docID2 = (*documents2)[doc2Index].DocumentID

		// 跳过不相同的 ID
		if docID1 < docID2 {
			doc1Index++
			continue
		} else if docID1 > docID2 {
			doc2Index++
			continue
		}

		// 首先检查文档是否应该被检索到
		if shouldDocBeRecalled != nil {
			if !shouldDocBeRecalled(docID1) {
				doc1Index++
				doc2Index++
				continue
			}
		}

		// ID 相同的情况
		locations := findStartLocationWithKeyDistance(
			(*documents1)[doc1Index].SortedStartLocations,
			(*documents2)[doc2Index].SortedStartLocations,
			distance,
		)

		func() {
			if len(locations) > 0 {
				retDocuments = append(retDocuments, DocumentWithLocations{
					DocumentID:     docID1,
					StartLocations: locations,
				})
			}
		}()

		doc1Index++
		doc2Index++
	}

	if len(retDocuments) == 0 {
		return nil, nil
	}

	return retDocuments, nil
}

func (index *NgramIndex) SearchOneKey(
	key IndexKey,
	shouldDocBeRetrived func(uint64) bool,
	isSymbol bool,
) ([]DocumentWithLocations, error) {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	var documents *SortedKeyedDocuments
	var ok bool

	// 对 key1 或者 key2 找不到文档的特殊情况做处理
	if isSymbol {
		if documents, ok = index.symbolIndexMap[key]; !ok {
			return nil, nil
		}
	} else {
		if documents, ok = index.indexMap[key]; !ok {
			return nil, nil
		}
	}
	for len(*documents) == 0 {
		return nil, nil
	}

	// 两个指针递增，找到相同的 documentID
	retDocuments := []DocumentWithLocations{}
	for docIndex := 0; docIndex < len(*documents); docIndex++ {
		doc := (*documents)[docIndex]
		docID := doc.DocumentID

		// 检查文档是否应该被检索到
		if shouldDocBeRetrived != nil {
			if !shouldDocBeRetrived(docID) {
				continue
			}
		}

		retDocuments = append(retDocuments, DocumentWithLocations{
			DocumentID:     docID,
			StartLocations: doc.SortedStartLocations,
		})
	}

	return retDocuments, nil
}

func findStartLocationWithKeyDistance(startLocations1 []uint32, startLocations2 []uint32, distance uint32) []uint32 {
	if len(startLocations1) == 0 || len(startLocations2) == 0 {
		return nil
	}
	lenStartLocations1 := len(startLocations1)
	lenStartLocations2 := len(startLocations2)

	retLocations := make([]uint32, 0, lenStartLocations1+lenStartLocations2)

	if distance == 0 {
		retLocations = append(retLocations, startLocations1...)
		retLocations = append(retLocations, startLocations2...)
		sort.Slice(retLocations, func(i, j int) bool {
			return retLocations[i] < retLocations[j]
		})
		return retLocations
	}

	var index1, index2 int
	for index1 < lenStartLocations1 && index2 < lenStartLocations2 {
		loc1 := startLocations1[index1]
		loc2 := startLocations2[index2]

		expectedLoc2 := loc1 + distance
		if expectedLoc2 < loc2 {
			index1++
			continue
		} else if expectedLoc2 > loc2 {
			index2++
			continue
		}

		// 距离匹配的情况
		retLocations = append(retLocations, loc1)
		index1++
		index2++
	}

	if len(retLocations) == 0 {
		return nil
	}

	return retLocations
}
