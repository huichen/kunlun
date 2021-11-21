// KLS：Kunlun Search
// 命令行代码搜索工具
package main

import (
	"os"
	"os/signal"

	flag "github.com/spf13/pflag"

	"kunlun/internal/kls"
	"kunlun/pkg/log"
)

var (
	fileExtensionFilter = flag.String("ext", "", "只读取这些后缀的文件（半角逗号分隔）")
	ignoreDirs          = flag.String("ignore", ".git,.idea,vendor,target", "忽略这些文件夹下的文件（半角逗号分隔）")
	ctags               = flag.String("ctags", "", "CTags 二进制文件路径")
)

func main() {
	flag.Parse()
	log.SetLogger(&log.EmptyLogger{})

	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = append(dirs, ".")
	}

	options := kls.NewKLSOptions().
		SetFileExtensionFilter(*fileExtensionFilter).
		SetIgnoreDirs(*ignoreDirs).
		SetIndexDirs(dirs).
		SetCtagBinaryPath(*ctags)

	app := kls.NewKLS(options)

	// 捕获 ctrl-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			app.Stop()
		}
	}()

	app.Run()
}
