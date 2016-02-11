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
	Word  *Word
}

func (f *Found) String() string {
	if f.Word.Fix != nil {
		return fmt.Sprintf("Found{Begin:%d, End:%d, Text:%q, Fix:%q}",
		f.Begin, f.End, f.Word.Text, *f.Word.Fix)
	}
	return fmt.Sprintf("Found{Begin:%d, End:%d, Text:%q}",
		f.Begin, f.End, f.Word.Text)
}

func (f *Found) OK() bool {
	return f.Word.Fix == nil
}

type ctx struct {
	fname string
	m     *ahocorasick.Matcher

	content string
	it      *ahocorasick.Iter
	loffs   []int

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
	c.setup(string(b))
	return nil
}

func (c *ctx) setup(s string) {
	c.content = s
	c.it = c.m.Iter()
	// it assumes that a line has 50 bytes in average.
	c.loffs = append(make([]int, 0, len(c.content)/50+1), 0)
}

func (c *ctx) find() error {
	var (
		lineTop = true
		lnum    = 1
	)
	for i, r := range c.content {
		if lineTop {
			if r == '\n' {
				lnum++
				c.loffs = append(c.loffs, i+1)
				// through
			} else if unicode.IsSpace(r) {
				if !c.it.Has(' ') {
					continue
				}
				r = ' '
			}
		} else {
			if r == '\n' {
				lineTop = true
				lnum++
				c.loffs = append(c.loffs, i+1)
				if !c.it.Has(' ') {
					continue
				}
				r = ' '
			}
		}
		lineTop = false
		ev := c.it.Put(r)
		if ev == nil {
			continue
		}
		for d := ev.Next(); d != nil; d = ev.Next() {
			w, _ := d.Value.(*Word)
			_, n := utf8.DecodeRuneInString(c.content[i:])
			top := c.top(i+n, w.Text)
			if top < 0 {
				return fmt.Errorf(
					"match failure for %q in file %s at offset %d",
					w.Text, c.fname, i+n)
			}
			err := c.push(&Found{
				Begin: top,
				End:   i + n,
				Word:  w,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *ctx) printFounds() error {
	has := false
	for _, f := range c.founds {
		if f.OK() {
			continue
		}
		has = true
		c.put(f)
	}
	if has {
		return ErrFound
	}
	return nil
}

func (c *ctx) push(f *Found) error {
	debug.Printf("push: %s", f)
	for {
		if len(c.founds) == 0 {
			// case 1 in doc/optmize-found-words.pdf
			debug.Printf("  case 1")
			c.founds = append(c.founds, f)
			break
		}
		last := c.founds[len(c.founds)-1]
		if f.End < last.End {
			return fmt.Errorf(
				"word %q ended at %d is before end of last word %q at %d",
				f.Word.Text, f.End, last.Word.Text, last.End)
		} else if f.End == last.End {
			if f.Begin > last.Begin {
				// case 4 in doc/optmize-found-words.pdf
				debug.Printf("  case 4: %s", last)
				break
			} else if f.Begin == last.Begin {
				// case 3 in doc/optmize-found-words.pdf with special.
				debug.Printf("  case 3: %s", last)
				if last.OK() != f.OK() {
					return fmt.Errorf(
						"word %q is registered as both good and bad word",
						f.Word.Text)
				}
				break
			}
			// case 2 in doc/optmize-found-words.pdf
			debug.Printf("  case 2: %s", last)
			c.founds = c.founds[:len(c.founds)-1]
		} else {
			if f.Begin > last.Begin {
				// case 6 in doc/optmize-found-words.pdf
				debug.Printf("  case 6: %s", last)
				c.founds = append(c.founds, f)
				break
			}
			// case 5 in doc/optmize-found-words.pdf
			debug.Printf("  case 5: %s", last)
			c.founds = c.founds[:len(c.founds)-1]
		}
	}
	return nil
}

func (c *ctx) put(f *Found) {
	lnum := c.lnum(f.Begin)
	fmt.Printf("%s:%d: %s >> %s\n", c.fname, lnum, f.Word.Text, *f.Word.Fix)
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
			debug.Printf("over backtrack: w=%q", w)
			return -1
		}
		wr, wn := utf8.DecodeLastRuneInString(w)
		cr, cn := utf8.DecodeLastRuneInString(c.content[:tail])
		tail -= cn
		if unicode.IsSpace(wr) {
			if !unicode.IsSpace(cr) {
				// no spaces which required.
				debug.Printf("not space: tail=%d w=%q cr=%q", tail, w, cr)
				return -1
			}
			w = w[:len(w)-wn]
			continue
		}
		if unicode.IsSpace(cr) {
			continue
		}
		w = w[:len(w)-wn]
		if cr != wr {
			// didn't match runes.
			debug.Printf("not match: tail=%d w=%q cr=%q wr=%q",
				tail, w, cr, wr)
			return -1
		}
	}
	return tail
}

func find(m *ahocorasick.Matcher, path string) error {
	c := &ctx{m: m, fname: path}
	if err := c.load(); err != nil {
		return err
	}
	if err := c.find(); err != nil {
		return err
	}
	return c.printFounds()
}
