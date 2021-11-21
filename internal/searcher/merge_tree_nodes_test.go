package searcher

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	"kunlun/internal/query"
	"kunlun/pkg/types"
)

func TestMergeTreeNodesAnd(t *testing.T) {
	flag.Parse()

	q, _ := query.Parse("(a and b)")
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, &[]types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 4}},
		},
	}, context.query.QueryResults[2])
}

func TestMergeTreeNodesOr(t *testing.T) {
	flag.Parse()

	q, _ := query.Parse("(a or b)")
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 4,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, &[]types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 4}},
		},
		{
			DocumentID: 2,
			Sections:   []types.Section{{0, 4}},
		},
		{
			DocumentID: 4,
			Sections:   []types.Section{{0, 4}},
		},
	}, context.query.QueryResults[2])
}

func TestMergeTreeNodesDeep(t *testing.T) {
	flag.Parse()

	q, _ := query.Parse("(a or (b and c))")
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 4,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 2,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, &[]types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 4}},
		},
	}, context.query.QueryResults[3])
	assert.Equal(t, &[]types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 4}},
		},
		{
			DocumentID: 4,
			Sections:   []types.Section{{0, 4}},
		},
	}, context.query.QueryResults[4])

}

func TestMergeTreeNodesPartial(t *testing.T) {
	flag.Parse()

	q, _ := query.Parse("a and b and (c or d)")
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
				nil,
				nil,
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, query.TreeQuery, q.SubQueries[1].Type)
	assert.Equal(t, 2, len(q.SubQueries))
	assert.Equal(t, &[]types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 4}},
		},
	}, context.query.QueryResults[0])

	q, _ = query.Parse("b and -a and (c or d)")
	context = Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
				nil,
				nil,
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, query.TreeQuery, q.SubQueries[1].Type)
	assert.Equal(t, 2, len(q.SubQueries))
	assert.Equal(t, &[]types.DocumentWithSections{
		{
			DocumentID: 3,
			Sections:   []types.Section{{0, 4}},
		},
	}, context.query.QueryResults[0])

	q, _ = query.Parse("-a and -b and (c or d)")
	context = Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
				nil,
				nil,
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, query.TreeQuery, q.SubQueries[1].Type)
	assert.Equal(t, 2, len(q.SubQueries))
	assert.Equal(t, &[]types.DocumentWithSections{
		{
			DocumentID: 1,
			Sections:   []types.Section{{0, 4}},
		},
		{
			DocumentID: 3,
			Sections:   []types.Section{{0, 4}},
		},
	}, context.query.QueryResults[0])
	assert.Equal(t, true, q.SubQueries[0].Negate)
}

func TestMergeTreeNodesShortcutAnd(t *testing.T) {
	flag.Parse()

	q, _ := query.Parse("a and b and (c or d)")
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{
					types.DocumentWithSections{
						DocumentID: 4,
						Sections:   []types.Section{{0, 4}},
					},
				},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
				nil,
				nil,
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, &[]types.DocumentWithSections{}, context.query.QueryResults[5])
}

func TestMergeTreeNodesShortcutOr(t *testing.T) {
	flag.Parse()

	q, _ := query.Parse("a or b or (c and d)")
	context := Context{
		query: &SearchQuery{
			QueryResults: []*[]types.DocumentWithSections{
				{},
				{
					types.DocumentWithSections{
						DocumentID: 1,
						Sections:   []types.Section{{0, 4}},
					},
					types.DocumentWithSections{
						DocumentID: 3,
						Sections:   []types.Section{{0, 4}},
					},
				},
				nil,
				nil,
				nil,
				nil,
			},
		},
	}
	internalMergeTreeNodes(&context, q)
	assert.Equal(t, 2, len(q.SubQueries))
}
