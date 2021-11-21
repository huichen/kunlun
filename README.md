昆仑代码搜索引擎
======

* 提供了 [KWS](/cmd/kws) 网页搜索界面
* 提供了 [KLS](/cmd/kls)，一个命令行界面的代码搜索瑞士军刀程序
* 采用对商业应用友好的[Apache License v2](/LICENSE)发布

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
	"log"

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
	kgn.Finish()                          // 开始搜索前必须先调用该函数

	// 构建索引
	kgn.IndexDir(*dir)

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

是不是很简单！


# 其它

* [联系方式](/doc/feedback.md)
