package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"unicode"
	"unicode/utf8"

	"github.com/koron/go-debug"
	"github.com/koron/nvcheck/internal/ahocorasick"
)

var (
	// ErrFound indicate "found variability" by forFounds.
	ErrFound = errors.New("found variability")

	errCont = errors.New("continue")
)

type ctx struct {
	fname string
	m     *ahocorasick.Matcher

	content string
	it      *ahocorasick.Iter
	loffs   []int
	lt      bool
	ln      int

	founds []*Found
}

func newCtx(m *ahocorasick.Matcher, path string) (*ctx, error) {
	c := &ctx{m: m, fname: path}
	if err := c.load(); err != nil {
		return nil, err
	}
	if err := c.find(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ctx) load() error {
	b, err := ioutil.ReadFile(c.fname)
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
	c.lt = true
	c.ln = 1
}

// parse parses a rune.
func (c *ctx) parse(i int, r rune) (rune, error) {
	if c.lt {
		if r == '\n' {
			c.ln++
			c.loffs = append(c.loffs, i+1)
			// through
		} else if unicode.IsSpace(r) {
			if !c.it.Has(' ') {
				return 0, errCont
			}
			r = ' '
		}
	} else {
		if r == '\n' {
			c.lt = true
			c.ln++
			c.loffs = append(c.loffs, i+1)
			if !c.it.Has(' ') {
				return 0, errCont
			}
			r = ' '
		}
	}
	c.lt = false
	return r, nil
}

func (c *ctx) find() error {
	for i, r := range c.content {
		r2, err := c.parse(i, r)
		if err == errCont {
			continue
		} else if err != nil {
			return err
		}
		ev := c.it.Put(r2)
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

func (c *ctx) forFounds(proc func(*Found) error) error {
	has := false
	for _, f := range c.founds {
		if f.OK() {
			continue
		}
		has = true
		err := proc(f)
		if err != nil {
			return err
		}
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
