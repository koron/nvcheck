package main

import (
	"bufio"
	"os"

	"github.com/koron/nvcheck/internal/ahocorasick"
)

type foundIter func() *Found

func nextFoundIter(founds []*Found) foundIter {
	return func() *Found {
		if len(founds) == 0 {
			return nil
		}
		f := founds[0]
		founds = founds[1:]
		return f
	}
}

func rewrite(m *ahocorasick.Matcher, path string) error {
	c, err := newCtx(m, path)
	if err != nil {
		return err
	}
	// open file to rewrite.
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var (
		buf  = bufio.NewWriter(f)
		iter = nextFoundIter(c.founds)
		next *Found
		curr *Found
		ir, iw int
	)
	next = iter()
	for i, r := range c.content {
		if curr == nil && next.IsBegin(i) {
			curr, ir, iw = next, 0, 0
		}
		for next.In(i) {
			next = iter()
		}
		if curr != nil {
			// TODO: rewrite founds
			_, _ = ir, iw
		}
		_, err := buf.WriteRune(r)
		if err != nil {
			return err
		}
	}
	// flush and close.
	err = buf.Flush()
	if err != nil {
		return err
	}
	return nil
}
