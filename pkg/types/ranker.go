package types

// 排序器接口，你可以实现自己的排序器，然后在 SearchRequest 中传入
type Ranker interface {
	Rank(response *SearchResponse)
}
