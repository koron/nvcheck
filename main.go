package main

import (
	"flag"
	"log"
)

var (
	dict = flag.String("d", "dict.yml", "variability dictionary")
)

func main() {
	flag.Parse()
	d, err := loadDict(*dict)
	if err != nil {
		log.Fatal(err)
	}
	m, err := toMatcher(d)
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range flag.Args() {
		findFile(m, n)
	}
}
