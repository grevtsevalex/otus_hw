package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(s string) []string {
	words := strings.Fields(s)

	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}

	type word struct {
		Word  string
		Count int
	}

	wds := make([]word, 0, len(counts))
	for k, v := range counts {
		wds = append(wds, word{Word: k, Count: v})
	}

	sort.Slice(wds, func(i, j int) bool {
		if wds[i].Count == wds[j].Count {
			return wds[i].Word < wds[j].Word
		}
		return wds[i].Count > wds[j].Count
	})

	res := make([]string, 0, 10)
	for i, v := range wds {
		if i == 10 {
			break
		}
		res = append(res, v.Word)
	}

	return res
}
