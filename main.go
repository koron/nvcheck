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
	m, err := d.toM()
	if err != nil {
		log.Fatal(err)
	}

	var found bool
	for _, n := range flag.Args() {
		err := find(m, n)
		if err != nil {
			if err == ErrFound {
				found = true
				continue
			}
			log.Fatal(err)
		}
	}
	if found {
		os.Exit(1)
	}
}
