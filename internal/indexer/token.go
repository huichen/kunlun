package indexer

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/huichen/kunlun/internal/ngram_index"
)

// 从 token 得到两个搜索键和他们的距离，这两个搜索键是所有搜索键中频率最低的
// 如果 token 包含小于三个 rune，则返回一个搜索键（第二个为 0)
// 返回的第一个参数为 key1 在 token 中的起始位置
func (indexer *Indexer) getTwoKeysFromToken(token string) (uint32, ngram_index.IndexKey, ngram_index.IndexKey, uint32) {
	lowerCaseKeyword := []byte(strings.ToLower(token))
	runes := ngram_index.DecodeRunes(lowerCaseKeyword)
	if len(runes) == 0 {
		return 0, 0, 0, 0
	}

	keyMap := make(map[ngram_index.IndexKey]Range)
	keys := []keyWithFrequency{}
	iStart := uint32(0)
	for len(runes) > 0 {
		key, _ := ngram_index.RuneSliceToIndexKey(runes)
		if _, ok := keyMap[key]; !ok {
			keyMap[key] = Range{
				Start: iStart,
				End:   iStart,
			}
			keys = append(keys, keyWithFrequency{
				Key:  key,
				Freq: indexer.getKeyFrequency(key),
			})
		} else {
			r := keyMap[key]
			keyMap[key] = Range{
				Start: min(iStart, r.Start),
				End:   max(iStart, r.End),
			}
		}

		iStart += uint32(utf8.RuneLen(runes[0]))
		runes = runes[1:]
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].Freq < keys[j].Freq
	})

	if len(keys) == 0 {
		return 0, 0, 0, 0
	} else if len(keys) == 1 {
		return keyMap[keys[0].Key].Start, keys[0].Key, 0, 0
	}

	i := 0
	j := 1
	if keyMap[keys[i].Key].Start > keyMap[keys[j].Key].Start {
		i, j = j, i
	}
	return keyMap[keys[i].Key].Start, keys[i].Key, keys[j].Key, keyMap[keys[j].Key].End - keyMap[keys[i].Key].Start
}

type Range struct {
	Start uint32
	End   uint32
}

type keyWithFrequency struct {
	Key  ngram_index.IndexKey
	Freq uint64
}
