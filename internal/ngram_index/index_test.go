package ngram_index

import (
	"flag"
	"fmt"
	"sort"
	"testing"

	"github.com/huichen/kunlun/internal/common_types"
)

func TestIndexDocument(t *testing.T) {
	content := "this is a document"

	index := NewNgramIndex()

	err := index.IndexDocument(DocumentData{
		DocumentID:    1,
		Content:       []byte(content),
		SymbolEntries: nil,
	})

	if err != nil {
		logger.Fatal(err)
	}

	locations := []LocsWithKey{}
	for _, v := range index.indexMap.index {
		if len(*v.documents) > 1 {
			logger.Fatal(err)
		}
		locations = append(locations, LocsWithKey{
			Key:       v.key,
			Locations: (*v.documents)[0].SortedStartLocations,
		})
	}

	sort.Sort(SortLocsWithKey(locations))

	fmt.Println(content)
	for _, loc := range locations {
		fmt.Printf("key = %s, docs = %+v\n", IndexKeyToPrettyString(loc.Key), loc.Locations)
	}
}

func TestIndexShortDocument(t *testing.T) {
	content := "thi"

	index := NewNgramIndex()

	err := index.IndexDocument(DocumentData{
		DocumentID:    1,
		Content:       []byte(content),
		SymbolEntries: nil,
	})

	if err != nil {
		logger.Fatal(err)
	}

	locations := []LocsWithKey{}
	for _, v := range index.indexMap.index {
		if len(*v.documents) > 1 {
			logger.Fatal(err)
		}
		locations = append(locations, LocsWithKey{
			Key:       v.key,
			Locations: (*v.documents)[0].SortedStartLocations,
		})
	}

	sort.Sort(SortLocsWithKey(locations))

	fmt.Println(content)
	for _, loc := range locations {
		fmt.Printf("key = %s, docs = %+v\n", IndexKeyToPrettyString(loc.Key), loc.Locations)
	}

	content = "th"

	index = NewNgramIndex()

	err = index.IndexDocument(DocumentData{
		DocumentID:    1,
		Content:       []byte(content),
		SymbolEntries: nil,
	})

	if err != nil {
		logger.Fatal(err)
	}

	locations = []LocsWithKey{}
	for _, v := range index.indexMap.index {
		if len(*v.documents) > 1 {
			logger.Fatal(err)
		}
		locations = append(locations, LocsWithKey{
			Key:       v.key,
			Locations: (*v.documents)[0].SortedStartLocations,
		})
	}

	sort.Sort(SortLocsWithKey(locations))

	fmt.Println(content)
	for _, loc := range locations {
		fmt.Printf("key = %s, docs = %+v\n", IndexKeyToPrettyString(loc.Key), loc.Locations)
	}

	content = "t"

	index = NewNgramIndex()

	err = index.IndexDocument(DocumentData{
		DocumentID:    1,
		Content:       []byte(content),
		SymbolEntries: nil,
	})

	if err != nil {
		logger.Fatal(err)
	}

	locations = []LocsWithKey{}
	for _, v := range index.indexMap.index {
		if len(*v.documents) > 1 {
			logger.Fatal(err)
		}
		locations = append(locations, LocsWithKey{
			Key:       v.key,
			Locations: (*v.documents)[0].SortedStartLocations,
		})
	}

	sort.Sort(SortLocsWithKey(locations))

	fmt.Println(content)
	for _, loc := range locations {
		fmt.Printf("key = %s, docs = %+v\n", IndexKeyToPrettyString(loc.Key), loc.Locations)
	}
}

type LocsWithKey struct {
	Key       IndexKey
	Locations []uint32
}

type SortLocsWithKey []LocsWithKey

func (s SortLocsWithKey) Len() int {
	return len(s)
}
func (s SortLocsWithKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortLocsWithKey) Less(i, j int) bool {
	if s[i].Locations[0] == s[j].Locations[0] {
		return s[i].Key < s[j].Key
	}

	return s[i].Locations[0] < s[j].Locations[0]
}

func TestIndexCTags(t *testing.T) {
	flag.Parse()

	content := "a1 a2 a3\nb1 b2 b3\nc1 c2 c3"

	index := NewNgramIndex()

	err := index.IndexDocument(DocumentData{
		DocumentID: 1,
		Content:    []byte(content),
		SymbolEntries: []*common_types.CTagsEntry{
			{Sym: "a2", Line: 1},
			{Sym: "b1", Line: 2},
			{Sym: "c3", Line: 3},
		},
	})

	if err != nil {
		logger.Fatal(err)
	}

	locations := []LocsWithKey{}
	for _, v := range index.indexMap.index {
		if len(*v.documents) > 1 {
			logger.Fatal(err)
		}
		locations = append(locations, LocsWithKey{
			Key:       v.key,
			Locations: (*v.documents)[0].SortedStartLocations,
		})
	}
	sort.Sort(SortLocsWithKey(locations))

	fmt.Println(content)
	for _, loc := range locations {
		fmt.Printf("key = %s, docs = %+v\n", IndexKeyToPrettyString(loc.Key), loc.Locations)
	}
}

func TestIndexSkip(t *testing.T) {
	flag.Parse()

	content := "a1 a2 a3\nb1 b2 b3\nc1 c2 c3"

	index := NewNgramIndex()

	err := index.IndexDocument(DocumentData{
		DocumentID: 1,
		Content:    []byte(content),
		SymbolEntries: []*common_types.CTagsEntry{
			{Sym: "a2", Line: 1},
			{Sym: "b", Line: 2},
			{Sym: "c2 ", Line: 3},
		},
		SkipIndexUnigram: true,
	})

	if err != nil {
		logger.Fatal(err)
	}

	locations := []LocsWithKey{}
	for _, v := range index.indexMap.index {
		if len(*v.documents) > 1 {
			logger.Fatal(err)
		}
		locations = append(locations, LocsWithKey{
			Key:       v.key,
			Locations: (*v.documents)[0].SortedStartLocations,
		})
	}
	sort.Sort(SortLocsWithKey(locations))

	fmt.Println(content)
	for _, loc := range locations {
		fmt.Printf("key = %s, docs = %+v\n", IndexKeyToPrettyString(loc.Key), loc.Locations)
	}
}
func TestIndexCTagsMismatch(t *testing.T) {
	flag.Parse()

	content := "a1 a2 a3\nb1 b2 b3\nc1 c2 c3"

	index := NewNgramIndex()

	err := index.IndexDocument(DocumentData{
		DocumentID: 1,
		Content:    []byte(content),
		SymbolEntries: []*common_types.CTagsEntry{
			{Sym: "a4", Line: 1},
			{Sym: "a2", Line: 1},
			{Sym: "a5", Line: 1},
			{Sym: "b0", Line: 2},
			{Sym: "b1", Line: 2},
			{Sym: "b4", Line: 2},
			{Sym: "c7", Line: 3},
			{Sym: "c3", Line: 3},
			{Sym: "c9", Line: 3},
		},
	})

	if err != nil {
		logger.Fatal(err)
	}

	locations := []LocsWithKey{}
	for _, v := range index.indexMap.index {
		if len(*v.documents) > 1 {
			logger.Fatal(err)
		}
		locations = append(locations, LocsWithKey{
			Key:       v.key,
			Locations: (*v.documents)[0].SortedStartLocations,
		})
	}
	sort.Sort(SortLocsWithKey(locations))

	fmt.Println(content)
	for _, loc := range locations {
		fmt.Printf("key = %s, docs = %+v\n", IndexKeyToPrettyString(loc.Key), loc.Locations)
	}
}
