/*
Package trie implements a trie with []rune key.

The Tries do not synchronize access (not thread-safe). A typical use case is
to perform Puts and Deletes upfront to populate the Trie, then perform Gets
very quickly.
*/
package trie
