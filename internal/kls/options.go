package kls

type KLSOptions struct {
	indexLocations      []string
	fileExtensionFilter string
	ignoreDirs          string
	numContextLines     int
	ctagBinaryPath      string
}

func NewKLSOptions() *KLSOptions {
	return &KLSOptions{}
}

func (options *KLSOptions) SetIndexDirs(dirs []string) *KLSOptions {
	options.indexLocations = dirs
	return options
}

func (options *KLSOptions) SetFileExtensionFilter(filter string) *KLSOptions {
	options.fileExtensionFilter = filter
	return options
}

func (options *KLSOptions) SetIgnoreDirs(dirs string) *KLSOptions {
	options.ignoreDirs = dirs
	return options
}

func (options *KLSOptions) SetCtagBinaryPath(path string) *KLSOptions {
	options.ctagBinaryPath = path
	return options
}
