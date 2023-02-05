# hyphenator

A package to find legit hyphenation points for words.

It provides a golang implementation of the word hyphenation algorithm
proposed by Frank M. Liang, which is distributed with Tex.

Language-specific hyphenation rules are provided by pattern files that can be downloaded from https://www.tug.org/tex-hyphen/ 
or from the CTAN TeX archive (e.g. https://www.ctan.org/pkg/ushyph).

Currently, this package is distributed with support for the following languages:
* American English: en-US patterns (2005)
* British English: en-GB patterns (1996)
* German: de patterns (1996)
* French:  fr-FR patterns
* Spanish: es patterns

## Maintainance

To update pattern files or support new languages, download and add the desired files into the folder "languages/tex", with the ".tex" extension,
then run:
```
go generate ./...
```

This codegen will strip the original files from comments, etc and produce similar but tighter pattern files in folder "languages"

Stripped files are thereafter built with the package as an embedded FS.
