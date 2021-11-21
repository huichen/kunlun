package ctags

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/huichen/kunlun/pkg/log"
	"github.com/huichen/kunlun/pkg/types"
)

var (
	logger = log.GetLogger()
)

type CTagsParser struct {
	options *types.CTagsParserOptions
	tempDir string

	cmd     *exec.Cmd
	in      io.WriteCloser
	out     *scanner
	outPipe io.ReadCloser
}

func NewCTagsParser(options *types.CTagsParserOptions) (*CTagsParser, error) {
	if options == nil {
		return nil, nil
	}

	// 检查二进制文件是否存在
	if _, err := os.Stat(options.BinaryPath); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("ctags 二进制文件 %s 不存在", options.BinaryPath)
	}

	// 创建临时文件夹
	tempDir, err := ioutil.TempDir("", "kunlun-ctags")
	if err != nil {
		return nil, err
	}

	// 对 Linux 系统，使用沙箱模式，这要求 Universal-Ctags 支持 seccomp 模式
	// 如果不支持，不会提取符号信息
	opt := "default"
	if runtime.GOOS == "linux" {
		opt = "sandbox"
	}

	// 设置输入输出等并启动命令
	cmd := exec.Command(
		options.BinaryPath,
		"--_interactive="+opt, // 交互模型，要求 ctags 编译时开启 libjansson 支持
		"--fields=*")
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		in.Close()
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	parser := CTagsParser{
		options: options,
		tempDir: tempDir,
		cmd:     cmd,
		in:      in,
		out:     &scanner{r: bufio.NewReaderSize(out, 4096)},
		outPipe: out,
	}
	var init reply
	if err := parser.read(&init); err != nil {
		return nil, err
	}

	return &parser, nil
}

func (cp *CTagsParser) Close() {
	cp.cmd.Process.Kill()
	cp.outPipe.Close()
	cp.in.Close()
}
