package ctags

import (
	"bufio"
	"bytes"
	"io"
)

type scanner struct {
	r    *bufio.Reader
	line []byte
	err  error
}

func (s *scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	var (
		err  error
		line []byte
	)

	for err == nil && len(line) == 0 {
		line, err = s.r.ReadSlice('\n')
		for err == bufio.ErrBufferFull {
			line = nil
			_, err = s.r.ReadSlice('\n')
		}
		line = bytes.TrimSuffix(line, []byte{'\n'})
		line = bytes.TrimSuffix(line, []byte{'\r'})
	}

	s.line, s.err = line, err
	return len(line) > 0
}

func (s *scanner) Bytes() []byte {
	return s.line
}

func (s *scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
