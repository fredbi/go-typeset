package trie

// Trier exposes the Trie structure capabilities.
type Trier interface {
	Get(key []rune) interface{}
	Put(key []rune, value interface{}) bool
	Delete(key []rune) bool
}
