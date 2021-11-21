昆仑代码搜索命令行工具（Kunlun commandLine Search）
=======

KLS 提供了一个命令行界面的代码搜索瑞士军刀程序，方便你在本地终端或者 SSH 中快速搜索代码

![](https://github.com/huichen/kunlun/blob/master/doc/kls.png)

#### 如何使用

进入 cmd/ls 目录，输入

```
go run main.go <你的代码目录列表>
```

其中代码目录列表用空格分隔，比如

```
go run main.go /usr/local/include ~/mygitrepo
```

进入文本界面后输入搜索表达式回车，然后用 tab 键在各个窗口跳转，其中“文件内容”窗口支持 VIM 模式的滚动和翻页，按 "/" 进入搜索框。在某些 terminal 中也支持鼠标点击。

#### 静态编译、交叉编译

你可以使用下面的命令将 KLS 编译成一个没有依赖的独立可执行程序，然后 scp 到服务器上去就可以在 ssh 的命令行界面中做代码搜索了

```
GOOS=linux go build -ldflags="-extldflags=-static"
```