package walker

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"

	"github.com/huichen/kunlun/pkg/types"
)

type walkerInfo struct {
	path string

	// 如果为内存读取的远端代码仓库，输入下面的值
	isGitRemoteRepo bool
	content         []byte
	repoRemoteURL   string

	// 退出信号
	exit bool
}

func (dw *IndexWalker) fileProcessor(shard int) {
	for {
		info := <-dw.walkerInfoChan
		if info.exit {
			return
		}
		path := info.path

		// 文件后缀过滤
		fileExtension := strings.TrimPrefix(filepath.Ext(path), ".")
		if len(dw.options.DisallowedFileExtensions) > 0 {
			if _, ok := dw.options.DisallowedFileExtensions[fileExtension]; ok {
				continue
			}
		}
		if len(dw.options.AllowedFileExtensions) > 0 {
			if _, ok := dw.options.AllowedFileExtensions[fileExtension]; !ok {
				// 当文件后缀过滤已经设置，并且没有命中，则退出
				dw.stats.FilteredByExtension++
				continue
			}
		}

		// . 文件
		if dw.options.FilterDotPrefix && strings.HasPrefix(strings.TrimSpace(filepath.Base(path)), ".") {
			dw.stats.FilteredByDotPrefix++
			continue
		}

		// 读取文件内容
		var content []byte
		var err error
		if info.content == nil {
			content, err = os.ReadFile(path)
			if err != nil {
				dw.stats.FilteredByError++
				dw.fileChan <- WalkerFileInfo{
					Error: err,
				}
				continue
			}
		} else {
			content = info.content
		}

		// 文件尺寸过滤
		if len(content) > dw.options.MaxFileSize {
			dw.stats.FilteredBySize++
			continue
		}

		// 是否是二进制文件
		if checkIfBinaryFile(content) {
			dw.stats.FilteredByBinaryType++
			continue
		}

		// 图片
		if enry.IsImage(path) {
			dw.stats.FilteredByImage++
			continue
		}

		// vendor
		if enry.IsGenerated(path, content) {
			dw.stats.FilteredByVendor++
			continue
		}

		// 生成的文件
		if enry.IsGenerated(path, content) {
			dw.stats.FilteredByGenerated++
			continue
		}

		// 过滤行比较多的文件
		fileLines := numLines(content)
		if fileLines > dw.options.MaxFileLines {
			dw.stats.FilteredByLines++
			continue
		}

		// 语言过滤，先检查黑名单，再匹配白名单
		lang := getLanguage(path, content)
		if len(dw.options.DisallowedCodeLanguages) > 0 {
			if _, ok := dw.options.DisallowedCodeLanguages[lang]; ok {
				continue
			}
		}
		if len(dw.options.AllowedCodeLanguages) > 0 {
			if _, ok := dw.options.AllowedCodeLanguages[lang]; !ok {
				dw.stats.FilteredByLanguage++
				continue
			}
		} else {
			if !dw.options.AllowUnknownLanguage && lang == "unknown" {
				dw.stats.FilteredByLanguage++
				continue
			}
		}
		if lang != "" {
			dw.statsLock.Lock()
			var flb types.FilesLinesBytes
			var ok bool
			if flb, ok = dw.stats.Languages[lang]; !ok {
				dw.stats.Languages[lang] = types.FilesLinesBytes{
					NumFiles: 1,
					NumLines: fileLines,
					NumBytes: len(content),
				}
			} else {
				flb.NumFiles = flb.NumFiles + 1
				flb.NumLines = flb.NumLines + fileLines
				flb.NumBytes = flb.NumBytes + len(content)
				dw.stats.Languages[lang] = flb
			}
			dw.statsLock.Unlock()
		}

		// 总行数
		dw.stats.TotalLinesOfCode += fileLines

		// 文件的 repo path
		var repoLocalPath, pathInRepo string
		if !info.isGitRemoteRepo {
			repoLocalPath, pathInRepo = dw.getFileRepoPath(path)
		} else {
			pathInRepo = path
		}

		// ctags
		var entries []*types.CTagsEntry
		if dw.ctagsParsers != nil {
			entries, _ = dw.ctagsParsers[shard].Parse(path, content)
		}

		// 最终通过，发送到下游
		dw.stats.IndexedFiles++

		dw.fileChan <- WalkerFileInfo{
			AbsPath:       path,
			Size:          int64(len(content)),
			Language:      lang,
			Content:       content,
			PathInRepo:    pathInRepo,
			RepoLocalPath: repoLocalPath,
			RepoRemoteURL: info.repoRemoteURL,
			CTagsEntries:  entries,
		}
	}
}

func (dw *IndexWalker) getFileRepoPath(filePath string) (string, string) {
	repoPath := ""
	pathInRepo := ""

	dw.reposLock.RLock()
	defer dw.reposLock.RUnlock()

	for repo := range dw.repos {
		if strings.HasPrefix(filePath, repo+"/") {
			repoPath = repo
			pathInRepo = strings.TrimPrefix(filePath, repo+"/")
			break
		}
	}

	return repoPath, pathInRepo
}

func checkIfBinaryFile(content []byte) bool {
	for _, c := range content {
		if c == 0 {
			return true
		}
	}

	return enry.IsBinary(content)
}

func getLanguage(path string, content []byte) string {
	lang := strings.ToLower(enry.GetLanguage(path, content))

	if lang == "" {
		return "unknown"
	}

	lang = strings.Join(strings.Split(lang, " "), "_")
	lang = strings.Join(strings.Split(lang, "."), "_")
	lang = strings.Join(strings.Split(lang, "-"), "_")
	lang = strings.ReplaceAll(lang, "++", "pp")
	return lang
}

func numLines(content []byte) int {
	lines := 0
	var c byte
	for _, c = range content {
		if c == '\n' {
			lines++
		}
	}
	if c != 0 && c != '\n' {
		lines++
	}
	return lines
}
