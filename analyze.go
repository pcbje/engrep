package main

import (
	"./engrep"
	"math"
	"math/rand"
	"time"
	"fmt"
	"encoding/json"
)

func init() {
	if true {
		rand.Seed(time.Now().UnixNano())
	}
}

func RandStringRunes(n int, letterRunes []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Lev(a_str string, b_str string) int {
	matrix := [][]int{}

	if len(b_str) == 0 {
		return len(a_str)
	}

	if len(a_str) == 0 {
		return len(b_str)
	}

	for i := 0; i <= len(a_str)+1; i++ {
		row := []int{}
		for j := 0; j <= len(b_str)+1; j++ {
			if i == 0 {
				row = append(row, j)
			} else if j == 0 {
				row = append(row, i)
			} else {
				row = append(row, 0)
			}
		}

		matrix = append(matrix, row)
	}

	d := 0
	for i := 1; i <= len(a_str); i++ {
		for j := 1; j <= len(b_str); j++ {
			if a_str[i-1] == b_str[j-1] {
				d = matrix[i-1][j-1]
				matrix[i][j] = d
			} else {
				d = int(math.Min(
					float64(matrix[i-1][j-1]),
					math.Min(float64(matrix[i][j-1]),
						float64(matrix[i-1][j])),
				)) + 1

				matrix[i][j] = d
			}
		}
	}

	return matrix[len(a_str)][len(b_str)]
}

type Experiment struct {
	Patterns  []string
	Texts     []string
	Distances map[string]map[string]int
}

func (e Experiment) Verify(k int, str string, hits []string, patterns []string) bool {

	return false
}

func GenerateExperiment(chars []rune, patts int, max_len int, text_len int, stringcount int, fullsubstring bool) Experiment {
	patterns := []string{}
	strings := []string{}

	for i := 0; i < patts; i++ {
		patterns = append(patterns, RandStringRunes(max_len+1, chars))
	}

	for i := 0; i < stringcount; i++ {
		strings = append(strings, RandStringRunes(text_len, chars))
	}

	cache := map[string]map[string]int{}

	for _, str := range strings {
		cache[str] = map[string]int{}

		for i := max_len; i <= len(str); i++ {
			x := 0
			y := i

			if !fullsubstring {
				x = i - max_len
				y = i - max_len + 1
			}

			for j := x; j < y; j++ {
				substring := str[j:i]
				for _, pattern := range patterns {

					d := Lev(substring, pattern)

					key := substring

					if d < 4 {
						if _, ok := cache[str][key]; !ok {
							cache[str][key] = 99
						}

						if d < cache[str][key] {
							cache[str][key] = d
						}
					}
				}
			}
		}
	}

	return Experiment{
		Patterns:  patterns,
		Texts:     strings,
		Distances: cache,
	}
}

func main() {
	tests := 1
	allchars := []rune("0123456789abcdefghijklnmnopqrstuvwxyz")
	max_patterns := 1000
	max_len := 20
	text_lens := 100
	fmt.Println("k", "pattern_length", "text_len", "alphabet_size", "patterns", "average_active_states", "states_with_error")
	for patts := 0; patts <= max_patterns; patts += 100 {
		println(patts)
		for x := 0; x < tests; x++ {
			for y := 12; y <= 24; y+=6 {
				for text_len := 100; text_len <= text_lens; text_len+=100 {
					chars := allchars[0:y]
					stringcount := 1

					for k := 1; k < 4; k++ {
						for j:=0; j < 3; j++ {
							experiment := GenerateExperiment(chars, patts+1, max_len, text_len, stringcount, true)

							for _, str := range experiment.Texts {
								trie := engrep.CreateEngrep(k, true, engrep.CreateDawg(k))
								trie.AddReferences(experiment.Patterns)
								active_states, hits, errors := trie.Scan(str+" ", max_len, func(s int, e int, actual string, reference string, d int) {})
								jsonString0, _ := json.Marshal(active_states)
								jsonString, _ := json.Marshal(hits)
								jsonString2, _ := json.Marshal(errors)
								xpatts := patts
								if xpatts == 0 {
									xpatts = 1
								}
								fmt.Println(k, max_len, text_len, y, xpatts, string(jsonString0), string(jsonString), string(jsonString2))
							}
						}
					}
				}
			}
		}
	}
}
