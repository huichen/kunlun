package ngram_index

type SingleDocIndexMap struct {
	index []SingleDocIndexMapItem
}

type SingleDocIndexMapItem struct {
	key       IndexKey
	locations *[]uint32
}

func (im *SingleDocIndexMap) find(key IndexKey) (int, bool) {
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

func (im *SingleDocIndexMap) insert(key IndexKey, loc uint32) {
	idx, found := im.find(key)

	if !found {
		if idx == len(im.index) {
			im.index = append(im.index, SingleDocIndexMapItem{
				key:       key,
				locations: &[]uint32{loc},
			})
			return
		}
		im.index = append(im.index[:idx+1], im.index[idx:]...)
		im.index[idx] = SingleDocIndexMapItem{
			key:       key,
			locations: &[]uint32{loc},
		}
		return
	}
	locations := im.index[idx].locations

	if (*locations)[len(*locations)-1] >= loc {
		logger.Fatal("不能插入一个比最后一个元素小的元素")
	}

	*locations = append(*locations, loc)
}
