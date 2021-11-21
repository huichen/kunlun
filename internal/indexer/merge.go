package indexer

// 对两个有序数组取交集，如果第二个数组为空，则返回第一个数组
func andMerge(docID1 []uint64, docID2 []DocumentWithLines) []DocumentWithLines {
	if len(docID2) == 0 || len(docID1) == 0 {
		return docID2
	}

	retIDs := []DocumentWithLines{}

	var index1, index2 int
	for index1 < len(docID1) && index2 < len(docID2) {
		id1 := docID1[index1]
		id2 := docID2[index2].DocumentID

		if id1 < id2 {
			index1++
		} else if id1 > id2 {
			index2++
		}

		retIDs = append(retIDs, docID2[index2])
		index1++
		index2++
	}

	return retIDs
}
