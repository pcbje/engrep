package server

import (
	"../engrep"
	"../automata"
	 "sort"
	 "strings"
)

type Entry struct {
	Actual string
	Reference string
	Distance int
	Info string
}

type Server struct {
	engine *engrep.Engrep
	auto *automata.Automata
}

func Build(patterns []string, maxk int) Server {
	sort.Strings(patterns)

	auto := automata.CreateAutomata(patterns)

	trie := engrep.CreateEngrep(maxk, true, engrep.CreateDawg(maxk))

	trie.AddReferences(patterns)

	return Server{
		engine: trie,
		auto: auto,
	}
}

func (s Server) Search(text string, k int) []Entry {
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
	}

	prev := ""
	res := []Entry{}

	s.engine.Scan(" "+text+" ", func(z int, e int, actual string, pre string, suf string, d int) {
		actual = strings.TrimSpace(actual)
		_, validPre := stop[pre]
		_, validSuf := stop[suf]

		if !validPre || !validSuf {
			return
		}

		if actual != prev {

			for _, found := range s.auto.FindAll(actual, k) {
				entry := Entry{
					Actual: actual,
					Reference: found.Match,
					Distance: found.Error,
					//Info: s.database[found.Match],
				}

				prev = actual
				res = append(res, entry)
			}
		}
	})

	return res
}
