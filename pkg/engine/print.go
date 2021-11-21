package engine

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"

	"github.com/huichen/kunlun/internal/ngram_index"

	"github.com/huichen/kunlun/pkg/types"
)

var (
	red   = color.New(color.FgHiMagenta).SprintfFunc()
	white = color.New(color.FgWhite).SprintfFunc()
	green = color.New(color.FgGreen).SprintfFunc()
)

func (engine *KunlunEngine) PrettyPrintSearchResponse(
	searchResponse *types.SearchResponse,
	colorPrint bool,
	printContent bool) {
	if searchResponse == nil || len(searchResponse.Repos) == 0 {
		return
	}
	if !colorPrint {
		color.NoColor = true
	} else {
		color.NoColor = false
	}
	for _, repo := range searchResponse.Repos {
		for _, doc := range repo.Documents {
			engine.prettyPrintDoc(doc.DocumentID, doc.Lines, printContent)
		}
	}
	if !printContent {
		fmt.Println()
	}

	// 打印搜索指标
	color.Green("检索到 %d 个仓库中的 %d 个结果\n", len(searchResponse.Repos), searchResponse.NumLines)
}

func (engine *KunlunEngine) prettyPrintDoc(documentID uint64, lines []types.Line, printContent bool) {
	meta := engine.indexer.GetMeta(documentID)

	// 打印文件名
	filename := getColoredString(meta.PathInRepo, nil, green, red)
	repoHost := ""
	if meta.Repo != nil {
		if meta.Repo.RemoteURL != "" {
			repoHost = meta.Repo.RemoteURL + ":"
		} else {
			repoHost = meta.Repo.LocalPath + ":"
		}
	} else {
		filename = getColoredString(meta.LocalPath, nil, green, red)
	}
	if printContent {
		fmt.Fprintf(color.Output, "%s%s\n", color.GreenString(repoHost), filename)
	}

	// 打印行
	if printContent {
		for _, line := range lines {
			content := string(line.Content)
			content = strings.Trim(content, "\n")
			lineContent := getColoredString(content, line.Highlights, white, red)
			fmt.Fprintf(color.Output, "%d\t%s\n", line.LineNumber+1, lineContent)
		}
	}

	if printContent {
		fmt.Print("\n")
	}
}

// 获得 str 的彩色显示
// 正常字体用 contentColor 颜色
// 高亮部分用 highlightColor 颜色
func getColoredString(
	str string, highlights []types.Section,
	contentColor func(format string, a ...interface{}) string,
	highlightColor func(format string, a ...interface{}) string,
) string {
	runes := ngram_index.DecodeRunes([]byte(str))
	retStr := ""
	idxLoc := 0
	iRune := uint32(0)
	iStr := uint32(0)
	for ; iRune < uint32(len(runes)) && idxLoc < len(highlights); iRune++ {
		if iStr < highlights[idxLoc].Start {
			retStr += contentColor(string(runes[iRune]))
			iStr += uint32(utf8.RuneLen(runes[iRune]))
		} else if iStr >= highlights[idxLoc].Start && iStr < highlights[idxLoc].End {
			retStr += highlightColor(string(runes[iRune]))
			iStr += uint32(utf8.RuneLen(runes[iRune]))
		} else {
			idxLoc++
			iRune--
		}
	}

	for ; iRune < uint32(len(runes)); iRune++ {
		retStr += contentColor(string(runes[iRune]))
	}
	return retStr
}
