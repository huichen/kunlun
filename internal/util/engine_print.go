package util

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kr/pretty"

	"github.com/huichen/kunlun/pkg/engine"
)

func PrintIndexerStats(kgn *engine.KunlunEngine) {
	stats := kgn.GetIndexerStats()
	fmt.Printf("索引统计：\n"+
		"索引分片数\t\t%d\n"+
		"总索引大小\t\t%d MB\n"+
		"总内容大小\t\t%d MB\n"+
		"成功索引文档数\t\t%d\n"+
		"索引出错文档数\t\t%d\n",
		stats.IndexerShards,
		int(float32(stats.TotalIndexSize)/1024.0/1024.0),
		int(float32(stats.TotalContentSize)/1024.0/1024.0),
		stats.TotalDocumentCount,
		stats.FailedAddingSymbol)
}

func PrintWalkerStats(kgn *engine.KunlunEngine) {
	stats := kgn.GetWalkerStats()
	stats.Languages = nil
	stats.Message = ""
	stats.CurrentFile = ""
	fmt.Printf("遍历统计：\n%# v\n", pretty.Formatter(stats))

	stats = kgn.GetWalkerStats()
	var langStats []LanguageWithStats
	for l, s := range stats.Languages {
		langStats = append(langStats, LanguageWithStats{
			Lang:  l,
			Files: s.NumFiles,
			Lines: s.NumLines,
			Bytes: s.NumBytes,
		})
	}

	sort.Slice(langStats, func(i, j int) bool {
		l1 := langStats[i].Lines
		l2 := langStats[j].Lines
		if l1 != l2 {
			return l1 > l2
		}

		return strings.Compare(langStats[i].Lang, langStats[j].Lang) < 0
	})

	fmt.Printf("语言统计（过滤后实际索引的部分）：\n%28s\t文件数\t行数\t字节数\n", "语言")
	totalFiles := 0
	totalLines := 0
	totalBytes := 0
	for _, lang := range langStats {
		fmt.Printf("%28s\t%d\t%d\t%d\n", lang.Lang, lang.Files, lang.Lines, lang.Bytes)
		totalFiles += lang.Files
		totalLines += lang.Lines
		totalBytes += lang.Bytes
	}
	fmt.Printf("%28s\t%d\t%d\t%d\n", "总计", totalFiles, totalLines, totalBytes)

}

type LanguageWithStats struct {
	Lang  string
	Files int
	Lines int
	Bytes int
}
