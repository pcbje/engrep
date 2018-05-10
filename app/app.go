package main

import (
  "os"
  "./server"
	"strconv"
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func GetPatterns(path string) []string {
	var patterns map[string]string

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}

	json.Unmarshal(bytes, &patterns)

	ps := []string{}

	for p, _ := range patterns {
		if len(p) > 11 {
			ps = append(ps, p)
		}
	}

	return ps
}

func main() {
	listen := os.Args[1]
	names := GetPatterns(os.Args[2])
	maxk, err := strconv.Atoi(os.Args[3])

	f, err := os.OpenFile("engrep-app.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	if err != nil {
		log.Panic(err)
	}

	s := server.CreateServer(names, maxk)

	http.Handle("/search", s.Limit(http.HandlerFunc(s.Search)))
	http.Handle("/create", s.Limit(http.HandlerFunc(s.Create)))
	http.Handle("/", server.NoCache(http.FileServer(http.Dir("./app/client"))))

	log.Print(fmt.Sprintf("Listening on: %v", listen))

	if len(os.Args) == 6 {
		// 4: /etc/letsencrypt/live/www.yourdomain.com/fullchain.pem
		// 5: /etc/letsencrypt/live/www.yourdomain.com/privkey.pem
		http.ListenAndServeTLS(listen, os.Args[4], os.Args[5], nil)
	} else {
		http.ListenAndServe(listen, nil)
	}
}
