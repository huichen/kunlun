package indexer

import (
	"testing"

	"github.com/huichen/kunlun/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestMergeSection(t *testing.T) {

	sections := []types.Section{
		{1, 2}, {4, 7}, {4, 8}, {6, 9}, {11, 12},
	}
	sections, _ = mergeSections(sections)
	assert.Equal(t, []types.Section{
		{1, 2}, {4, 9}, {11, 12},
	}, sections)

	sections = []types.Section{
		{1, 2}, {4, 7}, {4, 7}, {7, 9}, {11, 12},
	}
	sections, _ = mergeSections(sections)
	assert.Equal(t, []types.Section{
		{1, 2}, {4, 9}, {11, 12},
	}, sections)

	sections = []types.Section{
		{1, 2}, {4, 7}, {4, 7}, {7, 11}, {11, 12},
	}
	sections, _ = mergeSections(sections)
	assert.Equal(t, []types.Section{
		{1, 2}, {4, 12},
	}, sections)

	sections = []types.Section{
		{1, 2}, {4, 7}, {5, 6}, {7, 11}, {11, 12},
	}
	sections, _ = mergeSections(sections)
	assert.Equal(t, []types.Section{
		{1, 2}, {4, 12},
	}, sections)

	sections = []types.Section{
		{1, 2}, {4, 3}, {4, 7}, {7, 11}, {11, 12},
	}
	_, err := mergeSections(sections)
	logger.Error(err)
	assert.NotNil(t, err)

	sections = []types.Section{
		{1, 2}, {4, 6}, {4, 5}, {7, 11}, {11, 12},
	}
	_, err = mergeSections(sections)
	logger.Error(err)
	assert.NotNil(t, err)

}
