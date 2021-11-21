package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"kunlun/internal/searcher"
	"kunlun/internal/util"
	"kunlun/pkg/engine"
	"kunlun/pkg/log"
	"kunlun/pkg/types"
)

var (
	logger = log.GetLogger()

	// 索引参数
	git               = flag.String("git", "", "远程 git 仓库")
	dir               = flag.String("dir", "", "索引这个文件夹下的所有文件")
	allowedExtensions = flag.String("ext", "", "只读取这些后缀的文件（半角逗号分隔）")
	allowedLanguages  = flag.String("lang", "", "只读取这些编程语言的文件（半角逗号分隔）")
	ignoreDirs        = flag.String("ignore", ".git,vendor,target", "忽略这些文件夹下的文件（半角逗号分隔）")
	ctags             = flag.String("ctags", "", "ctags 二进制文件地址")
	maxFileSize       = flag.Int("size", 0, "最大文件尺寸")
	indexShards       = flag.Int("shards", 0, "检索器分片，设置为 0 则使用全部 CPU")

	// 搜索参数
	query = flag.String("q", "gcc", "待搜索的关键词")

	// 显示参数
	print = flag.Bool("print", false, "在结果中打印文件内容，否则只打印文件名")
	lines = flag.Int("lines", 0, "单文件最多返回多少结果")

	// 其他参数
	searchThreads = flag.Int("threads", 1, "搜索并发数")
	searchRuns    = flag.Int("runs", 100, "每个并发搜索多少次")
)

func main() {
	flag.Parse()
	if *git == "" && *dir == "" {
		logger.Fatal("-git 和 -dir 参数不能同时为空")
	}

	icpuf, err := os.Create("indexer.cpu_profile")
	if err != nil {
		logger.Fatal(err)
	}
	imemf, err := os.Create("indexer.mem_profile")
	if err != nil {
		logger.Fatal(err)
	}
	scpuf, err := os.Create("searcher.cpu_profile")
	if err != nil {
		logger.Fatal(err)
	}
	smemf, err := os.Create("searcher.mem_profile")
	if err != nil {
		logger.Fatal(err)
	}

	// ctags 选项
	ctagsOptions := types.NewCtagsParserOptions().SetBinaryPath(*ctags)

	// 文件遍历器选项
	walkerOptions := types.NewIndexWalkerOptions().
		SetAllowedFileExtensions(*allowedExtensions).
		SetAllowedCodeLanguages(*allowedLanguages).
		SetMaxFileSize(*maxFileSize).
		SetIgnoreDirs(*ignoreDirs).
		SetCTagsParserOptions(ctagsOptions)

	// 索引器选项
	indexerOptions := types.NewIndexerOptions().
		SetNumIndexerShards(*indexShards)

	// 搜索引擎
	engineOptions := types.NewEngineOptions().
		SetIndexerOptions(indexerOptions).
		SetWalkerOptions(walkerOptions)
	kgn, err := engine.NewKunlunEngine(engineOptions)
	if err != nil {
		logger.Fatal(err)
	}

	// 添加索引
	startTime1 := time.Now()
	pprof.StartCPUProfile(icpuf)
	if *dir != "" {
		kgn.IndexDir(*dir)
	}
	if *git != "" {
		kgn.IndexGitRemote(*git)
	}
	kgn.Finish()
	pprof.StopCPUProfile()
	endTime1 := time.Now()
	pprof.WriteHeapProfile(imemf)

	// 打印遍历统计指标
	util.PrintWalkerStats(kgn)

	// 打印索引统计指标
	util.PrintIndexerStats(kgn)

	// 生成查询请求
	request := types.SearchRequest{
		Query:               *query,
		ReturnLineContent:   *print,
		NumContextLines:     2,
		MaxLinesPerDocument: *lines,
	}

	// 启动 search worker
	startChan := make(chan bool, *searchThreads)
	retChan := make(chan float32, *searchThreads)
	for i := 0; i < *searchThreads; i++ {
		go searchWorker(kgn, request, startChan, retChan)
	}

	// 并发查询
	startTime2 := time.Now()
	pprof.StartCPUProfile(scpuf)
	totalLatency := float32(0)
	for i := 0; i < *searchThreads; i++ {
		startChan <- true
	}
	for i := 0; i < *searchThreads; i++ {
		totalLatency += (<-retChan)
	}
	pprof.StopCPUProfile()
	endTime2 := time.Now()
	pprof.WriteHeapProfile(smemf)

	// 单次查询并打印结果
	resp, err := kgn.Search(request)
	if err != nil {
		logger.Fatal(err)
	}
	kgn.PrettyPrintSearchResponse(resp, true, *print)

	// 打印搜索指标
	q, _ := searcher.ParseQuery(*query)
	fmt.Println()
	fmt.Printf("搜索表达式：%s\n", q.TrimmedQuery)
	if q.LanguageQuery != nil {
		fmt.Printf("lang 修饰词：%s\n", q.LanguageQuery)
	}
	if q.RepoQuery != nil {
		fmt.Printf("repo 修饰词：%s\n", q.RepoQuery)
	}
	if q.FileQuery != nil {
		fmt.Printf("file 修饰词：%s\n", q.FileQuery)
	}
	if q.Case {
		fmt.Printf("case:yes")
	}
	fmt.Printf("索引耗时 %.3f s\n", float32(endTime1.Sub(startTime1).Milliseconds())/1000)
	fmt.Printf("做了 %d 次正则表达式匹配\n", resp.NumRegexMatches)
	fmt.Printf("使用 %d 个协程并发搜索，每个协程请求 %d 次\n", *searchThreads, *searchRuns)
	fmt.Printf("搜索 QPS %.2f，平均延时 %.3f ms\n",
		float32(*searchThreads**searchRuns)/float32(endTime2.Sub(startTime2).Microseconds())*1000,
		totalLatency/(float32(*searchThreads))/1000)
}

func searchWorker(kgn *engine.KunlunEngine, request types.SearchRequest, startChan chan bool, retChan chan float32) {
	<-startChan

	startTime := time.Now()
	for i := 0; i < *searchRuns; i++ {
		_, err := kgn.Search(request)

		if err != nil {
			logger.Fatal(err)
		}
	}
	endTime := time.Now()

	retChan <- float32(endTime.Sub(startTime).Microseconds()) / float32(*searchRuns)
}
