package ngram_index

type IndexMap struct {
	index []IndexMapItem
}

type IndexMapItem struct {
	key       IndexKey
	documents *SortedKeyedDocuments
}

func (im *IndexMap) find(key IndexKey) (int, bool) {
	b := 0
	e := len(im.index) - 1

	if len(im.index) == 0 {
		return 0, false
	} else if key < im.index[0].key {
		return 0, false
	} else if key == im.index[0].key {
		return 0, true
	} else if key > im.index[e].key {
		return e + 1, false
	} else if key == im.index[e].key {
		return e, true
	}

	for e > b+1 {
		m := (e + b) / 2
		mkey := im.index[m].key
		if key == mkey {
			return m, true
		} else if key < mkey {
			e = m
		} else {
			b = m
		}
	}

	return e, false
}

func (im *IndexMap) insert(key IndexKey, document KeyedDocument) {
	idx, found := im.find(key)

	if !found {
		if idx == len(im.index) {
			im.index = append(im.index, IndexMapItem{
				key:       key,
				documents: &SortedKeyedDocuments{document},
			})
			return
		}
		im.index = append(im.index[:idx+1], im.index[idx:]...)
		im.index[idx] = IndexMapItem{
			key:       key,
			documents: &SortedKeyedDocuments{document},
		}
		return
	}

	im.index[idx].documents.insert(document)

}

func (sd *SortedKeyedDocuments) find(docID uint64) (int, bool) {
	b := 0
	e := len(*sd) - 1

	if len(*sd) == 0 {
		return 0, false
	} else if docID < (*sd)[b].DocumentID {
		return 0, false
	} else if docID == (*sd)[b].DocumentID {
		return 0, true
	} else if docID > (*sd)[e].DocumentID {
		return e + 1, false
	} else if docID == (*sd)[e].DocumentID {
		return e, true
	}

	for e > b+1 {
		m := (e + b) / 2
		mID := (*sd)[m].DocumentID
		if docID == mID {
			return m, true
		} else if docID < mID {
			e = m
		} else {
			b = m
		}
	}

	return e, false
}

func (sd *SortedKeyedDocuments) insert(document KeyedDocument) {
	idx, found := sd.find(document.DocumentID)

	if found {
		logger.Fatal("不能重复索引文档")
	}

	if idx == len(*sd) {
		*sd = append(*sd, document)
		return
	}
	*sd = append((*sd)[:idx+1], (*sd)[idx:]...)
	(*sd)[idx] = document
}
