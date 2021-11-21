package walker

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

// 遍历  dirPath 目录下的所有文件，并将文件信息通过 fileChan 通道回传给上游
func (dw *IndexWalker) WalkGitRemote(gitURL string) {
	if dw.pullMode {
		// dryrun 模式下不爬取远程仓库
		return
	}

	dw.stats.Message = "正在爬取远程代码仓库：" + gitURL

	defer func() {
		dw.fileChan <- WalkerFileInfo{
			WalkingDone: true,
		}
	}()

	var r *git.Repository
	var err error
	pubKeys := dw.getPubKeys()
	if pubKeys != nil {
		r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:             gitURL,
			Depth:           1,
			Auth:            pubKeys,
			InsecureSkipTLS: true,
		})
	} else {
		r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:             gitURL,
			Depth:           1,
			InsecureSkipTLS: true,
		})
	}

	if err != nil {
		dw.stats.CurrentError = "远程代码仓库抓取失败 - " + err.Error()
		dw.fileChan <- WalkerFileInfo{
			Error: err,
		}
		return
	}

	ref, err := r.Head()
	if err != nil {
		dw.stats.CurrentError = "远程代码仓库抓取失败 - " + err.Error()
		dw.fileChan <- WalkerFileInfo{
			Error: err,
		}
		return
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		dw.stats.CurrentError = "远程代码仓库抓取失败 - " + err.Error()
		dw.fileChan <- WalkerFileInfo{
			Error: err,
		}
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		dw.stats.CurrentError = "远程代码仓库抓取失败 - " + err.Error()
		dw.fileChan <- WalkerFileInfo{
			Error: err,
		}
		return
	}

	dw.fileChan <- WalkerFileInfo{
		IsRepo:        true,
		RepoRemoteURL: gitURL,
	}
	dw.stats.IndexedRepos++

	err = tree.Files().ForEach(func(f *object.File) error {
		content, err := f.Contents()
		if err != nil {
			return err
		}

		dw.walkerInfoChan <- walkerInfo{
			path:            f.Name,
			isGitRemoteRepo: true,
			content:         []byte(content),
			repoRemoteURL:   gitURL,
		}

		return nil
	})
}
