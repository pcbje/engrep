package main

import (
	"./server"
	"fmt"
	"os"
	"strconv"
	"bufio"
	"io/ioutil"
	"time"
	 "encoding/json"
	 "log"
	 "net/http"
	 "strings"
	 "math/rand"
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

func GetPatterns(path string) []string {
	var patterns map[string]string

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	json.Unmarshal(bytes, &patterns)

	ps := []string{}

	for p, _ := range patterns {
		ps = append(ps, p)
	}

	return ps
}

type Engine struct {
	server server.Server
	last time.Time
}

type Server struct {
	maxk int
	engine map[string]*Engine
	timeout map[string]time.Time
}

type Response struct {
	Results []server.Entry `json:"results"`
	Took string `json:"took"`
}

func (s Server) create(w http.ResponseWriter, r *http.Request) {
	var patterns []string

	reader := bufio.NewReader(r.Body)

	bytes := make([]byte, 1024 * 1024)
	rr, _ := reader.Read(bytes)

	err := json.Unmarshal(bytes[0:rr], &patterns)

	if err != nil {
		log.Panic("Could not decode json object:", err)
	}

	if len(patterns) > 10000 {
		log.Panic("Max 10000 patterns")
	}

	letters := []rune("abcdefghijklmnopqrstuvxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	key := RandStringRunes(16, letters)

	_, ok := s.engine[key]

	for ok {
		key = RandStringRunes(16, letters)
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

	s.engine[key] = &Engine{server: server.Build(patterns, maxk), last: time.Now()}

	fmt.Fprintf(w, key)
}

func (s Server) search(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	engine := r.URL.Query().Get("e")

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

	bytes := make([]byte, 1024 * 10)
	reader.Read(bytes)

	text := string(bytes)

	results := s.engine[engine].server.Search(text, k)

	response := Response{Results: results, Took: fmt.Sprintf("%s", time.Since(start))}
	jsonBytes, _ := json.MarshalIndent(response, "", "  ")

	s.engine[engine].last = start

	fmt.Fprintf(w, string(jsonBytes))
}

func (s Server) removeInactive() {
	for true {
		time.Sleep(60 * time.Second)
		for key, engine := range s.engine {
			if key == "demo" {
				continue
			}

			// Idle for an hour
			if time.Since(engine.last).Seconds() > 60*60 {
				delete(s.engine, key)
				log.Print("Removed ", key)
			}
		}
	}
}

func (s Server) Limit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				fmt.Println("Error:", recovered)
				http.Error(w, fmt.Sprintf("Aaaaaaaargh (%v)", recovered), http.StatusInternalServerError)
			}
		}()

		start := time.Now()

		ip := strings.Split(r.RemoteAddr, ":")[0]

		if prev, ok := s.timeout[ip]; ok {
			if time.Since(prev).Seconds() < 2.0 {
				s.timeout[ip] = start
				log.Panic("Max one request every two seconds")
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
		h.ServeHTTP(w, r)
	})
}

func main() {
	listen := os.Args[1]
	names := GetPatterns(os.Args[2])
	maxk, err := strconv.Atoi(os.Args[3])

	if err != nil {
		log.Panic(err)
	}

	s := Server{
		engine: map[string]*Engine{"demo": &Engine{server: server.Build(names, maxk)}},
		maxk: maxk,
		timeout: map[string]time.Time{},
	}

	http.Handle("/search", s.Limit(http.HandlerFunc(s.search)))
	http.Handle("/create", s.Limit(http.HandlerFunc(s.create)))
	http.Handle("/", NoCache(http.FileServer(http.Dir("./server/client"))))

	log.Print(fmt.Sprintf("Listening on: %v", listen))

	go s.removeInactive()

	http.ListenAndServe(listen, nil)
}
