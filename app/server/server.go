package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
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

type Dictionary struct {
	server   Engrep
	last     time.Time
	patterns int
}

type Server struct {
	maxk    int
	engine  map[string]*Dictionary
	timeout map[string]time.Time
}

type Response struct {
	Results  []Entry `json:"results"`
	Took     string  `json:"took"`
	Error    string  `json:"error"`
	Patterns int     `json:"patterns"`
}

func (s Server) Create(w http.ResponseWriter, r *http.Request) {
  var patterns []string

	reader := bufio.NewReader(r.Body)

	max_patterns := 10000

	bytes := make([]byte, max_patterns*64)
	rr, _ := reader.Read(bytes)

  raw_patterns := strings.Split(string(bytes[0:rr]), "\n")

  for _, pattern := range raw_patterns {
    if len(strings.TrimSpace(pattern)) > 0 {
      patterns = append(patterns, pattern)
    }
  }

	if len(patterns) == 0 {
		log.Panic("Min 1 pattern")
	}

	if len(patterns) > max_patterns {
		log.Panic("Max 10000 patterns")
	}

	letters := []rune("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	ok := true
	key := ""

	for ok {
		key = RandStringRunes(8, letters)
		_, ok = s.engine[key]
	}

	maxkstr := r.URL.Query().Get("k")
	maxk, err := strconv.Atoi(maxkstr)

	if err != nil {
		log.Panic("k is not a number")
	}

	if maxk > 2 {
		log.Panic("max k=2")
	}

	s.engine[key] = &Dictionary{server: Build(patterns, maxk), last: time.Now(), patterns: len(patterns)}

	s.log(r, fmt.Sprintf("created: %s, max k: %s, patterns: %d", key, maxkstr, len(patterns)))

	fmt.Fprintf(w, key)
}

func (s Server) log(r *http.Request, message string) {
	ip := strings.Split(r.RemoteAddr, ":")[0]

	log.Print(fmt.Sprintf("[%s] %s", ip, message))
}

func (s Server) Search(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	engine := r.URL.Query().Get("d")

	if engine == "" {
		engine = "demo"
	}

	if _, ok := s.engine[engine]; !ok {
		log.Panic("not found")
	}

	kstr := r.URL.Query().Get("k")

	k := 99
	if kstr == "" {
		k = s.maxk
	}
	k, err := strconv.Atoi(kstr)

	if err != nil {
		log.Panic("k is not a number")
	}

	if k > s.maxk {
		log.Panic("k is bigger than ", s.maxk)
	}

	reader := bufio.NewReader(r.Body)

	bytes := make([]byte, 1024*1024)

	rr, _ := reader.Read(bytes)

	bytes = bytes[0:rr]

	text := "        " + string(bytes) + "         "

	results := s.engine[engine].server.Search(text, k)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	response := Response{
		Results:  results,
		Error:    "",
		Took:     fmt.Sprintf("%s", time.Since(start)),
		Patterns: s.engine[engine].patterns,
	}
	jsonBytes, _ := json.MarshalIndent(response, "", "  ")

	s.engine[engine].last = start

	s.log(r, fmt.Sprintf("searched. engine: %s, k: %s, patterns: %d, len: %d took: %s", engine, kstr, s.engine[engine].patterns, len(text), time.Since(start)))

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	fmt.Fprintf(w, string(jsonBytes))
}

func (s Server) removeInactive() {
	for true {
		time.Sleep(60 * time.Second)
		var m runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m)
		log.Print("Allocated memory: ", m.Alloc)

		rem := 0
		for key, engine := range s.engine {
			if key == "demo" {
				continue
			}

			// Idle for an hour
			if time.Since(engine.last).Seconds() > 60*60 {
				delete(s.engine, key)
				log.Print("Removed ", key)
				rem++
			}
		}

		if rem > 0 {
			runtime.GC()
		}
	}
}

func (s Server) Limit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				fmt.Println("Error:", recovered)
				http.Error(w, fmt.Sprintf("Argh: %v", recovered), http.StatusInternalServerError)
			}
		}()

		start := time.Now()

		ip := strings.Split(r.RemoteAddr, ":")[0]

		if prev, ok := s.timeout[ip]; ok {
			if time.Since(prev).Seconds() < 0.0001 {
				s.timeout[ip] = start
				s.log(r, "too frequent")

				response := Response{
					Took:  fmt.Sprintf("%s", time.Since(start)),
					Error: "Max one request per 2 seonds...",
				}
				jsonBytes, _ := json.MarshalIndent(response, "", "  ")

				fmt.Fprintf(w, string(jsonBytes))
				return
			}
		}

		s.timeout[ip] = start

		h.ServeHTTP(w, r)
	})
}

func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Accept-Charset", "UTF-8")
		h.ServeHTTP(w, r)
	})
}

func CreateServer(names []string, maxk int) Server {
	s := Server{
		engine:  map[string]*Dictionary{"demo": &Dictionary{server: Build(names, maxk), patterns: len(names)}},
		maxk:    maxk,
		timeout: map[string]time.Time{},
	}

	go s.removeInactive()

	return s
}
