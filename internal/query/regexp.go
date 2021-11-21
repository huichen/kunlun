package query

import (
	"regexp/syntax"
)

// 解析正则表达式，尽可能得到所有可以做文本匹配的子字符串
func parseRegexp(r *syntax.Regexp) []string {
	// 基础情况：没有子正则表达式
	if len(r.Sub) == 0 {
		switch r.Op {
		case syntax.OpLiteral:
			// literal string 的情况，返回原始字符串
			token := string(r.Rune)
			if token == "" {
				return nil
			}
			return []string{token}
		case syntax.OpCapture:
		default:
			// 其他情况返回空
			return nil
		}
	}

	// 有子表达式的情况，先删除一些可选的情况
	switch r.Op {
	case syntax.OpStar:
		return nil
	case syntax.OpQuest:
		return nil
	case syntax.OpRepeat:
		if r.Min <= 0 {
			return nil
		}

		// 对字符重复的特殊情况特殊处理
		if r.Min <= r.Max {
			if len(r.Sub) == 1 && r.Sub[0].Op == syntax.OpLiteral {
				token := string(r.Sub[0].Rune)
				if len(token) != 0 {
					repeatedToken := ""
					for i := 0; i < r.Min; i++ {
						repeatedToken += token
					}
					return []string{repeatedToken}
				}
			}
		}
	}

	// 递归得到所有子 token，并添加
	subTokens := []string{}
	for _, s := range r.Sub {
		tokens := parseRegexp(s)
		if tokens != nil {
			subTokens = append(subTokens, tokens...)
		}
	}

	if len(subTokens) == 0 {
		return nil
	}

	return subTokens
}
