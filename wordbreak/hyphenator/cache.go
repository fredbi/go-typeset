package hyphenator

import "sync"

var (
	mx sync.Mutex

	// A cache for dictionaries. Users of this package only pay the cost
	// of building the trie once for a given language.
	loadedPatterns map[string]*Dictionary
)

// load a preloaded trie Dictionary for some language patterns file.
func loadDictFromCache(patterns string) *Dictionary {
	mx.Lock()
	defer mx.Unlock()

	if loadedPatterns == nil {
		loadedPatterns = make(map[string]*Dictionary)
	}

	dict, ok := loadedPatterns[patterns]
	if !ok {
		dict, _ = LoadPatterns(patterns)
	}

	loadedPatterns[patterns] = dict

	return dict
}
