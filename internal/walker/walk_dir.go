package walker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/karrick/godirwalk"
)

// 遍历  dirPath 目录下的所有文件，并将文件信息通过 fileChan 通道回传给上游
func (dw *IndexWalker) WalkDir(dirPath string) {
	dw.stats.Message = "索引目录：" + dirPath

	defer func() {
		dw.fileChan <- WalkerFileInfo{
			WalkingDone: true,
		}
	}()

	absDirPath, err := filepath.Abs(filepath.Clean(dirPath))
	if err != nil {
		dw.fileChan <- WalkerFileInfo{
			Error: err,
		}
		return
	}

	opts := godirwalk.Options{
		Callback:            dw.returnToChan,
		FollowSymbolicLinks: true,
		AllowNonDirectory:   true,
		ErrorCallback:       dw.errorCallback,
		Unsorted:            false,
	}

	if err := godirwalk.Walk(absDirPath, &opts); err != nil {
		dw.fileChan <- WalkerFileInfo{
			Error: err,
		}
		return
	}
}

func (dw *IndexWalker) errorCallback(path string, err error) godirwalk.ErrorAction {
	dw.stats.FilteredByError++
	dw.fileChan <- WalkerFileInfo{
		Error: err,
	}
	return godirwalk.SkipNode
}

func (dw *IndexWalker) returnToChan(path string, de *godirwalk.Dirent) error {
	dw.stats.CurrentFile = path
	dw.stats.Message = fmt.Sprintf("索引了%d个文件：%s", dw.stats.IndexedFiles, path)

	// 检查是否是文件夹
	isDir, err := de.IsDirOrSymlinkToDir()
	if err != nil {
		dw.stats.FilteredByError++
		dw.fileChan <- WalkerFileInfo{
			Error: err,
		}
		return nil
	}

	// 对文件夹做处理
	if isDir {
		// 检查是否需要跳过
		base := de.Name()
		if _, ok := dw.options.IgnoreDirs[base]; ok {
			dw.stats.FilteredDirs++
			return godirwalk.SkipThis
		}

		// 点开头的目录
		if dw.options.FilterDotPrefix && strings.HasPrefix(base, ".") {
			dw.stats.FilteredByDotPrefix++
			return godirwalk.SkipThis
		}

		// 存在 .git 目录的情况，读取仓库信息
		if _, err := os.Stat(path + "/.git"); !os.IsNotExist(err) {
			err := dw.processGitRepo(path)
			if err == godirwalk.SkipThis {
				return godirwalk.SkipThis
			}
			if err != nil {
				logger.Error(err)
				dw.stats.GitDirError++
				return godirwalk.SkipThis
			}
		}

		dw.stats.IndexedDirs++
		return nil
	}

	// 文件格式，发送到处理器
	if !dw.pullMode {
		dw.walkerInfoChan <- walkerInfo{path: path}
	}
	return nil
}

func (dw *IndexWalker) processGitRepo(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil
	}
	config, err := repo.Config()
	if err != nil {
		return err
	}

	var remoteURL string
	if config.Remotes != nil {
		if origin, ok := config.Remotes["origin"]; ok {
			// 只处理有一个 origin 的仓库
			if len(origin.URLs) == 1 {
				remoteURL = origin.URLs[0]

				// 白名单存在，且没有命中名单，则直接跳过
				if dw.options.AllowedRepoRemoteURLs != nil {
					if _, ok := dw.options.AllowedRepoRemoteURLs[remoteURL]; !ok {
						return godirwalk.SkipThis
					}
				}

				if dw.pullMode {
					tree, err := repo.Worktree()
					if err != nil {
						return err
					}
					logger.Infof("pulling 远程仓库 %s", remoteURL)
					err = tree.Pull(&git.PullOptions{
						RemoteName: "origin",
						Auth:       dw.getPubKeys(),
					})
					if err != nil && err != git.NoErrAlreadyUpToDate {
						return err
					}
				}
			} else {
				logger.Error("仓库 %s 的 origin 不止一个 URL，无法识别 remote URL", path)
			}
		}
	}

	repoLocalPath := strings.TrimSuffix(path, "/")
	if !dw.pullMode {
		dw.fileChan <- WalkerFileInfo{
			IsRepo:        true,
			RepoLocalPath: repoLocalPath,
			RepoRemoteURL: remoteURL,
		}
	}
	dw.stats.IndexedRepos++
	dw.reposLock.Lock()
	dw.repos[repoLocalPath] = true
	dw.reposLock.Unlock()

	return nil
}
