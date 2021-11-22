package types

// 遍历器统计指标
type IndexWalkerStats struct {
	IndexedRepos int
	IndexedDirs  int
	IndexedFiles int

	FilteredDirs         int
	FilteredBySize       int
	FilteredByLines      int
	FilteredByExtension  int
	FilteredByLanguage   int
	FilteredByBinaryType int
	FilteredByVendor     int
	FilteredByGenerated  int
	FilteredByImage      int
	FilteredByDotPrefix  int
	FilteredByError      int

	Languages        map[string]FilesLinesBytes
	TotalLinesOfCode int

	Message      string
	CurrentFile  string
	CurrentError string

	GitDirError int
}

type FilesLinesBytes struct {
	NumFiles int
	NumLines int
	NumBytes int
}
