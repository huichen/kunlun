package kls

import (
	"log"
	"os"
	"time"

	"kunlun/pkg/engine"
	"kunlun/pkg/types"
)

func (kls *KLS) buildIndex() {
	kls.printSuggestion("%s", "[yellow]索引创建中，请稍等 ...")

	// ctags
	ctagsOptions := types.NewCtagsParserOptions().
		SetBinaryPath(kls.options.ctagBinaryPath)

	// 文件遍历器选项
	walkerOptions := types.NewIndexWalkerOptions().
		SetAllowedFileExtensions(kls.options.fileExtensionFilter).
		SetIgnoreDirs(kls.options.ignoreDirs).
		SetCTagsParserOptions(ctagsOptions)

	// 索引器选项
	indexerOptions := types.NewIndexerOptions()

	// 搜索引擎
	engineOptions := types.NewEngineOptions().
		SetWalkerOptions(walkerOptions).
		SetIndexerOptions(indexerOptions)

	var err error
	kls.kgn, err = engine.NewKunlunEngine(engineOptions)
	if err != nil {
		log.Fatal(err)
	}
	kls.printSuggestion("索引创建中，请稍等...")

	go func() {
		oldMessage := ""
		for !kls.indexingFinished {
			stats := kls.kgn.GetWalkerStats()

			if oldMessage != stats.Message && stats.Message != "" {
				kls.printSuggestion("%s", stats.Message)
				oldMessage = stats.Message
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()

	start := time.Now()
	for _, loc := range kls.options.indexLocations {
		_, err := os.Stat(loc)
		if err == nil {
			kls.kgn.IndexDir(loc)
		} else {
			kls.kgn.IndexGitRemote(loc)
		}
	}
	kls.kgn.Finish()
	kls.indexingFinished = true

	errMsg := kls.kgn.GetWalkerStats().CurrentError
	if errMsg != "" {
		kls.printSuggestion("[red]索引出错：%s", errMsg)
	} else {
		kls.printSuggestion("[green]索引创建完毕，共索引了%d个文件，耗时%.3fs",
			kls.kgn.GetWalkerStats().IndexedFiles,
			float32(time.Since(start).Milliseconds())/1000.)
	}
}
