package handler

import "lgdSearch/pkg/trie"

func HotSearch() []trie.HotSearchMessage {
	data := trie.GetHotSearch().Array()
	result := make([]trie.HotSearchMessage, len(data))
	for index, val := range data {
		if val != nil {
			result[index] = *val
		}
	}
	return result
}
