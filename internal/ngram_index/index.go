package ngram_index

import (
	"bytes"
	"errors"
	"sort"
	"sync"
	"unicode/utf8"

	"kunlun/pkg/log"
	"kunlun/pkg/types"
)

const (
	// 'a' - 'A'
	Atoa = 32
)

var (
	logger = log.GetLogger()
)

// 保存了 IndexKey -> 文档和键起始位置数组 的反向索引
type NgramIndex struct {
	// 读写锁，保证索引是线程安全的
	indexLock sync.RWMutex

	// 反向索引
	indexMap       map[IndexKey]*SortedKeyedDocuments
	symbolIndexMap map[IndexKey]*SortedKeyedDocuments

	// 保存 IndexKey frequency
	keyFrequencies map[IndexKey]uint64

	// 反向索引的大小，只统计了 start locations 占用的字节数
	totalIndexSize uint64

	// 用于做递增判断
	maxDocumentID uint64

	// 排序触发次数
	sortTriggered uint64
}

// 键->文档 反向索引中的键
// 由 RuneNgram 中的三个 unicode rune 组合而成
// 组合方式如下（见 RuneNgramToIndexKey 函数）
// 20 ~ 0  位 = RuneNgram[2]
// 41 ~ 21 位 = RuneNgram[1]
// 62 ~ 42 位 = RuneNgram[0]
type IndexKey uint64

// 保存了 IndexKey 的 unicode rune
// 可以用 RuneNgramToIndexKey 和 IndexKeyToRuneNgram 两个函数相互切换
type RuneNgram [3]rune

// 对某个 IndexKey，保存了一个文档中所有该 key 的起始位置
type KeyedDocument struct {
	DocumentID uint64

	// IndexKey 在该文档中的所有起始位置，按照升序排列
	SortedStartLocations []uint32
}

// 按照 DocumentID 正序排列的数组
type SortedKeyedDocuments []KeyedDocument

func NewNgramIndex() *NgramIndex {
	return &NgramIndex{
		indexMap:       make(map[IndexKey]*SortedKeyedDocuments),
		symbolIndexMap: make(map[IndexKey]*SortedKeyedDocuments),
		keyFrequencies: make(map[IndexKey]uint64),
	}
}

func (index *NgramIndex) IndexDocument(documentID uint64, content []byte, entries []*types.CTagsEntry) error {
	// 临时缓存，从文档中获得后一次性加入索引
	indexCache := make(map[IndexKey]*[]uint32)
	symbolIndexCache := make(map[IndexKey]*[]uint32)

	var ngram RuneNgram
	var ngramSize [3]int
	var startLocation uint32
	var contentIndex int
	var entryIndex int
	var lineIndex int
	var lineStart int
	var entryInLine int
	var symbolLength int
	var addSymbolIndex bool
	var stopAddingSymbolIndex bool
	var err error

	if len(entries) > 0 && entries[0].Line-1 == 0 {
		entryInLine = bytes.Index(content, []byte(entries[0].Sym))
		if entryInLine == -1 {
			stopAddingSymbolIndex = true
		}
	}

	for contentIndex < len(content) {
		// 得到 content 的对应位置起的第一个 UTF8 字符和该字符的字节尺寸
		r, rsize := utf8.DecodeRune(content[contentIndex:])
		contentIndex += rsize

		// sanity check
		if r == 0 {
			return errors.New("文件是二进制文件")
		}

		// 将 r 转为小写
		if r >= 'A' && r <= 'Z' {
			r += Atoa
		}

		// 更新 ngram
		ngram[0], ngram[1], ngram[2] = ngram[1], ngram[2], r
		ngramSize[0], ngramSize[1], ngramSize[2] = ngramSize[1], ngramSize[2], rsize

		// 略过文件头部直到得到第一个 ngram
		if ngram[0] == 0 {
			continue
		}

		// 处理 symbols
		if entries != nil && !stopAddingSymbolIndex {
			var updateErr error
			stopAddingSymbolIndex, addSymbolIndex, lineIndex, lineStart, entryInLine, entryIndex, symbolLength, updateErr = updateSymbolLocations(content, startLocation, lineIndex, lineStart, entryInLine, entryIndex, entries)
			if updateErr != nil {
				stopAddingSymbolIndex = true
			}
		}

		// 从 ngram 得到反向索引的 key
		key := RuneNgramToIndexKey(ngram)
		keyBigram := RuneBigramToIndexKey(ngram)
		keyUnigram := RuneUnigramToIndexKey(ngram)

		// 将 key -> (docID, location) 加入索引
		addIndexToCache(key, startLocation, &indexCache)
		addIndexToCache(keyBigram, startLocation, &indexCache)
		addIndexToCache(keyUnigram, startLocation, &indexCache)
		if entries != nil && !stopAddingSymbolIndex && addSymbolIndex {
			entryLocation := startLocation - uint32(lineStart) - uint32(entryInLine)
			if ngramSize[0]+int(entryLocation) <= symbolLength {
				addIndexToCache(keyUnigram, startLocation, &symbolIndexCache)
			}
			if ngramSize[0]+ngramSize[1]+int(entryLocation) <= symbolLength {
				addIndexToCache(keyBigram, startLocation, &symbolIndexCache)
			}
			if ngramSize[0]+ngramSize[1]+ngramSize[2]+int(entryLocation) <= symbolLength {
				addIndexToCache(key, startLocation, &symbolIndexCache)
			}
		}
		// 更新 content 和 startLocation
		startLocation += uint32(ngramSize[0])
	}

	// 处理剩余的 ngram
	for i := 0; i < 2; i++ {
		ngram[0], ngram[1], ngram[2] = ngram[1], ngram[2], 0
		ngramSize[0], ngramSize[1], ngramSize[2] = ngramSize[1], ngramSize[2], 0
		if ngram[0] != 0 {
			keyUnigram := RuneUnigramToIndexKey(ngram)
			addIndexToCache(keyUnigram, startLocation, &indexCache)

			if ngram[1] != 0 {
				keyBigram := RuneBigramToIndexKey(ngram)
				addIndexToCache(keyBigram, startLocation, &indexCache)
			}

			// 对 symbol 做处理
			if entries != nil && !stopAddingSymbolIndex {
				var updateErr error
				stopAddingSymbolIndex, addSymbolIndex, lineIndex, lineStart, entryInLine, entryIndex, symbolLength, updateErr = updateSymbolLocations(content, startLocation, lineIndex, lineStart, entryInLine, entryIndex, entries)
				if updateErr != nil {
					err = updateErr
					stopAddingSymbolIndex = true
				}

				if entries != nil && !stopAddingSymbolIndex && addSymbolIndex {
					entryLocation := startLocation - uint32(lineStart) - uint32(entryInLine)
					if ngramSize[0]+int(entryLocation) <= symbolLength {
						addIndexToCache(keyUnigram, startLocation, &symbolIndexCache)
					}
					if ngram[1] != 0 {
						keyBigram := RuneBigramToIndexKey(ngram)
						if ngramSize[0]+ngramSize[1]+int(entryLocation) <= symbolLength {
							addIndexToCache(keyBigram, startLocation, &symbolIndexCache)
						}
					}
				}
			}

			startLocation += uint32(ngramSize[0])
		}
	}

	// 最后一次性将缓存添加到索引
	index.addCacheToIndex(documentID, &indexCache, false)
	index.addCacheToIndex(documentID, &symbolIndexCache, true)
	return err
}

func (index *NgramIndex) addCacheToIndex(documentID uint64, cache *map[IndexKey]*[]uint32, isSymbol bool) {
	index.indexLock.Lock()
	defer index.indexLock.Unlock()

	needSort := false
	if documentID <= index.maxDocumentID {
		// 在极少数情况下，需要对结果做重排序
		needSort = true
	}
	index.maxDocumentID = documentID

	for key, locations := range *cache {
		// 仅对 index cache 添加 key frequency
		if !isSymbol {
			if _, ok := index.keyFrequencies[key]; ok {
				index.keyFrequencies[key] = index.keyFrequencies[key] + 1
			} else {
				index.keyFrequencies[key] = 1
			}
		}

		index.totalIndexSize += uint64(len(*locations) * 4)

		var documents *SortedKeyedDocuments
		var ok bool
		if !isSymbol {
			documents, ok = index.indexMap[key]
		} else {
			documents, ok = index.symbolIndexMap[key]
		}

		if !ok {
			// key 不存在的情况
			docs := SortedKeyedDocuments{
				KeyedDocument{
					DocumentID:           documentID,
					SortedStartLocations: *locations,
				},
			}
			if !isSymbol {
				index.indexMap[key] = &docs
			} else {
				index.symbolIndexMap[key] = &docs
			}
			continue
		}

		documentsNeedSort := false
		if needSort {
			// 进一步检查是否真的需要重排序
			if documentID < (*documents)[len(*documents)-1].DocumentID {
				documentsNeedSort = true
			}
		}

		*documents = append(*documents, KeyedDocument{
			DocumentID:           documentID,
			SortedStartLocations: *locations,
		})

		if documentsNeedSort {
			// 极少触发，不应该对索引速度产生显著影响
			index.sortTriggered++
			sort.Sort(*documents)
		}
	}
}

func addIndexToCache(key IndexKey, start uint32, cache *map[IndexKey]*[]uint32) {
	var locations *[]uint32
	var ok bool
	if locations, ok = (*cache)[key]; !ok {
		// key 不存在的情况
		locs := []uint32{start}
		(*cache)[key] = &locs
		return
	}

	*locations = append(*locations, start)
}

func updateSymbolLocations(
	content []byte, startLocation uint32, lineIndex int, lineStart int, entryInLine int, entryIndex int, entries []*types.CTagsEntry,
) (bool, bool, int, int, int, int, int, error) {
	var stopAddingSymbolIndex bool
	addSymbolIndex := false
	var err error
	var symbolLength int
	if entryIndex >= len(entries) {
		stopAddingSymbolIndex = true
		goto continuePoint
	}

	// 对换行符做处理
	if content[startLocation] == '\n' {
		lineStart = int(startLocation + 1)
		if lineStart < len(content) && entries[entryIndex].Line-1 == lineIndex+1 {
			entryInLine = bytes.Index(content[lineStart:], []byte(entries[entryIndex].Sym))
			if entryInLine == -1 {
				err = errors.New("无法添加 symbol")
			}
		}

		lineIndex++
		goto continuePoint
	}

	if lineIndex < entries[entryIndex].Line-1 {
		goto continuePoint
	}

	// 行内
	symbolLength = len(entries[entryIndex].Sym)
	if lineIndex == entries[entryIndex].Line-1 {
		if lineStart+entryInLine > int(startLocation) {
			goto continuePoint
		} else if lineStart+entryInLine+len(entries[entryIndex].Sym) <= int(startLocation) {
			goto continuePoint
		} else if lineStart+entryInLine+len(entries[entryIndex].Sym)-1 == int(startLocation) {
			entryIndex++
		}
		addSymbolIndex = true
	} else {
		goto continuePoint
	}

continuePoint:

	return stopAddingSymbolIndex, addSymbolIndex, lineIndex, lineStart, entryInLine, entryIndex, symbolLength, err
}
