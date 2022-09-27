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

	var euCheck bool
	var checkUrl string
	var verbose bool
	var checkForCookie bool
	var maxPages uint64
	var help bool
	flag.StringVar(&checkUrl, "url", "", "url to check")
	flag.Uint64Var(&maxPages, "max", 10000, "max pages to check")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.BoolVar(&checkForCookie, "cookie", false, "check for cookies")
	flag.BoolVar(&euCheck, "eu", false, "warn for using problematic external resources in EU")
	flag.BoolVar(&help, "help", false, "help information")
	flag.Parse()

	if help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	if len(checkUrl) == 0 {
		fmt.Println("No url specified.")
		os.Exit(0)
	}

	externals := make(map[string]fugu.Privacy)
	scanner := fugu.NewCollector(maxPages, checkForCookie, checkUrl, externals, verbose)
	scanner.Collector.Visit(checkUrl)

	privacies := fugu.FromExternals(Products, externals)

	r := fugu.ReportFromScanner(checkUrl, checkForCookie, euCheck, privacies, scanner)

	fugu.PrintReport(r)

}
