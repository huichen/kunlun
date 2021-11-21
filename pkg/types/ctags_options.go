package types

const (
	defaultCtagsParsingTimeout = 10000
)

type CTagsParserOptions struct {
	BinaryPath string
}

func NewCtagsParserOptions() *CTagsParserOptions {
	options := &CTagsParserOptions{}
	return options
}

func (cp *CTagsParserOptions) SetBinaryPath(path string) *CTagsParserOptions {
	cp.BinaryPath = path
	return cp
}
