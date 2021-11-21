package query

import "fmt"

// 通过分隔符拆成多个 segments
// 返回解析出来的分段，一个分段可以是 token，也可以是一个表达式（需要进一步解析）
func segment(pattern string) ([]string, error) {
	segments := []string{}

	// 遍历时存储当前分段的起始位置
	var segmentStart int

	// 遍历时标记当前游标是否在一个分段中
	inSegment := false

	for idx := 0; idx < len(pattern); idx++ {
		char := pattern[idx]
		if char == ' ' || char == '\t' || char == '\n' {
			// 这三个字符当做间隔符处理
			if inSegment {
				segments = append(segments, pattern[segmentStart:idx])
				inSegment = false
			}
		} else if char == '"' {
			hasMinusPrefix := (idx > 0 && pattern[idx-1] == '-')
			hasColonPrefix := (idx > 0 && pattern[idx-1] == ':')

			// 处理遗留 segment
			if !hasMinusPrefix && !hasColonPrefix && inSegment {
				segments = append(segments, pattern[segmentStart:idx])
			}

			if !hasMinusPrefix && !hasColonPrefix {
				segmentStart = idx
			}

			// 找到对应的 '"'
			var segmentEnd int
			idx++
			for idx < len(pattern) {
				if pattern[idx] == '"' {
					if pattern[idx-1] != '\\' {
						segmentEnd = idx + 1
						break
					}
				}
				idx++
			}

			if segmentEnd > segmentStart {
				segments = append(segments, pattern[segmentStart:segmentEnd])
				inSegment = false
			} else {
				return nil, fmt.Errorf("query 解析失败，位置 %d：'\"' 没有找到对应的 '\"'", segmentStart)
			}
		} else if char == '(' {
			hasMinusPrefix := (idx > 0 && pattern[idx-1] == '-')

			// 处理遗留 segment
			if !hasMinusPrefix && inSegment {
				segments = append(segments, pattern[segmentStart:idx])
			}

			if !hasMinusPrefix {
				segmentStart = idx
			}

			// 找到对应的 ')'
			parenthesesLevel := 0
			inQuote := false
			var segmentEnd int
			idx++
			for idx < len(pattern) {
				if pattern[idx] == '(' && !inQuote {
					parenthesesLevel++
				} else if pattern[idx] == ')' && !inQuote {
					if parenthesesLevel == 0 {
						segmentEnd = idx + 1
						break
					} else {
						parenthesesLevel--
					}
				} else if pattern[idx] == '"' && (idx == 0 || pattern[idx-1] != '\\') {
					inQuote = !inQuote
				}
				idx++
			}

			// 如果 () 之间的部分不为空，则添加
			if segmentEnd != 0 {
				segments = append(segments, pattern[segmentStart:segmentEnd])
				inSegment = false
			} else {
				return nil, fmt.Errorf("query 解析失败，位置 %d：'(' 没有找到对应的 ')'", segmentStart)
			}
		} else if char == ')' {
			return nil, fmt.Errorf("query 解析失败，位置 %d：')' 没有找到对应的 '('", idx)
		} else {
			if !inSegment {
				if char == '-' {
					if idx > 0 && pattern[idx-1] == '-' {
						return nil, fmt.Errorf("- 不能连用")
					}
				}
				inSegment = true
				segmentStart = idx
			}
		}
	}

	// 添加最后一个剩余的 segment
	if inSegment {
		segments = append(segments, pattern[segmentStart:])
	}

	return segments, nil
}
