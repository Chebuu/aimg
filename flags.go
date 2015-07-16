package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/ts"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-w width] file [file...]\n", os.Args[0])
	flag.PrintDefaults()
}

func getColumns(widthstr string) int {
	numstr := ""
	format := ""
	if strings.HasSuffix(widthstr, "%") {
		if len(widthstr) < 2 {
			fmt.Fprintf(os.Stderr, "Invalid percentage.\n")
			os.Exit(1)
		}
		numstr = widthstr[:len(widthstr)-1]
		format = "percent"
	} else if strings.Contains(widthstr, ".") {
		numstr = widthstr
		format = "decimal"
	} else {
		numstr = widthstr
		format = "columns"
	}

	num := 0
	if format == "decimal" {
		f, err := strconv.ParseFloat(numstr, 64)
		handleErr(err)
		num = int(100.0 * f)
	} else {
		var err error
		num, err = strconv.Atoi(numstr)
		handleErr(err)
	}

	if format == "columns" && num > 0 {
		return num
	}

	size, err := ts.GetSize()
	if err != nil {
		fmt.Fprintln(os.Stderr, err, "\nYou may need to "+
			"set width manually using -w num")
		os.Exit(2)
	}

	if format == "columns" {
		return size.Col()
	}

	cols := float64(size.Col())
	return int(cols * (float64(num) / 100.0))
}
