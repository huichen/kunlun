package common_types

// ctags 解析出来的文档中的符号信息
type CTagsEntry struct {
	Sym      string
	Path     string
	Line     int
	Kind     string
	Language string
}
