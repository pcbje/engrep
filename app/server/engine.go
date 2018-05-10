package server

import (
	"../../engrep"
	"../../automata"
	 "sort"
	 "strings"
)

type Engrep struct {
	engine *engrep.Engrep
	auto *automata.Automata
}

type Entry struct {
	Actual string
	Reference string
	Distance int
	Info string
	Offset int
}


func Build(patterns []string, maxk int) Engrep {
	sort.Strings(patterns)

	auto := automata.CreateAutomata(patterns)

	trie := engrep.CreateEngrep(maxk, true, engrep.CreateDawg(maxk))

	trie.AddReferences(patterns)

	return Engrep{
		engine: trie,
		auto: auto,
	}
}

func (s Engrep) Search(text string, k int) []Entry {
	stop := map[string]bool{
		".":	true,
		",":	true,
		"?":	true,
		" ":	true,
		"!":	true,
		"\"":	true,
		"'":	true,
		"-":	true,
		"#":	true,
		"\n":	true,
		"\r":	true,
		"[":	true,
		"]":	true,
		"{":	true,
		"}":	true,
		">":	true,
		"<":	true,
		"\x00":	true,
	}

	prev := ""
	res := []Entry{}

	cache := map[string][]Entry{}

	s.engine.Scan(text, k, func(z int, e int, actual string, pre []rune, suf []rune, d int) {
		actual = strings.TrimSpace(actual)

		validPre := len(pre) == 0
		validSuf := len(suf) == 0

		if !validPre {
			for _, s := range pre {
					_, v := stop[string(s)]
					validPre = validPre || v
					if validPre {break}
			}
		}

		if !validSuf {
			for _, s := range suf {
					_, v := stop[string(s)]
					validSuf = validSuf || v
					if validSuf {break}
			}
		}

		if !validPre || !validSuf {
			//println(validPre, validSuf, len(suf), string(pre), string(suf))
			return
		}



		if actual != prev {
			if _, ok := cache[actual]; !ok {
				x := []Entry{}

				for _, found := range s.auto.FindAll([]rune(actual), k) {
					x = append(x, Entry{
						Actual: actual,
						Reference: found.Match,
						Distance: found.Error,
						Offset: z,
						//Info: s.database[found.Match],
					})
				}

				cache[actual] = x
			}

			//println(actual)

			prev = actual
			res = append(res, cache[actual]...)
		}
	})

	return res
}
