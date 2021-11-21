package types

type IndexRepoInfo struct {
	// 仓库在操作系统中的路径
	RepoLocalPath string

	// 仓库的远程路径
	RepoRemoteURL string
}
