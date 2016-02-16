package main

import (
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
	// TODO: rewrite founds
	return nil
}
