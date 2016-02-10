package main

import (
	"bytes"
	"fmt"
)

type Fix struct {
	ln    int
	found *Found
}

type Result map[string][]Fix

type Summary map[string]Result

func (s Summary) Add(name string, ln int, m *Found) {
	_, ok := s[m.Fix]
	if !ok {
		s[m.Fix] = make(Result)
	}
	s[m.Fix][name] = append(s[m.Fix][name], Fix{ln, m})
}

func (s Summary) String() string {
	var buf bytes.Buffer
	for k, v := range s {
		fmt.Fprintln(&buf, k)
		for kk, vv := range v {
			fmt.Fprintln(&buf, "  "+kk)
			for _, f := range vv {
				fmt.Fprintf(&buf, "    %s at %d\n", f.found.Text, f.ln)
			}
		}
	}
	return buf.String()
}
