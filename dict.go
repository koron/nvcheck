package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Dict map[string][]string

func loadDict(name string) (Dict, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
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
