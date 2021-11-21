package indexer

import (
	"errors"

	"kunlun/pkg/types"
)

// 搜索得到的关键词 section 转化成容易阅读的带高亮行信息
func (indexer *Indexer) GetLinesFromSections(
	documentID uint64,
	sections []types.Section,
	contextLines int,
	maxLinesPerFile int,
) ([]types.Line, int, error) {

	var err error
	sections, err = mergeSections(sections)
	if err != nil {
		return nil, 0, err
	}

	lines := []types.Line{}
	var startLine, endLine, sizeB, e, b uint32
	for _, sec := range sections {
		// 首先得到该区间起始位置所属的行
		startLine, b, sizeB, err = indexer.getLine(documentID, sec.Start)
		if err != nil {
			return nil, 0, err
		}

		// 然后得到区间（闭区间）结束位置所属行
		possibleEndInSameLine := sec.End - sec.Start + b
		if possibleEndInSameLine <= sizeB {
			// 先简单判断是否在同一行，这是大多数情况
			endLine = startLine
			e = possibleEndInSameLine - 1
		} else {
			// 如果不在一行，搜结束位置所属行
			endLine, e, _, err = indexer.getLine(documentID, sec.End-1) // 最后一个在区间内的元素
			if err != nil {
				return nil, 0, err
			}
		}

		if startLine == endLine {
			// 一行内的情况，也是多数情况
			lines, err = appendLine(lines, startLine, types.Section{b, e + 1})
			if err != nil {
				return nil, 0, err
			}
		} else {
			lines, err = appendLine(lines, startLine, types.Section{b, sizeB})
			if err != nil {
				return nil, 0, err
			}

			for iLine := startLine + 1; iLine < endLine; iLine++ {
				lineSize := indexer.getLineSize(documentID, uint32(iLine))
				lines, err = appendLine(lines, iLine, types.Section{0, lineSize})
				if err != nil {
					return nil, 0, err
				}

			}
			lines, err = appendLine(lines, endLine, types.Section{0, e + 1})
			if err != nil {
				return nil, 0, err
			}

		}
	}
	totalResults := len(lines)
	if maxLinesPerFile > 0 && totalResults > maxLinesPerFile {
		lines = lines[:maxLinesPerFile]
	}

	// 添加上下文信息
	retLine := make([]types.Line, 0, len(lines)*(2*contextLines+1))
	maxLineNumber := indexer.getMaxLineNumber(documentID)
	curLine := -1
	for id, l := range lines {
		for ln := int(l.LineNumber) - contextLines; ln <= int(l.LineNumber)+contextLines; ln++ {
			if ln < 0 || ln > int(maxLineNumber) {
				// 行号越界
				continue
			}
			if ln <= curLine {
				// 本行已经输出
				continue
			}
			if id+1 < len(lines) && ln >= int(lines[id+1].LineNumber) {
				// 本行超出下一个待输出的匹配行的行号，不继续
				continue
			}
			curLine = ln

			if ln == int(l.LineNumber) {
				// 匹配行
				retLine = append(retLine, l)
			} else {
				// 匹配行的上下文
				retLine = append(retLine, types.Line{
					LineNumber: uint32(ln),
				})
			}
		}
	}
	return retLine, totalResults, nil
}

func appendLine(lines []types.Line, lineNumber uint32, highlight types.Section) ([]types.Line, error) {
	if len(lines) == 0 {
		lines = append(lines, types.Line{
			LineNumber: lineNumber,
			Highlights: []types.Section{highlight},
		})
		return lines, nil
	}

	oldLineNumber := lines[len(lines)-1].LineNumber
	if oldLineNumber > lineNumber {
		return nil, errors.New("appendLine 没有递增")
	} else if oldLineNumber == lineNumber {
		lines[len(lines)-1].Highlights = append(lines[len(lines)-1].Highlights, highlight)
	} else if oldLineNumber < lineNumber {
		lines = append(lines, types.Line{
			LineNumber: lineNumber,
			Highlights: []types.Section{highlight},
		})
	}

	return lines, nil
}

// 获得某个文档某行的字节大小
func (indexer *Indexer) getLineSize(docID uint64, line uint32) uint32 {
	indexer.indexerLock.RLock()
	defer indexer.indexerLock.RUnlock()

	return indexer.unsafeGetLineSize(docID, line)
}

// 获得某个文档的最大行号
func (indexer *Indexer) getMaxLineNumber(docID uint64) uint32 {
	indexer.indexerLock.RLock()
	defer indexer.indexerLock.RUnlock()

	lineStarts := indexer.documentIDToMetaMap[docID].LineStartLocations
	return uint32(len(lineStarts) - 1)
}

// 获得某个文档某行的内容
func (indexer *Indexer) GetLineContent(docID uint64, line uint32) ([]byte, uint32) {
	if !indexer.finished {
		logger.Fatal("必须先调用 Finish 函数")
	}

	meta := indexer.documentIDToMetaMap[docID]
	lineStarts := meta.LineStartLocations
	b := lineStarts[line]
	var e uint32
	if line == uint32(len(lineStarts)-1) {
		// 最后一行的情况
		e = uint32(meta.Size)
	} else {
		e = lineStarts[line+1]
	}

	return (*indexer.documentIDToContentMap[docID])[b:e], b
}

// 线程不安全版本，请不要直接调用
func (indexer *Indexer) unsafeGetLineSize(docID uint64, line uint32) uint32 {
	contentSize := indexer.documentIDToMetaMap[docID].Size
	lineStarts := indexer.documentIDToMetaMap[docID].LineStartLocations
	if int(line) == len(lineStarts)-1 {
		return uint32(contentSize) - lineStarts[line]
	}

	return lineStarts[line+1] - lineStarts[line]
}

// 二分法找到 loc 位置所属的行，行中位置和行长度，如果没找到则返回错误
func (indexer *Indexer) getLine(docID uint64, loc uint32) (uint32, uint32, uint32, error) {
	if !indexer.finished {
		logger.Fatal("必须先调用 Finish 函数")
	}

	if meta, ok := indexer.documentIDToMetaMap[docID]; !ok {
		return 0, 0, 0, errors.New("没找到文档")
	} else {
		if loc >= uint32(meta.Size) {
			return 0, 0, 0, errors.New("位置越界")
		}
	}

	lineStarts := indexer.documentIDToMetaMap[docID].LineStartLocations
	if len(lineStarts) == 0 {
		return 0, 0, 0, errors.New("文档无内容")
	}

	b := 0
	e := len(lineStarts) - 1
	if loc >= lineStarts[e] {
		// 属于最后一行，直接返回
		return uint32(e), loc - lineStarts[e], uint32(indexer.documentIDToMetaMap[docID].Size) - lineStarts[e], nil
	}
	var line int
	for e-b > 1 {
		m := (e + b) / 2
		ms := lineStarts[m]
		if ms < loc {
			b = m
		} else if ms == loc {
			b = m
			break
		} else {
			e = m
		}
	}
	line = b

	return uint32(line), loc - lineStarts[line], lineStarts[line+1] - lineStarts[line], nil
}
