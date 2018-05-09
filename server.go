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
)

func GetPatterns(path string) map[string]string {
	var patterns map[string]string

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	json.Unmarshal(bytes, &patterns)

	return patterns
}

type Server struct {
	maxk int
	engine server.Server
	timeout map[string]time.Time
}

type Response struct {
	Results []server.Entry `json:"results"`
	Took string `json:"took"`
}

func (s Server) search(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	defer func() {
		recovered := recover()
		if recovered != nil {
			fmt.Println("Error:", recovered)
			http.Error(w, fmt.Sprintf("Aaaaaaaa (%v)", recovered), http.StatusInternalServerError)
		}
	}()

	ip := strings.Split(r.RemoteAddr, ":")[0]

	if prev, ok := s.timeout[ip]; ok {
		if time.Since(prev).Seconds() < 2.0 {
			s.timeout[ip] = start
			log.Panic("Max one request every two seconds")
		}
	}

	s.timeout[ip] = start

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

	results := s.engine.Search(text, k)

	response := Response{Results: results, Took: fmt.Sprintf("%s", time.Since(start))}
	jsonBytes, _ := json.Marshal(response)
	fmt.Fprintf(w, string(jsonBytes))
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
		engine: server.Build(names, maxk),
		maxk: maxk,
		timeout: map[string]time.Time{},
	}

	http.Handle("/search", http.HandlerFunc(s.search))
	http.Handle("/", NoCache(http.FileServer(http.Dir("./server/client"))))

	log.Print(fmt.Sprintf("Listening on: %v", listen))

	http.ListenAndServe(listen, nil)

	/*reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')


	res := s.Search(text)

	for _, r := range res {
		println(r.Found, "->", r.Reference, "=", r.Error)
	}*/
}
