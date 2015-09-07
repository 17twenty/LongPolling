package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var messages chan string = make(chan string, 100)

func PushHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("push Request received")

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, "+
		"Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, "+
		"Cache-Control, X-Requested-With")
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		w.WriteHeader(400)
	}

	messages <- string(body)
}

func PollResponse(w http.ResponseWriter, req *http.Request) {
	log.Println("poll Request received")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, "+
		"Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, "+
		"Cache-Control, X-Requested-With")

	io.WriteString(w, <-messages) // Blocks if no message! We call this the LOOONG POLL
}

func main() {
	http.HandleFunc("/poll", PollResponse)
	http.HandleFunc("/push", PushHandler)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
