//go:build ignore

package main

// Generate runes_tables.go from unicode property tables published at https://unicode.org.

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type (
	// codePointValue represents the East Asian Width unicode property,
	// as defined by https://www.unicode.org/reports/tlow1.
	//
	// The default is "N".
	codePointValue string

	runesRange struct {
		lo rune
		hi rune
	}
)

const (
	packageName   = "runes"
	generatedFile = "runes_tables.go"
	unicodeSource = "https://unicode.org/Public/15.0.0/ucd"

	codePointValueA  codePointValue = "A"  // ambiguous
	codePointValueF  codePointValue = "F"  // full-width
	codePointValueH  codePointValue = "H"  // half-width
	codePointValueN  codePointValue = "N"  // neutral (not East-Asian characters)
	codePointValueNa codePointValue = "Na" // narrow
	codePointValueW  codePointValue = "W"  // wide
)

func main() {
	f := new(bytes.Buffer)
	fmt.Fprint(f, "// Code generated by build_tables.go. DO NOT EDIT.\n\n")
	fmt.Fprintf(f, "package %s\n\n", packageName)

	targetWidthProperties := fmt.Sprintf("%s/%s", unicodeSource, "EastAsianWidth.txt")
	fmt.Fprintf(f, "\n// unicode EastAsianWidth properties retrieved from %s\n", targetWidthProperties)

	targetEmojiProperties := fmt.Sprintf("%s/%s", unicodeSource, "emoji/emoji-data.txt")
	fmt.Fprintf(f, "// unicode emoji properties retrieved from %s\n\n", targetEmojiProperties)

	// retrieve unicode property EastAsianWidth
	resp, err := http.Get(targetWidthProperties)
	if err != nil {
		log.Fatal("retrieving %s: %v", targetWidthProperties, err)
	}
	defer resp.Body.Close()

	_ = eastasian(f, resp.Body)

	// retrieve unicode emojis with Extended_Pictographics
	resp, err = http.Get(targetEmojiProperties)
	if err != nil {
		log.Fatal("retrieving %s: %v", targetEmojiProperties, err)
	}
	defer resp.Body.Close()

	_ = emoji(f, resp.Body)

	out, err := format.Source(f.Bytes())
	if err != nil {
		log.Fatal("formatting generated source: %v.\n%s", err, f.String())
	}
	err = ioutil.WriteFile(generatedFile, out, 0666)
	if err != nil {
		log.Fatal("writing generated file: %v", err)
	}
}

func generate(out io.Writer, v string, arr []runesRange) {
	fmt.Fprintf(out, "var %s = table{\n\t", v)

	for i := 0; i < len(arr); i++ {
		fmt.Fprintf(out, "{0x%04X, 0x%04X},", arr[i].lo, arr[i].hi)
		if i < len(arr)-1 {
			if i%3 == 2 {
				fmt.Fprint(out, "\n\t")
			} else {
				fmt.Fprint(out, " ")
			}
		}
	}

	fmt.Fprintln(out, "\n}")
}

func shapeup(p []runesRange) []runesRange {
	arr := p

	for i := 0; i < len(arr)-1; i++ {
		if arr[i].hi+1 == arr[i+1].lo {
			lo := arr[i].lo
			arr = append(arr[:i], arr[i+1:]...)
			arr[i].lo = lo
			i--
		}
	}

	return arr
}

func eastasian(out io.Writer, in io.Reader) error {
	scanner := bufio.NewScanner(in)

	dbl := []runesRange{}
	amb := []runesRange{}
	cmb := []runesRange{}
	na := []runesRange{}
	nu := []runesRange{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		var (
			low, high rune
			class     string
		)
		n, err := fmt.Sscanf(line, "%x..%x;%s ", &low, &high, &class)
		if err != nil || n == 2 {
			n, err = fmt.Sscanf(line, "%x;%s ", &low, &class)
			if err != nil || n != 2 {
				continue
			}
			high = low
		}

		if strings.Index(line, "COMBINING") != -1 {
			cmb = append(cmb, runesRange{
				lo: low,
				hi: high,
			})
		}

		switch codePointValue(class) {
		case codePointValueW, codePointValueF:
			dbl = append(dbl, runesRange{
				lo: low,
				hi: high,
			})
		case codePointValueA:
			amb = append(amb, runesRange{
				lo: low,
				hi: high,
			})
		case codePointValueNa:
			na = append(na, runesRange{
				lo: low,
				hi: high,
			})
		case codePointValueN:
			nu = append(nu, runesRange{
				lo: low,
				hi: high,
			})
		}
	}

	cmb = shapeup(cmb)
	generate(out, "combining", cmb)
	fmt.Fprintln(out)

	dbl = shapeup(dbl)
	generate(out, "doublewidth", dbl)
	fmt.Fprintln(out)

	amb = shapeup(amb)
	generate(out, "ambiguous", amb)
	fmt.Fprint(out)

	na = shapeup(na)
	generate(out, "narrow", na)
	fmt.Fprintln(out)

	nu = shapeup(nu)
	generate(out, "neutral", nu)
	fmt.Fprintln(out)

	return nil
}

func emoji(out io.Writer, in io.Reader) error {
	scanner := bufio.NewScanner(in)
	arr := []runesRange{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.Contains(line, "Extended_Pictographic") {
			continue
		}

		var low, high rune
		n, err := fmt.Sscanf(line, "%x..%x ", &low, &high)
		if err != nil || n == 1 {
			n, err = fmt.Sscanf(line, "%x ", &low)
			if err != nil || n != 1 {
				continue
			}
			high = low
		}
		if high < 0xFF {
			continue
		}

		arr = append(arr, runesRange{
			lo: low,
			hi: high,
		})
	}

	arr = shapeup(arr)
	generate(out, "emoji", arr)

	return nil
}
