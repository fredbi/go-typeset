//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var texComment = regexp.MustCompile(`^(.+?)\s*%.*$`)

func main() {
	if err := filepath.WalkDir("tex", func(pth string, d fs.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".tex") {
			return nil
		}

		return strip(pth)
	}); err != nil {
		log.Fatal(err)
	}
}

func strip(pth string) error {
	file, err := os.Open(pth)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	destination := filepath.Base(pth)
	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func() {
		_ = output.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case isTeXComment(line):
			continue

		default:
			stripped := texComment.ReplaceAllString(line, "$1")
			fmt.Fprintln(output, strings.TrimSpace(stripped))
		}
	}

	return nil
}

func isTeXComment(line string) bool {
	return len(line) == 0 ||
		strings.HasPrefix(line, "%")
}
