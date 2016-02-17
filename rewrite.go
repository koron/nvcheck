package main

import (
	"bufio"
	"os"

	"github.com/koron/nvcheck/internal/ahocorasick"
)

func rewrite(m *ahocorasick.Matcher, path string) error {
	c := &ctx{m: m, fname: path}
	if err := c.load(); err != nil {
		return err
	}
	if err := c.find(); err != nil {
		return err
	}
	// open file to rewrite.
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var (
		buf = bufio.NewWriter(f)
		founds = c.founds
		next *Found
		curr *Found
	)
	nextFound := func() *Found {
		if len(founds) == 0 {
			return nil
		}
		f := founds[0]
		founds = founds[1:]
		return f
	}
	nextFound()
	for i, r := range c.content {
		// TODO: rewrite founds
		if curr != nil {
			// TODO:
		}
		if curr == nil {
			for next != nil {
				if  i== next.Begin {
					// TODO:
				} else if i >= next.End {
					next = nextFound()
				}
			}
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
