package main

import (
	"flag"
	"log"
	"os"
)

var (
	dict    = flag.String("d", "dict.yml", "variability dictionary")
	replace = flag.Bool("r", false, "replace words to stdout")
	inplace = flag.Bool("i", false, "replace words in place")
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
	proc := find
	if *inplace {
		proc = replaceInPlace
	} else if *replace {
		proc = replaceToStdout
	}
	for _, n := range flag.Args() {
		err := proc(m, n)
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
