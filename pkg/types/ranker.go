package types

type Ranker interface {
	Rank(response *SearchResponse)
}
