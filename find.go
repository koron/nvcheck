package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/koron/go-debug"
	"github.com/koron/nvcheck/internal/ahocorasick"
)

var (
	ErrFound = errors.New("found variability")
)

type Found struct {
	Begin int
	End   int
	Text  string
	Fix   string
}

type ctx struct {
	fname string
	m     *ahocorasick.Matcher

	content string
	it      *ahocorasick.Iter
	loffs   []int

	has    bool
	founds []*Found
}

func (c *ctx) load() error {
	f, err := os.Open(c.fname)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	c.content = string(b)
	c.it = c.m.Iter()
	// it assumes that a line has 50 bytes in average.
	c.loffs = append(make([]int, 0, len(c.content)/50+1), 0)
	return nil
}

func (c *ctx) find() error {
	if err := c.load(); err != nil {
		return err
	}
	var (
		lineTop = true
		lnum = 1
	)
	for i, r := range c.content {
		if lineTop {
			if r == '\n' {
				lnum++
				c.loffs = append(c.loffs, i+1)
				// through
			} else if unicode.IsSpace(r) {
				continue
			}
		} else {
			if r == '\n' {
				lineTop = true
				lnum++
				c.loffs = append(c.loffs, i+1)
				continue
			}
		}
		lineTop = false
		ev := c.it.Put(r)
		if ev == nil {
			c.flush(i)
			continue
		}
		for d := ev.Next(); d != nil; d = ev.Next() {
			w, _ := d.Value.(*Word)
			_, n := utf8.DecodeRuneInString(c.content[i:])
			top := c.top(i+n, w.Text)
			if top < 0 {
				return fmt.Errorf("match failure for %q in file %s at offset %d", w.Text, c.fname, i+n)
			}
			if w.Fix != nil {
				c.flush(i)
				c.founds = append(c.founds, &Found{
					Begin: top,
					End:   i + n,
					Text:  w.Text,
					Fix:   *w.Fix,
				})
				continue
			}
			c.flush(top)
		}
	}
	c.flush(len(c.content) - 1)
	if c.has {
		return ErrFound
	}
	return nil
}

func (c *ctx) flush(top int) {
	if len(c.founds) <= 0 {
		return
	}
	debug.Printf("flush: %d", top)
	for _, f := range c.founds {
		if top <= f.Begin {
			debug.Printf("  IGN: %#v", f)
			continue
		}
		debug.Printf("  HIT: %#v", f)
		lnum := c.lnum(f.Begin)
		fmt.Printf("%s:%d: %s >> %s\n", c.fname, lnum, f.Text, f.Fix)
	}
	c.has = true
	c.founds = c.founds[:0]
}

func (c *ctx) lnum(off int) int {
	return c.searchLoffs(off, 0, len(c.loffs)) + 1
}

func (c *ctx) searchLoffs(off, start, end int) int {
	if start+1 >= end {
		return start
	}
	mid := (start + end) / 2
	pivot := c.loffs[mid]
	if off < pivot {
		return c.searchLoffs(off, start, mid)
	}
	return c.searchLoffs(off, mid, end)
}

// top returns offset to start of an match.
func (c *ctx) top(tail int, w string) int {
	for len(w) > 0 {
		if tail <= 0 {
			return -1
		}
		r1, n1 := utf8.DecodeLastRuneInString(c.content[:tail])
		tail -= n1
		if unicode.IsSpace(r1) {
			continue
		}
		r2, n2 := utf8.DecodeLastRuneInString(w)
		w = w[:len(w)-n2]
		if r1 != r2 {
			return -1
		}
	}
	return tail
}

func find(m *ahocorasick.Matcher, path string) error {
	c := &ctx{m: m, fname: path}
	return c.find()
}
