package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/huichen/kunlun/internal/query"

	"github.com/huichen/kunlun/pkg/types"
)

func TestMergeSections(t *testing.T) {
	flag.Parse()

	sections := []*[]types.Section{
		{{1, 2}},
		{{3, 4}},
		{{7, 8}},
	}
	r, err := mergeSections(sections)
	assert.Nil(t, err)
	assert.Equal(t, []types.Section{
		{1, 2}, {3, 4}, {7, 8},
	}, r)

	sections = []*[]types.Section{
		{{1, 2}, {3, 4}, {2, 4}, {7, 8}},
		{{1, 2}, {3, 4}, {7, 8}},
	}
	_, err = mergeSections(sections)
	logger.Error(err)
	assert.NotNil(t, err)

	sections = []*[]types.Section{
		{{1, 2}, {3, 4}, {7, 8}},
		{{1, 2}, {2, 4}, {4, 8}, {11, 12}},
	}
	r, _ = mergeSections(sections)
	assert.Equal(t, []types.Section{
		{1, 2}, {2, 4}, {3, 4}, {4, 8}, {7, 8}, {11, 12},
	}, r)

	sections = []*[]types.Section{
		{{1, 2}, {3, 2}, {7, 8}},
		{{1, 2}, {2, 4}, {4, 8}, {11, 12}},
	}
	r, _ = mergeSections(sections)
	logger.Error(err)
	assert.NotNil(t, err)
}

func TestMergeQueries1(t *testing.T) {
	flag.Parse()

	queries := []*query.Query{
		{ID: 0},
		{ID: 1},
		{ID: 2},
	}
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {3, 4}, {7, 8},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
				},
			},
		},
	}

	r, err := mergeQueries(&context, queries, true)
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {7, 8}, {10, 11}, {11, 12},
			},
		},
		{
			DocumentID: 2,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {4, 5}, {5, 6}, {10, 11}, {11, 12},
			},
		},
		{
			DocumentID: 3,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
			},
		},
	}, r)

	r, err = mergeQueries(&context, queries, false)
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {7, 8}, {10, 11}, {11, 12},
			},
		},
	}, r)

	queries = []*query.Query{
		{ID: 0},
		{ID: 1, Negate: true},
	}
	context = Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 14}, {11, 18},
						},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 14}, {11, 18},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {3, 14}, {7, 18},
						},
					},
				},
			},
		},
	}

	r, err = mergeQueries(&context, queries, false)
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 3,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {4, 14}, {11, 18},
			},
		}}, r)

}

func TestMergeQueriesError(t *testing.T) {
	flag.Parse()

	queries := []*query.Query{
		{ID: 0},
		{ID: 1},
		{ID: 2},
	}
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {3, 4}, {7, 8},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
				},
			},
		},
	}

	_, err := mergeQueries(&context, queries, true)
	logger.Error(err)
	assert.NotNil(t, err)
}

func TestMergeQueriesEmpty(t *testing.T) {
	flag.Parse()

	queries := []*query.Query{
		{ID: 0},
		{ID: 1},
		{ID: 2},
	}
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
				},
			},
		},
	}

	r, _ := mergeQueries(&context, queries, true)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {4, 5}, {5, 6}, {10, 11}, {11, 12},
			},
		},
		{
			DocumentID: 2,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {4, 5}, {5, 6}, {10, 11}, {11, 12},
			},
		},
		{
			DocumentID: 3,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
			},
		},
	}, r)

	r, _ = mergeQueries(&context, queries, false)
	assert.Equal(t, []types.DocumentWithSections{}, r)
}

func TestMergeQueriesNegate(t *testing.T) {
	flag.Parse()

	queries := []*query.Query{
		{ID: 0},
		{ID: 1, Negate: true},
		{ID: 2},
	}

	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {3, 4}, {7, 8},
						},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
						},
					},
				},
			},
		},
	}

	r, err := mergeQueries(&context, queries, false)
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 3,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {5, 6}, {10, 11}, {11, 12},
			},
		},
	}, r)

}

func TestMergeMoreQueries(t *testing.T) {
	flag.Parse()

	queries := []*query.Query{
		{ID: 0},
		{ID: 1},
	}

	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 4,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {3, 4}, {7, 8},
						},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections: []types.Section{
							{1, 2}, {2, 3}, {4, 5}, {11, 12},
						},
					},
				},
			},
		},
	}

	r, err := mergeQueries(&context, queries, true)
	assert.Nil(t, err)
	assert.Equal(t, []types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {4, 5}, {11, 12},
			},
		},
		{
			DocumentID: 4,
			Sections: []types.Section{
				{1, 2}, {2, 3}, {3, 4}, {7, 8},
			},
		},
	}, r)

}
