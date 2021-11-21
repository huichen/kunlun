package types

type SearchRequest struct {
	// 原始搜索表达式
	Query string

	// 如果输入文档 ID，则仅从这个文档范围查询，并返回全文
	DocumentIDs []uint64

	// 当返回多个仓库时，一个仓库最多返回多少文档，如果小于等于零则都返回
	MaxDocumentsPerRepo int
	// 当返回单个仓库时，该最多返回多少文档，如果小于等于零则都返回
	MaxDocumentsInSingleRepoReturn int

	// 当小于等于0时单文件返回所有结果，否则最多返回这么多行
	MaxLinesPerDocument int

	// 是否返回行内容
	ReturnLineContent bool

	// 返回的行内容包括目标行上下各这么多行的上下文，方便阅读
	// 如果该值小于等于0，不展示上下文
	NumContextLines int

	// 如果需要分页，该值设置为大于 0 的值，一页最多包含这么多仓库
	// 如果小于等于 0 则返回所有结果
	PageSize int
	// 和 PageSize 配合使用
	PageNum int

	// 高亮 tags
	HightlightStartTag string
	HightlightEndTag   string

	// 超时设置，单位毫秒
	TimeoutInMs int

	// 外部传入的钩子函数，用于判断某个 repo 里的文件是否应该被召回
	ShouldRecallRepo func(uint64) bool

	// 排序器，如果不传入则用默认的排序器
	Ranker Ranker
}
