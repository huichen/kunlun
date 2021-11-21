package main

import (
	"flag"

	"kunlun/pkg/engine"
	"kunlun/pkg/types"
)

var (
	dir   = flag.String("d", "/usr/local/include", "索引这个文件夹下的所有文件")
	query = flag.String("q", "gcc", "搜索表达式")
)

func main() {
	flag.Parse()

	// 创建引擎
	kgn, _ := engine.NewKunlunEngine(nil) // 使用默认选项

	// 构建索引
	kgn.IndexDir(*dir)
	kgn.Finish() // 开始搜索前必须先调用该函数

	// 检索
	request := types.SearchRequest{
		Query:             *query,
		ReturnLineContent: true,
		NumContextLines:   2}
	resp, _ := kgn.Search(request)

	// 打印输出
	kgn.PrettyPrintSearchResponse(resp, true, true)
}
