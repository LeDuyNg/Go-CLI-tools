package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader, countLines bool, countBytes bool) (int, error) {
	scanner := bufio.NewScanner(r)

	switch {
	case countLines:
		scanner.Split(bufio.ScanLines)
	case countBytes:
		scanner.Split(bufio.ScanBytes)
	default:
		scanner.Split(bufio.ScanWords)
	}

	wc := 0
	for scanner.Scan() {
		wc++
	}

	return wc, scanner.Err()

}

func main() {
	lines := flag.Bool("l", false, "Count lines")
	bytes := flag.Bool("b", false, "Count bytes")
	flag.Parse()
	n, err := count(os.Stdin, *lines, *bytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(n)

}
