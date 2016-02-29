package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"

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

func replaceToStdout(m *ahocorasick.Matcher, path string) error {
	c, err := newCtx(m, path)
	if err != nil {
		return err
	}
	return replaceFounds(c, os.Stdout)
}

func replaceInPlace(m *ahocorasick.Matcher, path string) error {
	c, err := newCtx(m, path)
	if err != nil {
		return err
	}
	// open file to replace in place.
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return replaceFounds(c, f)
}

type unexpectedRune struct {
	ex, ac rune
}

func (e *unexpectedRune) Error() string {
	return fmt.Sprintf("unexpected rune: expected=%c actually=%c", e.ex, e.ac)
}

func replaceFounds(c *ctx, w io.Writer) error {
	var (
		buf     = bufio.NewWriter(w)
		iter    = nextFoundIter(c.founds)
		next    *Found
		ng, fix []rune
	)
	next = iter()
	for i, r := range c.content {
		if len(ng) == 0 && next.IsBeginAndFix(i) {
			ng = []rune(next.Word.Text)
			fix = []rune(*next.Word.Fix)
		}
		for next.In(i) {
			next = iter()
		}
		if len(ng) > 0 {
			if !unicode.IsSpace(r) {
				if r != ng[0] {
					return &unexpectedRune{ex: ng[0], ac: r}
				}
				ng = ng[1:]
				if len(ng) == 0 {
					if len(fix) > 0 {
						_, err := buf.WriteString(string(fix))
						if err != nil {
							return err
						}
					}
					continue
				}
				r, fix = fix[0], fix[1:]
			} else {
				if unicode.IsSpace(ng[0]) {
					// TODO: consume a rune from ng
					continue
				}
			}
		}
		_, err := buf.WriteRune(r)
		if err != nil {
			return err
		}
	}
	return buf.Flush()
}
