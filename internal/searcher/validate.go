package searcher

import (
	"errors"

	"github.com/huichen/kunlun/internal/query"
)

// 对解析得到的 query 做合法性校验
func validateQuery(q *query.Query) error {
	// - 和 OR 的关系校验
	err := validateNegateOrRelation(q)
	if err != nil {
		return err
	}

	return nil
}

func validateModifiers(q *query.Query) error {
	if q == nil {
		return nil
	}

	ca, err := internalValidateModifiers(q)
	if err != nil {
		return err
	}

	if ca > 1 {
		return errors.New("ca: 只能有一个")
	}
	return nil
}

// 返回参数：
// case: 个数
func internalValidateModifiers(q *query.Query) (int, error) {
	switch q.Type {
	case query.FileQuery:
		if q.RootDistance > 2 {
			return 0, errors.New("file: 修饰词不能超过两层")
		}
		return 0, nil
	case query.RepoQuery:
		if q.RootDistance > 2 {
			return 0, errors.New("repo: 修饰词不能超过两层")
		}
		return 0, nil
	case query.LanguageQuery:
		if q.RootDistance > 2 {
			return 0, errors.New("lang: 修饰词不能超过两层")
		}
		return 0, nil
	case query.CaseQuery:
		if q.RootDistance != 1 {
			return 0, errors.New("case: 修饰词必须在第一层")
		}
		if q.Negate {
			return 0, errors.New("case: 修饰词前不能有 NOT")
		}
		return 1, nil
	case query.TreeQuery:
		var ca int
		var numFile, numRepo, numLang int
		for _, sq := range q.SubQueries {
			if sq == nil {
				continue
			}
			switch sq.Type {
			case query.FileQuery:
				numFile++
			case query.RepoQuery:
				numRepo++
			case query.LanguageQuery:
				numLang++
			}
			c, err := internalValidateModifiers(sq)
			if err != nil {
				return 0, err
			}

			ca += c
		}

		if q.RootDistance > 0 {
			if numFile != len(q.SubQueries) && numFile != 0 {
				return 0, errors.New("file: 不能和其他操作符在二级 query 中并列")
			}
			if numRepo != len(q.SubQueries) && numRepo != 0 {
				return 0, errors.New("repo: 不能和其他操作符在二级 query 中并列")
			}
			if numLang != len(q.SubQueries) && numLang != 0 {
				return 0, errors.New("lang: 不能和其他操作符在二级 query 中并列")
			}
		}

		return ca, nil
	default:
	}
	return 0, nil
}

// 如果 - 和 or 一起出现，就报错，比如下面的例子是非法的
// 		a or -b
//		a and (b or -c)
// 我们禁止这样的表达式出现主要有两个考虑
// 1、“非”操作通常会召回很多文档，再与另一些召回的文档做“或”操作，这会召回更多文档，
//	  而且召回的文档通常没有特别好的解释含义
// 2、“非”和“或”的结合，对做数组归并带来实际困难，而“非”和“与”操作结合做数组归并容易很多
func validateNegateOrRelation(q *query.Query) error {
	if q == nil {
		return nil
	}

	if q.Type == query.TreeQuery {
		if q.Or {
			for _, sq := range q.SubQueries {
				if sq != nil && sq.Negate {
					return errors.New("'-' 操作符禁止在 OR 连接的表达式中直接（但可以嵌套使用）")
				}
			}
		}
		for _, sq := range q.SubQueries {
			err := validateNegateOrRelation(sq)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
