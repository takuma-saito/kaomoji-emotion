
package main
import ("fmt")

type TrieTree struct {
	parent *TrieTree
	path rune
	children map[rune]*TrieTree
}
var root TrieTree

// Trie 木を作成
func MakeTrie(tree *TrieTree, name string) {
	if tree == nil {panic("error occured")}
	if tree.children == nil {
		tree.children = make(map[rune]*TrieTree)
	}
	if len(name) == 0 {
		tree.children[TERMINAL] = nil
	} else {
		char, size := utf8.DecodeRuneInString(name)
		_, ok := tree.children[char]
		if !ok {
			tree.children[char] = new(TrieTree)
			tree.children[char].parent = tree
			tree.children[char].path = char
		}
		MakeTrie(tree.children[char], name[size:])
	}
}

// 単純な Cost 距離
func Cost(x string, y string) {
	
}

func main() {
	
}



