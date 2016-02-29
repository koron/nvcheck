package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/koron/nvcheck/internal/ahocorasick"
)

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

func replaceCheck(r, r0 rune) (consume, fix bool, err error) {
	if !unicode.IsSpace(r) {
		if r != r0 {
			return false, false, &unexpectedRune{ex: r0, ac: r}
		}
		return true, true, nil
	} else if unicode.IsSpace(r0) {
		return true, false, nil
	}
	return false, false, nil
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
			c, f, err := replaceCheck(r, ng[0])
			if err != nil {
				return err
			}
			if c {
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
			}
			if f {
				r, fix = fix[0], fix[1:]
			}
		}
		_, err := buf.WriteRune(r)
		if err != nil {
			return err
		}
	}
	return buf.Flush()
}
