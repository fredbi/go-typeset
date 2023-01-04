package hyphenator

import (
	"bufio"
	"fmt"
	"path"
	"sort"
	"strings"
	"unicode"

	"github.com/fredbi/go-typeset/wordbreak/hyphenator/internal/trie"
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

func atoiRune(r rune) int {
	switch r {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	default:
		return 0
	}
}

// readPattern reads a pattern in the patterns section.
func (dict *Dictionary) readPattern(line string) ([]rune, []int) {
	wasdigit := false                     // has the last char been a digit?
	pattern := make([]rune, 0, len(line)) // will become the pattern without positions
	positions := make([]int, 0, 10)       // we'll extract positions

	for _, char := range line { // iterate over the runes for this pattern
		if unicode.IsDigit(char) {
			d := atoiRune(char)
			positions = append(positions, d) // add to positions array
			wasdigit = true

			continue
		}

		// '.' or alphabetic rune
		pattern = append(pattern, unicode.ToLower(char))
		if wasdigit {
			wasdigit = false
		} else {
			positions = append(positions, 0) // append a 0
		}
	}

	return pattern, positions
}

// readExceptions processes a line from the exceptions section in a pattern file
// ("\hyphenation{").
//
// Exceptions are encoded as predefined hyphenation points for known words:
//
//	ex-cep-tion
//	ta-ble
func (dict *Dictionary) readException(line string) ([]rune, []int) {
	if isTeXComment(line) {
		return nil, nil
	}

	positions := make([]int, 0, 5)

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
	word := strings.ToLower(strings.ReplaceAll(line, "-", ""))

	return []rune(word), positions
}

// String returns the identifier of the pattern file (by default, this is the file name).
func (dict *Dictionary) String() string {
	return dict.Identifier
}
