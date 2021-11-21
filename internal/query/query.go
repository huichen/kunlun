package query

import (
	"fmt"
	"strings"

	"github.com/huichen/kunlun/pkg/log"
)

var (
	logger = log.GetLogger()
)

type QueryType int

const (
	UnknownQueryType QueryType = iota
	TokenQuery                 // 简单文本查询
	RegexQuery                 // 正则表达式查询
	TreeQuery                  // 复杂的查询树
	// 修饰词
	FileQuery     // 文件名查询
	RepoQuery     // 仓库名查询
	CaseQuery     // 大小写区分
	LanguageQuery // 编程语言查询
)

// 查询表达式
// 保存了从表达式解析出来的树状结构
type Query struct {
	// 表达式类型
	Type QueryType

	// 在树状结构中的唯一标识 ID，不同的查询树重新构建 ID
	ID int

	// 当 QueryType 为简单文本查询时，存储需要匹配的字符串
	// 其他类型下，该字段为空
	Token string

	// 表达式基础上取反操作
	Negate bool

	// 当 QueryType 为正则表达式的时候，这里存储正则表达式原始字符串
	RegexString string
	// 当 QueryType 为正则表达式的时候，这里存储所有解析出来的 token
	RegexTokens []string

	// 当 QueryType 为查询树时，SubQueries 存储解析出的递归子查询表达式
	SubQueries []*Query
	Or         bool // 如果为 true，则 SubQueries 之间取或，否则取和

	// 当 token 或者 regex 类型的 symbol
	IsSymbol bool

	// FileQuery
	FileRegexString string

	// RepoQuery
	RepoRegexString string

	// CaseQuery
	Case bool

	// LanguageQuery
	Language string

	// 子节点的统计和调试信息
	NumNodes             int    // 树中一共多少节点（包含根节点）
	MaxDepth             int    // 子节点最大深度，叶子节点深度为 0
	RootDistance         int    // 该节点到 root 的距离
	stringRepresentation string // 该 query 的字符串表达形式
}

// 从 pattern 中解析 Query
func Parse(pattern string) (*Query, error) {
	q, err := internalParse(pattern)
	if err != nil {
		return nil, err
	}

	// 下放 NOT 操作符
	populateDownNegate(q, false)

	// 扁平化
	Flatten(q)

	// 重新排序
	Sort(q)

	return q, nil
}

func internalParse(pattern string) (*Query, error) {
	// 去掉首尾空格
	pattern = strings.TrimSpace(pattern)

	// 先做切分
	segments, err := segment(pattern)
	if err != nil {
		return nil, err
	}

	// 没有找到 segment 的情况
	if len(segments) == 0 {
		return nil, nil
	}

	// 只有一个 segment 的情况
	if len(segments) == 1 {
		// 先检查是否为特殊运算符
		lowerSeg := strings.ToLower(segments[0])
		if lowerSeg == "or" || lowerSeg == "and" || lowerSeg == "-" {
			return nil, fmt.Errorf("不能只有一个运算符操作 %s", lowerSeg)
		}

		seg := segments[0]
		if seg[0] == '"' {
			// 首尾都是 " 的情况
			seg = seg[1 : len(seg)-1]
		} else if seg[0] == '(' {
			// 被 () 包裹的情况
			return internalParse(seg[1 : len(seg)-1])
		} else if seg[0] == '-' {
			// 以 - 开头
			q, err := internalParse(seg[1:])
			if err != nil {
				return nil, err
			}
			q.Negate = !q.Negate
			return q, err
		}

		return parseTerm(seg)
	}

	// 多个 segments 的情况
	// 存储一系列 andQuery
	andQueries := [][]*Query{}
	currentAndQuery := []*Query{}

	for idx := 0; idx < len(segments); idx++ {
		seg := segments[idx]

		lowerSec := strings.ToLower(seg)
		if lowerSec == "or" {
			// 为 or 时先把之前所有的 and queries 添加到结果
			if len(currentAndQuery) > 0 {
				andQueries = append(andQueries, currentAndQuery)
			}
			currentAndQuery = []*Query{}
		} else if lowerSec != "and" {
			q, err := internalParse(seg)
			if err != nil {
				return nil, err
			}
			if q != nil {
				currentAndQuery = append(currentAndQuery, q)
			}
		}
	}

	// 最后一个 or 之后剩余的加入 andQueries
	if len(currentAndQuery) > 0 {
		andQueries = append(andQueries, currentAndQuery)
	}

	query := Query{
		SubQueries: []*Query{},
		Type:       TreeQuery,
	}

	// 将多个 and queries 合并
	if len(andQueries) == 0 {
		return nil, nil
	} else if len(andQueries) == 1 {
		query.SubQueries = andQueries[0]
	} else {
		// 多个 and queries 的情况
		for _, andQ := range andQueries {
			if len(andQ) != 0 {
				query.SubQueries = append(query.SubQueries, &Query{
					SubQueries: andQ,
					Type:       TreeQuery,
				})
			}
		}
		query.Or = true
	}

	// 检查是否为空
	if len(query.SubQueries) == 0 {
		return nil, nil
	}

	return &query, nil
}
