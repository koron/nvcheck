package linereader

import (
	"bufio"
	"bytes"
	"io"
)

type LineReader struct {
	r   *bufio.Reader
	num int
}

func New(r io.Reader) *LineReader {
	return &LineReader{
		r: bufio.NewReader(r),
	}
}

func (lr *LineReader) LineNum() int {
	return lr.num
}

func (lr *LineReader) ReadLine() (*string, error) {
	b, err := lr.readLine()
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}
	lr.num++
	s := b.String()
	return &s, nil
}

func (lr *LineReader) readLine() (*bytes.Buffer, error) {
	b, isPrefix, err := lr.r.ReadLine()
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	bb := bytes.NewBuffer(b)
	if isPrefix {
		for {
			b, cont, err := lr.r.ReadLine()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			if _, err := bb.Write(b); err != nil {
				return nil, err
			}
			if !cont {
				break
			}
		}
	}
	return bb, nil
}
