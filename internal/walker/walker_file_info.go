package walker

import "github.com/huichen/kunlun/internal/common_types"

type WalkerFileInfo struct {
	// 返回是否是代码仓库，如果为否则为文件
	IsRepo bool

	// 文件的绝对路径地址
	AbsPath string

	// 文件大小
	Size int64

	// 检测到的语言
	Language string

	// 文件内容
	Content []byte

	// 文件的仓库和仓库内路径
	PathInRepo string

	// 下面两个字段：
	// 1、当检测到 git 仓库（IsRepo == true）的时候返回仓库的本地地址和远程地址
	// 2、当检测到文件时（IsRepo == false），返回所属仓库的本地地址和远程地址（如果有的话）
	RepoLocalPath string
	RepoRemoteURL string

	Error error

	// 一次 IndexDir 请求完成
	WalkingDone bool

	// 返回解析的 CTags 标签
	CTagsEntries []*common_types.CTagsEntry
}
