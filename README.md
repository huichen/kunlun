昆仑代码搜索引擎
======

* 支持千万行级别代码仓库的实时查询
* 超快的索引速度：2000+ 万行代码 100 秒索引完毕
* 支持丰富的[搜索表达式类型](/doc/query/.md)
* 提供了 [KWS](/cmd/kws) 网页搜索界面
* 提供了 [KLS](/cmd/kls) 命令行界面的代码搜索瑞士军刀程序
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
```

引擎提供了丰富的遍历、索引和搜索[选项](/pkg/types/engine_options.go)。


# 其它

* [联系方式](/doc/feedback.md)
