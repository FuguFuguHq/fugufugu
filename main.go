package main

import (
	"flag"
	"fmt"
	"fugufugu/fugu"
	"os"
)

var Products map[string]fugu.Product

func main() {
	Products = fugu.Products()

	var checkUrl string
	var verbose bool
	var checkForCookie bool
	var maxPages uint64
	flag.StringVar(&checkUrl, "url", "", "url to check")
	flag.Uint64Var(&maxPages, "max", 10000, "max pages to check")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.BoolVar(&checkForCookie, "cookie", false, "check for cookies")
	flag.Parse()

	if len(checkUrl) == 0 {
		fmt.Println("No url specified.")
		os.Exit(-1)
	}

	externals := make(map[string]fugu.Privacy)

	scanner := fugu.NewCollector(maxPages, checkForCookie, checkUrl, externals, verbose)
	scanner.Collector.Visit(checkUrl)

	privacies := fugu.FromExternals(Products, externals)

	r := fugu.ReportFromScanner(checkUrl, checkForCookie, privacies, scanner)
	fugu.PrintReport(r)
}
