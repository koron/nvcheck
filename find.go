package main

import (
	"fmt"

	"github.com/koron/nvcheck/internal/ahocorasick"
)

func find(m *ahocorasick.Matcher, path string) error {
	c := &ctx{m: m, fname: path}
	if err := c.load(); err != nil {
		return err
	}
	if err := c.find(); err != nil {
		return err
	}
	return c.forFounds(func(f *Found) error {
		lnum := c.lnum(f.Begin)
		fmt.Printf("%s:%d: %s >> %s\n", c.fname, lnum, f.Word.Text, *f.Word.Fix)
		return nil
	})
}
