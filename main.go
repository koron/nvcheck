package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/koron/nvcheck/internal/ahocorasick"
)

var (
	dict = flag.String("d", "dict.yml", "variability dictionary")
	fpat = flag.String("m", `\.(txt)$`, "pattern of the file")
	fre  *regexp.Regexp
)

func findFile(s Summary, m *ahocorasick.Matcher, name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return filepath.Walk(name, func(path string, fi os.FileInfo, err error) error {
			if fi != nil && fi.Mode().IsRegular() {
				if !fre.MatchString(path) {
					return nil
				}
				checkFile(s, m, path)
			}
			return nil
		})
	}
	if !fre.MatchString(name) {
		return nil
	}
	return checkFile(s, m, name)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	re, err := regexp.Compile(*fpat)
	if err != nil {
		log.Fatal(err)
	}
	fre = re

	d, err := loadDict(*dict)
	if err != nil {
		log.Fatal(err)
	}
	m, err := toMatcher(d)
	if err != nil {
		log.Fatal(err)
	}

	summary := make(Summary)
	for _, n := range flag.Args() {
		if err := findFile(summary, m, n); err != nil {
			log.Printf("failed to operate file: %v", n)
		}
	}
	if len(summary) > 0 {
		fmt.Println(summary)
		os.Exit(1)
	}
}
