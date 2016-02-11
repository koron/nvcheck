package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/koron/nvcheck/internal/ahocorasick"
	"gopkg.in/yaml.v2"
)

type Dict map[string][]string

type Word struct {
	Text string
	Fix  *string
}

func loadDict(name string) (Dict, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readDict(f)
}

func readDict(r io.Reader) (Dict, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var d = make(Dict)
	err = yaml.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d Dict) toM() (*ahocorasick.Matcher, error) {
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
