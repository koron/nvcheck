package main

import (
	"os"
	"unicode"

	"github.com/koron/nvcheck/internal/ahocorasick"
	"github.com/koron/nvcheck/internal/linereader"
)

type Word struct {
	Text string
	Fix  *string
}

type Found struct {
	Off  int
	Text string
	Fix  string
}

func toMatcher(d Dict) (*ahocorasick.Matcher, error) {
	m := ahocorasick.New()
	for k, v := range d {
		m.Add(k, &Word{Text: k})
		k2 := k
		for _, w := range v {
			m.Add(w, &Word{Text: w, Fix: &k2})
		}
	}
	if err := m.Compile(); err != nil {
		return nil, err
	}
	return m, nil
}

func checkFile(s Summary, m *ahocorasick.Matcher, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return nil
	}
	defer f.Close()
	it := m.Iter()
	lr := linereader.New(f)
	loff := 0
	var last *Found
	for {
		l, err := lr.ReadLine()
		if err != nil {
			return nil
		}
		if l == nil {
			break
		}
		for i, r := range *l {
			if unicode.IsSpace(r) {
				continue
			}
			ev := it.Put(r)
			if ev == nil {
				if last != nil {
					s.Add(name, lr.LineNum(), last)
					last = nil
				}
				continue
			}
			for {
				d := ev.Next()
				if d == nil {
					break
				}
				off := loff + i - d.Offset
				w, _ := d.Value.(*Word)
				if w.Fix != nil {
					if last != nil {
						s.Add(name, lr.LineNum(), last)
					}
					last = &Found{
						Off:  off,
						Text: w.Text,
						Fix:  *w.Fix,
					}
					continue
				}
				if last == nil {
					continue
				}
				if off <= last.Off {
					last = nil
					continue
				}
				s.Add(name, lr.LineNum(), last)
				last = nil
			}
		}
		if last != nil {
			s.Add(name, lr.LineNum(), last)
			last = nil
		}
		loff += len(*l)
	}
	return nil
}
