package gonuts

import (
	"strings"
)

// TrieNode represents a node in the Trie data structure.
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	value    interface{} // This can be used to store additional information at each node
}

// Trie is a tree-like data structure for efficient string operations.
type Trie struct {
	root *TrieNode
}

// NewTrie creates and returns a new Trie.
//
// Example:
//
//	trie := NewTrie()
func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}

// Insert adds a word to the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
func (t *Trie) Insert(word string) {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// InsertWithValue adds a word to the Trie with an associated value.
//
// Example:
//
//	trie := NewTrie()
//	trie.InsertWithValue("apple", 42)
func (t *Trie) InsertWithValue(word string, value interface{}) {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
	node.value = value
}

// BulkInsert efficiently inserts multiple words into the Trie.
//
// Example:
//
//	trie := NewTrie()
//	words := []string{"apple", "app", "application"}
//	trie.BulkInsert(words)
func (t *Trie) BulkInsert(words []string) {
	for _, word := range words {
		t.Insert(word)
	}
}

// Search checks if a word exists in the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	fmt.Println(trie.Search("apple"))  // Output: true
//	fmt.Println(trie.Search("app"))    // Output: false
func (t *Trie) Search(word string) bool {
	node := t.findNode(word)
	return node != nil && node.isEnd
}

// StartsWith checks if any word in the Trie starts with the given prefix.
//
// Example:
//
//	trie := NewTrie()
//	trie.Insert("apple")
//	fmt.Println(trie.StartsWith("app"))  // Output: true
//	fmt.Println(trie.StartsWith("ban"))  // Output: false
func (t *Trie) StartsWith(prefix string) bool {
	return t.findNode(prefix) != nil
}

// findNode is a helper function to find a node for a given word or prefix.
func (t *Trie) findNode(word string) *TrieNode {
	node := t.root
	for _, ch := range word {
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}

// AutoComplete returns a list of words that start with the given prefix.
//
// Example:
//
//	trie := NewTrie()
//	trie.BulkInsert([]string{"apple", "app", "application", "appreciate"})
//	suggestions := trie.AutoComplete("app", 3)
//	fmt.Println(suggestions)  // Output: [app apple application]
func (t *Trie) AutoComplete(prefix string, limit int) []string {
	node := t.findNode(prefix)
	if node == nil {
		return []string{}
	}

	result := []string{}
	t.dfs(node, prefix, &result, limit)
	return result
}

// dfs is a helper function for AutoComplete.
func (t *Trie) dfs(node *TrieNode, prefix string, result *[]string, limit int) {
	if len(*result) == limit {
		return
	}
	if node.isEnd {
		*result = append(*result, prefix)
	}
	for ch, child := range node.children {
		t.dfs(child, prefix+string(ch), result, limit)
	}
}

// WildcardSearch searches for words matching a pattern with wildcards.
// The '.' character in the pattern matches any single character.
//
// Example:
//
//	trie := NewTrie()
//	trie.BulkInsert([]string{"cat", "dog", "rat"})
//	matches := trie.WildcardSearch("r.t")
//	fmt.Println(matches)  // Output: [rat]
func (t *Trie) WildcardSearch(pattern string) []string {
	result := []string{}
	t.wildcardDfs(t.root, "", pattern, &result)
	return result
}

// wildcardDfs is a helper function for WildcardSearch.
func (t *Trie) wildcardDfs(node *TrieNode, current, pattern string, result *[]string) {
	if len(current) == len(pattern) {
		if node.isEnd {
			*result = append(*result, current)
		}
		return
	}

	if pattern[len(current)] == '.' {
		for ch, child := range node.children {
			t.wildcardDfs(child, current+string(ch), pattern, result)
		}
	} else {
		ch := rune(pattern[len(current)])
		if child, ok := node.children[ch]; ok {
			t.wildcardDfs(child, current+string(ch), pattern, result)
		}
	}
}

// LongestCommonPrefix finds the longest common prefix of all words in the Trie.
//
// Example:
//
//	trie := NewTrie()
//	trie.BulkInsert([]string{"flower", "flow", "flight"})
//	fmt.Println(trie.LongestCommonPrefix())  // Output: "fl"
func (t *Trie) LongestCommonPrefix() string {
	if t.root == nil {
		return ""
	}

	var sb strings.Builder
	node := t.root

	for len(node.children) == 1 {
		var ch rune
		for k := range node.children {
			ch = k
			break
		}
		sb.WriteRune(ch)
		node = node.children[ch]
		if node.isEnd {
			break
		}
	}

	return sb.String()
}
