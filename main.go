package main

import (
	"flag"
	"log"
	"os"
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

	var hasError bool
	for _, n := range flag.Args() {
		if !findFile(m, n) {
			hasError = true
		}
	}
	if hasError {
		os.Exit(1)
	}
}
