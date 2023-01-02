package hyphenator

import (
	"bufio"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/dghubble/trie"
)

const (
	folder = "languages"
)

// Dictionary represents a dictionary for hyphenation rules: patterns and exceptions.
//
// Dictionary knows how to load a TeX hyphenation patterns file from supported languages.
// Pattern files may be downloaded from https://www.tug.org/tex-hyphen/ or from the CTAN TeX archive.
// (e.g. https://www.ctan.org/pkg/ushyph).
//
// This packages comes by default with support for:
// * American English: en-US patterns (2005)
// * British English: en-GB patterns (1996)
// * German: de patterns (1996)
// * French:  fr-FR patterns
// * Spanish: es patterns
type Dictionary struct {
	exceptions trie.Trier // e.g., "computer" => [3,5] = "com-pu-ter"
	patterns   trie.Trier // where we store patterns and positions
	Identifier string     // Identifies the dictionary
}

func isTeXComment(line string) bool {
	return len(line) == 0 ||
		strings.HasPrefix(line, "%") ||
		strings.HasPrefix(line, `\`) ||
		strings.HasPrefix(line, "}")
}

// SupportedPatterns lists all currently supported language files with hyphenation patterns.
func SupportedPatterns() ([]string, error) {
	entries, err := texFS.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		result = append(result, entry.Name())
	}

	sort.Strings(result)

	return result, nil
}

// LoadPatterns loads a pattern file as a Dictionary
// from the embedded file system.
//
// Patterns are enclosed like so:
//
//	\patterns{ % some comment
//	 ...
//	.wil5i
//	.ye4
//	4ab.
//	a5bal
//	a5ban
//	abe2
//	 ...
//	}
//
// Odd numbers stand for possible discretionaries, even numbers forbid
// hyphenation. Digits belong to the character immediately after them, i.e.,
//
// Example:
//
//	"a5ban" => (a)(5b)(a)(n) => positions["aban"] = [0,5,0,0].
func LoadPatterns(patternfile string) (*Dictionary, error) {
	const (
		messageSection    = `\message{`
		exceptionsSection = `\hyphenation{`
	)

	file, err := texFS.Open(path.Join(folder, patternfile)) // known pattern files are embedded
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	dict := &Dictionary{
		exceptions: trie.NewRuneTrie(),
		patterns:   trie.NewRuneTrie(),
		Identifier: fmt.Sprintf("patterns: %s", patternfile), // default identifier
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, messageSection):
			// extract the patterns identifier
			dict.Identifier = strings.TrimPrefix(line, messageSection)

		case strings.HasPrefix(line, exceptionsSection):
			// decode the exceptions section
			for scanner.Scan() {
				line = scanner.Text()
				dict.exceptions.Put(dict.readException(line))
			}

		case isTeXComment(line):
			// ignore comments, TeX commands, etc.
			continue

		default:
			// decode a patterns section: ".ab1a" "abe4l3in", ...
			patterns := strings.Fields(line)
			for _, pattern := range patterns {
				dict.patterns.Put(dict.readPattern(pattern))
			}
		}
	}

	return dict, nil
}

// readPattern reads a pattern in the patterns section.
func (dict *Dictionary) readPattern(line string) (string, []int) {
	var (
		pattern   strings.Builder // will become the pattern without positions
		positions []int           // we'll extract positions
		wasdigit  bool            // has the last char been a digit?
	)

	for _, char := range line { // iterate over the runes for this pattern
		if unicode.IsDigit(char) {
			d, _ := strconv.Atoi(string(char))
			positions = append(positions, d) // add to positions array
			wasdigit = true

			continue
		}

		// '.' or alphabetic rune
		pattern.WriteRune(char)
		if wasdigit {
			wasdigit = false
		} else {
			positions = append(positions, 0) // append a 0
		}
	}

	return pattern.String(), positions
}

// readExceptions processes a line from the exceptions section in a pattern file
// ("\hyphenation{").
//
// Exceptions are encoded as predefined hyphenation points for known words:
//
//	ex-cep-tion
//	ta-ble
func (dict *Dictionary) readException(line string) (word string, positions []int) {
	if isTeXComment(line) {
		return
	}

	var washyphen bool
	for _, char := range line {
		switch {
		case hyphens.Contains(char):
			positions = append(positions, 1) // possible break point
			washyphen = true
		case washyphen: // skip letter
			washyphen = false
		default: // a letter without a '-'
			positions = append(positions, 0)
		}
	}
	word = strings.ReplaceAll(line, "-", "")

	return
}

// String returns the identifier of the pattern file (by default, this is the file name).
func (dict *Dictionary) String() string {
	return dict.Identifier
}
