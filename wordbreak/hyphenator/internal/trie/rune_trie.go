package trie

// RuneTrie is a trie of runes with []rune keys and interface{} values.
//
// Internal nodes have nil values: a stored nil value will thus not be distinguishable.
type RuneTrie struct {
	value    interface{}
	children map[rune]*RuneTrie
}

// NewRuneTrie allocates and returns a new RuneTrie.
func NewRuneTrie() *RuneTrie {
	return new(RuneTrie)
}

// Get returns the value stored at the given key.
//
// Returns nil for internal nodes or for nodes with a value of nil.
func (trie *RuneTrie) Get(key []rune) interface{} {
	node := trie
	for _, r := range key {
		node = node.children[r]
		if node == nil {
			return nil
		}
	}

	return node.value
}

// Put inserts the value into the trie at the given key, replacing any
// existing items.
//
// It returns true if the put adds a new value, and false if it replaces an existing value.
func (trie *RuneTrie) Put(key []rune, value interface{}) bool {
	node := trie
	for _, r := range key {
		child := node.children[r]
		if child == nil {
			if node.children == nil {
				node.children = map[rune]*RuneTrie{}
			}
			child = new(RuneTrie)
			node.children[r] = child
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.value == nil
	node.value = value

	return isNewVal
}

// Delete removes the value associated with the given key.
//
// Returns true if a node was found for the given key.
// If the node or any of its ancestors becomes childless as a result, it is removed from the trie.
func (trie *RuneTrie) Delete(key []rune) bool {
	path := make([]nodeRune, len(key)) // record ancestors to check later
	node := trie
	for i, r := range key {
		path[i] = nodeRune{r: r, node: node}
		node = node.children[r]
		if node == nil {
			// node does not exist
			return false
		}
	}
	// delete the node value
	node.value = nil
	// if leaf, remove it from its parent's children map. Repeat for ancestor
	// path.
	if node.isLeaf() {
		// iterate backwards over path
		for i := len(key) - 1; i >= 0; i-- {
			parent := path[i].node
			r := path[i].r
			delete(parent.children, r)
			if !parent.isLeaf() {
				// parent has other children, stop
				break
			}
			parent.children = nil
			if parent.value != nil {
				// parent has a value, stop
				break
			}
		}
	}

	return true // node (internal or not) existed and its value was nil'd
}

// A node of the RuneTrie with its the rune key and child to descend into.
type nodeRune struct {
	node *RuneTrie
	r    rune
}

func (trie *RuneTrie) isLeaf() bool {
	return len(trie.children) == 0
}
