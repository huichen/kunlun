package main

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/handlers"

	"github.com/huichen/kunlun/internal/api"
	"github.com/huichen/kunlun/internal/resource/engine"
	"github.com/huichen/kunlun/internal/util"
	"github.com/huichen/kunlun/pkg/log"
)

var (
	repoFolders     = flag.String("repo_folders", "", "半角逗号分隔的本地git仓库地址")
	staticFolder    = flag.String("static_folder", "fe/dist", "静态文件目录")
	port            = flag.String("port", ":8080", "端口号")
	useDebugHandles = flag.Bool("use_debug_handles", false, "是否打开 debug http handles")
)

func main() {
	flag.Parse()

	// 初始化 kunlun engine
	engine.Init()

	// 启动索引创建线程
	go buildIndex()

	// 捕获ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.GetLogger().Info("捕获Ctrl-c，退出服务器")
			os.Exit(0)
		}
	}()

	// API 服务路由
	m := http.NewServeMux()

	// 服务器状态
	m.HandleFunc("/healthz", api.Healthz)

	// 服务器状态
	m.HandleFunc("/api/search", api.Search)

	// 静态页面
	m.Handle("/", http.FileServer(http.Dir(*staticFolder)))

	// 开启调试接口
	if *useDebugHandles {
		m.HandleFunc("/debug/pprof/", pprof.Index)
		m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		m.HandleFunc("/debug/pprof/profile", pprof.Profile)
		m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		m.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// 启动服务
	server := &http.Server{
		Addr:    *port,
		Handler: handlers.CompressHandler(m),
	}
	log.GetLogger().Info("服务器启动")
	log.GetLogger().Fatal(server.ListenAndServe())
}

func buildIndex() {
	log.GetLogger().Info("索引开始")
	kgn := engine.GetEngine()
	fields := strings.Split(*repoFolders, ",")

	startTime := time.Now()
	for _, repo := range fields {
		kgn.IndexDir(repo)
	}
	kgn.Finish()

	loc := kgn.GetWalkerStats().TotalLinesOfCode
	ms := time.Since(startTime).Milliseconds()
	log.GetLogger().Infof("索引完毕，耗时 %.3f 秒，共索引 %d 行代码，平均每秒索引 %d 行代码", float32(ms)/1000.0, loc, loc*1000/int(ms))

	// 打印遍历统计指标
	util.PrintWalkerStats(kgn)

	// 打印索引统计指标
	util.PrintIndexerStats(kgn)
}
