package indexer

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/huichen/kunlun/internal/ngram_index"
)

func (indexer *Indexer) getTwoKeysFromToken(token string) (uint32, ngram_index.IndexKey, ngram_index.IndexKey, uint32) {
	lowerCaseKeyword := []byte(strings.ToLower(token))
	runes := ngram_index.DecodeRunes(lowerCaseKeyword)
	if len(runes) == 0 {
		return 0, 0, 0, 0
	}

	keyMap := make(map[ngram_index.IndexKey]Range)
	keys := []KeyWithFrequency{}
	iStart := uint32(0)
	for len(runes) > 0 {
		key, _ := ngram_index.RuneSliceToIndexKey(runes)
		if _, ok := keyMap[key]; !ok {
			keyMap[key] = Range{
				Start: iStart,
				End:   iStart,
			}
			keys = append(keys, KeyWithFrequency{
				Key:  key,
				Freq: indexer.GetKeyFrequency(key),
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

type KeyWithFrequency struct {
	Key  ngram_index.IndexKey
	Freq uint64
}
