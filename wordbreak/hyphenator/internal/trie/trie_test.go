package trie

import (
	"testing"
)

func TestRuneTrie(t *testing.T) {
	trie := NewRuneTrie()
	testTrie(t, trie)
}

func TestRuneTrieNilBehavior(t *testing.T) {
	trie := NewRuneTrie()
	testNilBehavior(t, trie)
}

func TestRuneTrieRoot(t *testing.T) {
	trie := NewRuneTrie()
	testTrieRoot(t, trie)

	trie = NewRuneTrie()
	if !trie.isLeaf() {
		t.Error("root of empty tree should be leaf")
	}
	trie.Put([]rune(""), "root")
	if !trie.isLeaf() {
		t.Error("root should not have children, only value")
	}
}

func testTrie(t *testing.T, trie Trier) {
	const firstPutValue = "first put"
	cases := []struct {
		key   []rune
		value interface{}
	}{
		{[]rune("fish"), 0},
		{[]rune("/cat"), 1},
		{[]rune("/dog"), 2},
		{[]rune("/cats"), 3},
		{[]rune("/caterpillar"), 4},
		{[]rune("/cat/gideon"), 5},
		{[]rune("/cat/giddy"), 6},
	}

	// get missing keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %v to be missing, found value %v", c.key, value)
		}
	}

	// initial put
	for _, c := range cases {
		if isNew := trie.Put(c.key, firstPutValue); !isNew {
			t.Errorf("expected key %v to be missing", c.key)
		}
	}

	// subsequent put
	for _, c := range cases {
		if isNew := trie.Put(c.key, c.value); isNew {
			t.Errorf("expected key %v to have a value already", c.key)
		}
	}

	// get
	for _, c := range cases {
		if value := trie.Get(c.key); value != c.value {
			t.Errorf("expected key %v to have value %v, got %v", c.key, c.value, value)
		}
	}

	// delete, expect Delete to return true indicating a node was nil'd
	for _, c := range cases {
		if deleted := trie.Delete(c.key); !deleted {
			t.Errorf("expected key %v to be deleted", c.key)
		}
	}

	// delete cleaned all the way to the first character
	// expect Delete to return false bc no node existed to nil
	for _, c := range cases {
		if deleted := trie.Delete([]rune{c.key[0]}); deleted {
			t.Errorf("expected key %v to be cleaned by delete", string(c.key[0]))
		}
	}

	// get deleted keys
	for _, c := range cases {
		if value := trie.Get(c.key); value != nil {
			t.Errorf("expected key %v to be deleted, got value %v", c.key, value)
		}
	}
}

func testNilBehavior(t *testing.T, trie Trier) {
	cases := []struct {
		key   []rune
		value interface{}
	}{
		{[]rune("/cat"), 1},
		{[]rune("/catamaran"), 2},
		{[]rune("/caterpillar"), nil},
	}
	expectNilValues := [][]rune{[]rune("/"), []rune("/c"), []rune("/ca"), []rune("/caterpillar"), []rune("/other")}

	// initial put
	for _, c := range cases {
		if isNew := trie.Put(c.key, c.value); !isNew {
			t.Errorf("expected key %v to be missing", c.key)
		}
	}

	// get nil
	for _, key := range expectNilValues {
		if value := trie.Get(key); value != nil {
			t.Errorf("expected key %v to have value nil, got %v", key, value)
		}
	}
}

func testTrieRoot(t *testing.T, trie Trier) {
	const firstPutValue = "first put"
	const putValue = "value"

	if value := trie.Get([]rune("")); value != nil {
		t.Errorf("expected key '' to be missing, found value %v", value)
	}
	if !trie.Put([]rune(""), firstPutValue) {
		t.Error("expected key '' to be missing")
	}
	if trie.Put([]rune(""), putValue) {
		t.Error("expected key '' to have a value already")
	}
	if value := trie.Get([]rune("")); value != putValue {
		t.Errorf("expected key '' to have value %v, got %v", putValue, value)
	}
	if !trie.Delete([]rune("")) {
		t.Error("expected key '' to be deleted")
	}
	if value := trie.Get([]rune("")); value != nil {
		t.Errorf("expected key '' to be deleted, got value %v", value)
	}
}
