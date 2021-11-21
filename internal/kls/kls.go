package kls

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/huichen/kunlun/pkg/engine"
	"github.com/huichen/kunlun/pkg/log"
	"github.com/huichen/kunlun/pkg/types"
)

var (
	logger = log.GetLogger()
)

type KLS struct {
	options *KLSOptions

	// 昆仑引擎和搜索返回结果
	kgn      *engine.KunlunEngine
	response *types.SearchResponse

	// 主应用
	app *tview.Application

	// 组件
	inputField     *tview.InputField
	suggestionText *tview.TextView
	repoList       *tview.List
	fileList       *tview.List
	fileContent    *tview.TextView
	fileName       *tview.TextView

	// 文件列表
	currentRepoID             int
	currentFileID             int
	currentHighlightLineIndex int

	// 状态
	indexingFinished bool
	inSearching      bool
}

func NewKLS(options *KLSOptions) *KLS {
	kls := &KLS{
		options: options,
	}

	// 主应用
	kls.app = tview.NewApplication().EnableMouse(true)

	// 搜索框
	kls.inputField = tview.NewInputField().
		SetLabel("[gold]搜索：").
		SetFieldWidth(60).
		SetDoneFunc(func(key tcell.Key) {
			inputText := kls.inputField.GetText()
			if inputText == "" {
				return
			} else if inputText == "q" {
				kls.app.Stop()
			}

			if !kls.indexingFinished {
				return
			}
			go kls.search()
		})

	// 搜索框右侧提示框
	kls.suggestionText = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	// 搜索结果列表
	kls.repoList = tview.NewList().ShowSecondaryText(false)
	kls.repoList.SetTitle("[gold]仓库列表")
	kls.repoList.SetBorder(true)
	kls.repoList.SetBorderPadding(0, 0, 1, 1)
	kls.fileList = tview.NewList().ShowSecondaryText(false)
	kls.fileList.SetTitle("[gold]文件列表")
	kls.fileList.SetBorder(true)
	kls.fileList.SetBorderPadding(0, 0, 1, 1)

	leftPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(kls.repoList, 0, 2, true).
		AddItem(kls.fileList, 0, 8, true)

	// 文件内容
	kls.fileName = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true)
	kls.fileContent = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true)
	filePanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(kls.fileName, 2, 1, false).
		AddItem(kls.fileContent, 0, 1, true)
	filePanel.SetBorder(true)
	filePanel.SetBorderPadding(0, 0, 1, 1)
	filePanel.SetTitle("[gold] 文件内容")

	// 页头
	header := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(kls.inputField, 40, 0, true).
		AddItem(nil, 2, 0, false).
		AddItem(kls.suggestionText, 0, 1, false)

	// 页中
	body := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(leftPanel, 0, 3, true).
		AddItem(nil, 1, 0, false).
		AddItem(filePanel, 0, 7, false)

	// 主窗口
	mainWindow := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 1, 1, true).
		AddItem(body, 0, 1, false)

	// 设置按键
	kls.SetInputFieldKeys()
	kls.SetRepoListKeys()
	kls.SetFileListKeys()
	kls.SetFileContentKey()

	kls.app.SetRoot(mainWindow, true)

	return kls
}

func (kls *KLS) Run() {
	go kls.buildIndex()

	if err := kls.app.Run(); err != nil {
		logger.Error(err)
		return
	}
}

func (kls *KLS) Stop() {
	kls.app.Stop()
}
