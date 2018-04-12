package main

import (
	"./engrep"
	"fmt"
	"sort"
	"math"
	"math/rand"
	"time"
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
	success := 0
	errors := 0
	tests := 100000
	allchars := []rune("012346789abcdef")
	for x := 0; x < tests; x++ {
		for y := 2; y < len(allchars); y++ {
			chars := allchars[0:y]
			patts := 5
			max_len := 12
			text_len := 20
			stringcount := 1

			experiment := GenerateExperiment(chars, patts, max_len, text_len, stringcount, true)

			for k := 1; k < 4; k++ {
				for _, str := range experiment.Texts {
					trie := engrep.CreateEngrep(k, true, engrep.CreateDawg(k))
					trie.AddReferences(experiment.Patterns)

					z := []string{}
					c := map[string]bool{}

					trie.Scan(str+" ", max_len, func(s int, e int, actual string, reference string, d int) {
						if _, ok := c[reference]; !ok {
							c[reference] = true

							z = append(z, reference)
						}
					})

					y := []string{}
					yy := []string{}

					for x, d := range experiment.Distances[str] {
						found := d > k
						for _, zz := range z {
							if found || Lev(zz, x) <= k {
								found = true
								break
							}

							if found || Lev(zz, x[1:]) <= k {
								found = true
								break
							}

							if found || Lev(zz, x[2:]) <= k {
								found = true
								break
							}

							if found || Lev(zz, x[:len(x)-1]) <= k {
								found = true
								break
							}

							if found || Lev(zz, x[:len(x)-2]) <= k {
								found = true
								break
							}
						}

						if d <= k && !found {
							y = append(y, x)
						}
						if d <= k {
							yy = append(yy, x)
						}
					}

					sort.Strings(y)

					if len(y) > 0 {
						fmt.Println(k, yy, y, z, str, experiment.Patterns, len(c))
						errors += 1
					} else {
						success += 1
					}
				}
			}
		}

		print(fmt.Sprintf("\r%.1f%% completed. %v tests passed, %v failed", float64(x*100)/float64(tests), success, errors))
	}

	if errors == 0 {
		println("\npassed")
	} else {
		println("\nfailed")
	}

}
