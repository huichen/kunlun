package ngram_index

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	runeMask = 1<<21 - 1
)

func RuneNgramToIndexKey(ngram RuneNgram) IndexKey {
	return IndexKey(uint64(ngram[0])<<42 | uint64(ngram[1])<<21 | uint64(ngram[2]))
}

func RuneBigramToIndexKey(ngram RuneNgram) IndexKey {
	return IndexKey(uint64(ngram[0])<<42 | uint64(ngram[1])<<21)
}

func RuneUnigramToIndexKey(ngram RuneNgram) IndexKey {
	return IndexKey(uint64(ngram[0]) << 42)
}

func IndexKeyToString(key IndexKey) string {
	rune0 := rune((key >> 42) & runeMask)
	rune1 := rune((key >> 21) & runeMask)
	rune2 := rune(key & runeMask)

	runeSlice := []rune{}
	if rune0 != 0 {
		runeSlice = append(runeSlice, rune0)
	}

	if rune1 != 0 {
		runeSlice = append(runeSlice, rune1)
	}

	if rune2 != 0 {
		runeSlice = append(runeSlice, rune2)
	}

	return string(runeSlice)
}

func IndexKeyToPrettyString(key IndexKey) string {
	str := IndexKeyToString(key)

	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\t", "\\t")

	return str
}

func BytesToIndexKey(bytes []byte) (IndexKey, error) {
	rs := DecodeRunes(bytes)
	if len(rs) > 3 || len(rs) == 0 {
		return 0, errors.New("字符串中的 rune 个数需要在 [1, 3] 之间")
	} else if len(rs) == 2 {
		return RuneBigramToIndexKey(RuneNgram{rs[0], rs[1]}), nil
	} else if len(rs) == 1 {
		return RuneUnigramToIndexKey(RuneNgram{rs[0]}), nil
	}

	return RuneNgramToIndexKey(RuneNgram{rs[0], rs[1], rs[2]}), nil
}

func StringToIndexKey(str string) (IndexKey, error) {
	return BytesToIndexKey([]byte(str))
}

func RuneSliceToIndexKey(runes []rune) (IndexKey, error) {
	if len(runes) == 0 {
		return 0, errors.New("字符串中的 rune 个数不能为 0")
	} else if len(runes) == 2 {
		return RuneBigramToIndexKey(RuneNgram{runes[0], runes[1]}), nil
	} else if len(runes) == 1 {
		return RuneUnigramToIndexKey(RuneNgram{runes[0]}), nil
	}
	return RuneNgramToIndexKey(RuneNgram{runes[0], runes[1], runes[2]}), nil
}

func DecodeRunes(content []byte) []rune {
	ret := []rune{}

	for {
		r, rs := utf8.DecodeRune(content)
		if rs == 0 {
			break
		}
		content = content[rs:]

		ret = append(ret, r)
	}

	return ret
}

type SortDocumentWithLocations []DocumentWithLocations

func (s SortDocumentWithLocations) Len() int {
	return len(s)
}
func (s SortDocumentWithLocations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortDocumentWithLocations) Less(i, j int) bool {
	return s[i].DocumentID < s[j].DocumentID
}

func (s SortedKeyedDocuments) Len() int {
	return len(s)
}
func (s SortedKeyedDocuments) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortedKeyedDocuments) Less(i, j int) bool {
	return s[i].DocumentID < s[j].DocumentID
}
