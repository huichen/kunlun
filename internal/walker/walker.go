package walker

import (
	"runtime"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/huichen/kunlun/internal/ctags"
	"github.com/huichen/kunlun/pkg/log"
	"github.com/huichen/kunlun/pkg/types"
)

var (
	logger = log.GetLogger()
)

// 文件遍历器，用于控制哪些仓库和文件需要添加到索引中
type IndexWalker struct {
	options *types.IndexWalkerOptions

	// 仓库
	repos     map[string]bool
	reposLock sync.RWMutex

	// 发送给下游处理
	fileChan chan WalkerFileInfo

	stats     types.IndexWalkerStats
	statsLock sync.RWMutex

	ctagsParsers []*ctags.CTagsParser

	walkerInfoChan chan walkerInfo

	// git public keys
	pubKeys *ssh.PublicKeys

	pullMode bool
}

func NewIndexWalker(options *types.IndexWalkerOptions, pullMode bool) (*IndexWalker, error) {
	if options == nil {
		options = types.NewIndexWalkerOptions()
	}

	wkr := &IndexWalker{
		pullMode:       pullMode,
		options:        options,
		fileChan:       make(chan WalkerFileInfo, runtime.NumCPU()*2),
		walkerInfoChan: make(chan walkerInfo, runtime.NumCPU()*2),
		repos:          make(map[string]bool),
	}
	wkr.stats.Languages = make(map[string]types.FilesLinesBytes)

	for i := 0; i < options.NumFileProcessors; i++ {
		go wkr.fileProcessor(i)
		parser, err := ctags.NewCTagsParser(options.CTagsParserOptions)
		if err != nil {
			return nil, err
		}
		if parser != nil {
			wkr.ctagsParsers = append(wkr.ctagsParsers, parser)
		}
	}

	return wkr, nil
}

func (dw *IndexWalker) GetFileChan() chan WalkerFileInfo {
	return dw.fileChan
}

func (dw *IndexWalker) Finish() {
	for _, parser := range dw.ctagsParsers {
		parser.Close()
	}

	for i := 0; i < dw.options.NumFileProcessors; i++ {
		dw.walkerInfoChan <- walkerInfo{exit: true}
	}
}
