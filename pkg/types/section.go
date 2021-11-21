package types

type DocumentWithSections struct {
	DocumentID uint64

	Sections []Section
}

type Section struct {
	Start uint32
	End   uint32
}
