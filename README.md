昆仑代码搜索引擎
======

* 支持一亿行代码、多仓库、复杂表达式高速检索和实时查询
* 超快的建索引速度：基于 [Unicode 三元组（trigram）的内存索引](/doc/index.md)，每分钟可索引超过 1000 万行代码
* 超快的搜索速度：多数 query 百毫秒内返回，[优化了正则表达式查询](/doc/regexp.md)，和竞品相比有 [10 倍速度提升](/doc/benchmark.md)
* 支持丰富的[搜索语言](/doc/query.md)，可以基于正则表达式、与或非逻辑运算、文件名、仓库名等搜索
* 支持上百种编程语言的[检测和查询](/doc/language.md)
* 支持丰富的[索引文件过滤选项](/doc/index_filter.md)
* 支持基于 ctags 的[符号（变量、函数、类名等）查询](/doc/ctags.md)
* 支持可扩展的[访问权限控制](/doc/acl.md)
* 提供了 [KWS](/cmd/kws) 基于网页的搜索服务
* 提供了 [KLS](/cmd/kls) 命令行下的代码搜索瑞士军刀程序
* 采用对商业应用友好的 [Apache License v2](/LICENSE) 发布

# 安装/更新

```
go get -u -v github.com/huichen/kunlun
```

# 使用

先看一个例子（来自[cmd/examples/simplest_example.go](/cmd/examples/simplest_example.go)）

```go
package main

import (
	"flag"

	"github.com/huichen/kunlun/pkg/engine"
	"github.com/huichen/kunlun/pkg/types"
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
```

引擎提供了一系列[文件遍历选项](/pkg/types/walker_options.go)、[索引选项](/pkg/types/indexer_options.go)和[搜索选项](/pkg/types/searcher_options.go)，可以在引擎启动时通过[参数传入](/pkg/types/engine_options.go)。

如果你想阅读昆仑的代码，可以先看看[这篇文档](/doc/codebase.md)。

# 其它

* [为什么要有昆仑](/doc/why.md)
* [联系方式](/doc/feedback.md)
