package types

type IndexRepoInfo struct {
	// 如果这个值不为 0 那么从外部传入 RepoID
	// 这允许使用 SearchRequest.ShouldRecallRepo 传入权限控制钩子函数
	RepoID uint64

	// 仓库在操作系统中的路径
	RepoLocalPath string

	// 仓库的远程路径
	RepoRemoteURL string
}
