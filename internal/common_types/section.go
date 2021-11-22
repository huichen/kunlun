package common_types

import "github.com/huichen/kunlun/pkg/types"

type DocumentWithSections struct {
	DocumentID uint64

	Sections []types.Section
}
