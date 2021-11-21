package engine

import (
	"flag"

	"github.com/huichen/kunlun/pkg/engine"
	"github.com/huichen/kunlun/pkg/log"
	"github.com/huichen/kunlun/pkg/types"
)

var (
	kgn *engine.KunlunEngine

	logger = log.GetLogger()
)

var (
	allowedExtensions    = flag.String("ext", "", "只读取这些后缀的文件（半角逗号分隔）")
	allowedLanguages     = flag.String("lang", "", "只读取这些编程语言的文件（半角逗号分隔）")
	ignoreDirs           = flag.String("ignore", ".git,vendor,target", "忽略这些文件夹下的文件（半角逗号分隔）")
	ctags                = flag.String("ctags", "", "ctags 二进制文件地址")
	maxFileSize          = flag.Int("size", 0, "最大文件尺寸")
	indexShards          = flag.Int("shards", 0, "检索器分片，设置为 0 则使用全部 CPU")
	disallowedLanguages  = flag.String("dis_langs", "svg", "语言黑名单，半角逗号分隔")
	disallowedExtensions = flag.String("dis_exts", "svg", "后缀黑名单，半角逗号分隔")
)

func Init() {
	// ctags 选项
	ctagsOptions := types.NewCtagsParserOptions().SetBinaryPath(*ctags)

	// 文件遍历器选项
	walkerOptions := types.NewIndexWalkerOptions().
		SetAllowedFileExtensions(*allowedExtensions).
		SetAllowedCodeLanguages(*allowedLanguages).
		SetMaxFileSize(*maxFileSize).
		SetIgnoreDirs(*ignoreDirs).
		SetCTagsParserOptions(ctagsOptions).
		SetDisallowedCodeLanguages(*disallowedLanguages).
		SetDisallowedFileExtensions(*disallowedExtensions)

	// 索引器选项
	indexerOptions := types.NewIndexerOptions().
		SetNumIndexerShards(*indexShards)

	// 检索器选项
	searcherOptions := types.NewSearcherOptions()

	// 搜索引擎
	engineOptions := types.NewEngineOptions().
		SetIndexerOptions(indexerOptions).
		SetWalkerOptions(walkerOptions).
		SetSearcherOptions(searcherOptions)

	var err error
	kgn, err = engine.NewKunlunEngine(engineOptions)
	if err != nil {
		logger.Fatal(err)
	}

}

func GetEngine() *engine.KunlunEngine {
	if kgn == nil {
		logger.Fatal("GetEngine 之前必须调用 Init")
	}
	return kgn
}
