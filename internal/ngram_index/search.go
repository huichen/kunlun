package ngram_index

import "sort"

type DocumentWithLocations struct {
	DocumentID     uint64
	StartLocations []uint32
}

// 在索引中检索所有包含 key1 和 key2，且两者距离（头到头）为 distance 的文档和 key1 的起始位置
// 使用 shouldDocBeRecalled 钩子函数对文档进行过滤
// 当 isSymbol 为 true 时搜索符号索引表，否则搜索全文索引表
func (index *NgramIndex) SearchTwoKeys(
	key1 IndexKey, key2 IndexKey, distance uint32,
	shouldDocBeRecalled func(uint64) bool,
	isSymbol bool,
) ([]DocumentWithLocations, error) {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	var documents1, documents2 *SortedKeyedDocuments

	// 对 key1 或者 key2 找不到文档的特殊情况做处理
	if isSymbol {
		idx1, found1 := index.symbolIndexMap.find(key1)
		if !found1 {
			return nil, nil
		}

		idx2, found2 := index.symbolIndexMap.find(key2)
		if !found2 {
			return nil, nil
		}

		documents1 = index.symbolIndexMap.index[idx1].documents
		documents2 = index.symbolIndexMap.index[idx2].documents
	} else {
		idx1, found1 := index.indexMap.find(key1)
		if !found1 {
			return nil, nil
		}

		idx2, found2 := index.indexMap.find(key2)
		if !found2 {
			return nil, nil
		}

		documents1 = index.indexMap.index[idx1].documents
		documents2 = index.indexMap.index[idx2].documents
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
			*(*documents1)[doc1Index].SortedStartLocations,
			*(*documents2)[doc2Index].SortedStartLocations,
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

// 在索引中检索所有包含 key 的文档和 key 的起始位置
// 使用 shouldDocBeRecalled 钩子函数对文档进行过滤
// 当 isSymbol 为 true 时搜索符号索引表，否则搜索全文索引表
func (index *NgramIndex) SearchOneKey(
	key IndexKey,
	shouldDocBeRetrived func(uint64) bool,
	isSymbol bool,
) ([]DocumentWithLocations, error) {
	index.indexLock.RLock()
	defer index.indexLock.RUnlock()

	var documents *SortedKeyedDocuments

	// 对 key1 或者 key2 找不到文档的特殊情况做处理
	if isSymbol {
		idx, found := index.symbolIndexMap.find(key)
		if !found {
			return nil, nil
		}
		documents = index.symbolIndexMap.index[idx].documents
	} else {
		idx, found := index.indexMap.find(key)
		if !found {
			return nil, nil
		}
		documents = index.indexMap.index[idx].documents
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
			StartLocations: *doc.SortedStartLocations,
		})
	}

	return retDocuments, nil
}

// 两个正序排列的起始位置数组，查找所有 startLocations1 元素，满足：
// 在 startLocations2 中存在到该元素距离等于 distance 元素
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
