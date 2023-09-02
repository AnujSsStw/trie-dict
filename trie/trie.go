package trie

import (
	"strings"
)

type Node struct {
	val      string
	end      bool
	children [26]*Node
}

type Trie struct {
	RootNode *Node
}

func (t *Trie) InsertText(text string) {
	zero := []rune("a")[0]
	text = strings.ToLower(text)
	curr := t.RootNode

	for _, e := range text {
		idx := e - zero
		if curr.children[idx] == nil {
			curr.children[idx] = &Node{val: string(rune(e))}
		}
		curr = curr.children[idx]
		curr.end = false
	}
	curr.end = true
}

func Autocomplete(node *Node, prefix string, sugg *[]string) {
	if node == nil {
		return
	}
	// if len(*sugg) >= 10 {
	// 	return
	// }

	if node.end {
		*sugg = append(*sugg, prefix)
	}

	for i, child := range node.children {
		if child != nil {
			su := string(rune('a' + i))
			Autocomplete(child, prefix+su, sugg)
		}
	}

}
func printAutoSuggestions(root *Node, text string) []string {
	pCrawl := root
	text = strings.ToLower(text)
	sugg := []string{}

	for _, e := range text {
		idx := e - 'a'
		if pCrawl.children[idx] != nil {
			pCrawl = pCrawl.children[idx]
		} else {
			return sugg
		}
	}
	Autocomplete(pCrawl, text, &sugg)
	return sugg
}

func Trieee(autoCum string, a *Trie) []string {
	var ans []string
	if len(autoCum) > 3 {
		ans = printAutoSuggestions(a.RootNode, autoCum)
	}

	return ans
}
