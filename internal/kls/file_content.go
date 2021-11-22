package kls

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/huichen/kunlun/internal/ngram_index"
	"github.com/huichen/kunlun/pkg/types"
)

const (
	omitAfterChars = 300
)

// 重新绘制右侧的文件内容展示窗口
func (kls *KLS) redrawFileContent() {
	if kls.response == nil {
		return
	}

	kls.currentHighlightLineIndex = 0

	kls.fileContent.Clear()

	if kls.currentFileID != -1 {
		kls.displayFileContent()
		return
	}

	kls.fileName.Clear()
	kls.printFilename("[gold]搜索结果[white]%s\n", "")

	for _, repo := range kls.response.Repos {
		output := fmt.Sprintf("仓库：%s", repo.String())
		fmt.Fprintf(kls.fileContent, "[green]%s[white]\n", output)

		for _, doc := range repo.Documents {
			filename := doc.Filename
			output := fmt.Sprintf("文件路径：%s", filename)
			lang := doc.Language
			if lang != "" {
				output = fmt.Sprintf("语言：%s\n%s", lang, output)
			}

			fmt.Fprintf(kls.fileContent, "[green]%s[white]\n", output)
			prevLine := 0
			for _, line := range doc.Lines {

				if line.LineNumber > uint32(prevLine)+1 {
					fmt.Fprintf(kls.fileContent, "|\n")
				}
				content := strings.Trim(string(line.Content), "\n")
				shortened := ""
				if len(content) > omitAfterChars {
					content = content[:omitAfterChars]
					shortened = " [yellow]... (省略)[white]"
				}
				fmt.Fprintf(
					kls.fileContent,
					"[blue]%d[white]\t%s%s\n",
					line.LineNumber+1,
					getColoredString(content, line.Highlights, "white", "red"),
					shortened)
				prevLine = int(line.LineNumber)
			}
			fmt.Fprintf(kls.fileContent, "\n")
		}
	}
	kls.fileContent.ScrollToBeginning()
}

func (kls *KLS) displayFileContent() {
	if kls.currentFileID == -1 {
		return
	}

	doc := kls.response.Repos[0].Documents[kls.currentFileID]
	docID := doc.DocumentID
	content := string(kls.kgn.GetContent(docID))

	filename := doc.Filename
	output := fmt.Sprintf("文件路径：%s", filename)

	lang := doc.Language
	if lang != "" {
		output = fmt.Sprintf("语言：%s\n%s", lang, output)
	}

	kls.printFilename("[green]%s\n", output)
	lines := strings.Split(content, "\n")

	matchedLineIndex := 0
	lineIndex := 0
	firstMatchedLine := -1

	for matchedLineIndex < len(doc.Lines) && lineIndex < len(lines) {
		if lineIndex < int(doc.Lines[matchedLineIndex].LineNumber) {
			content := lines[lineIndex]
			shortened := ""
			if len(content) > omitAfterChars {
				content = content[:omitAfterChars]
				shortened = " [yellow]... (省略)[white]"
			}
			fmt.Fprintf(kls.fileContent, "[blue]%d[white]\t%s%s\n", lineIndex+1, content, shortened)
			lineIndex++
			continue
		}

		if len(doc.Lines[matchedLineIndex].Highlights) != 0 && firstMatchedLine == -1 {
			firstMatchedLine = lineIndex
		}

		content := lines[lineIndex]
		shortened := ""
		if len(content) > omitAfterChars {
			content = content[:omitAfterChars]
			shortened = " [yellow]... (省略)[white]"
		}
		fmt.Fprintf(kls.fileContent,
			"[blue]%d[white]\t%s%s\n",
			lineIndex+1,
			getColoredString(
				content,
				doc.Lines[matchedLineIndex].Highlights,
				"white",
				"red"),
			shortened)

		lineIndex++
		matchedLineIndex++
	}

	for lineIndex < len(lines) {
		fmt.Fprintf(kls.fileContent, "[blue]%d[white]\t%s\n", lineIndex+1, lines[lineIndex])
		lineIndex++
	}

	if firstMatchedLine >= 0 {
		toLine := int(firstMatchedLine - 10)
		if toLine < 0 {
			toLine = 0
		}
		kls.fileContent.ScrollTo(toLine, 0)
	} else {
		kls.fileContent.ScrollToBeginning()
	}
}

func getColoredString(
	str string, highlights []types.Section,
	contentColor string,
	highlightColor string,
) string {
	runes := ngram_index.DecodeRunes([]byte(str))
	retStr := ""
	idxLoc := 0
	iRune := uint32(0)
	iStr := uint32(0)
	for ; iRune < uint32(len(runes)) && idxLoc < len(highlights); iRune++ {
		if iStr < highlights[idxLoc].Start {
			retStr += string(runes[iRune])
			iStr += uint32(utf8.RuneLen(runes[iRune]))
		} else if iStr >= highlights[idxLoc].Start && iStr < highlights[idxLoc].End {
			retStr += fmt.Sprintf("[%s]%s[%s]", highlightColor, string(runes[iRune]), contentColor)
			iStr += uint32(utf8.RuneLen(runes[iRune]))
		} else {
			idxLoc++
			iRune--
		}
	}

	for ; iRune < uint32(len(runes)); iRune++ {
		retStr += string(runes[iRune])
	}
	return retStr
}

func (kls *KLS) printFilename(format string, message ...interface{}) {
	kls.fileName.Clear()
	fmt.Fprintf(kls.fileName, format, message...)
}
