package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jamesorlakin/cacheyd/pkg/mux"
	"github.com/jamesorlakin/cacheyd/pkg/service"
)

type JsonableRequest struct {
	Method        string
	URL           url.URL
	Proto         string
	Header        http.Header
	ContentLength int64
	Host          string
	RemoteAddr    string
	RequestURI    string
}

func main() {
	log.Println("Starting cacheyd")

	port := 3000
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		portInt, err := strconv.Atoi(portEnv)
		if err == nil {
			port = portInt
		}
	}
	listenStr := ":" + strconv.Itoa(port)
	log.Printf("Listening on %s", listenStr)

	router := mux.NewRouter(&service.CacheydService{})

	everything := func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		router.ServeHTTP(w, r)
	}

	err := http.ListenAndServe(listenStr, http.HandlerFunc(everything))
	if err != nil {
		log.Panic(err)
	}
}

func logRequest(r *http.Request) {
	// http.Request contains methods which the JSON marshaller doesn't like.
	jsonRequest := JsonableRequest{
		Method:        r.Method,
		URL:           *r.URL,
		Proto:         r.Proto,
		Header:        r.Header,
		ContentLength: r.ContentLength,
		Host:          r.Host,
		RemoteAddr:    r.RemoteAddr,
		RequestURI:    r.RequestURI,
	}
	log.Printf("Received HTTP request: %+v\n", jsonRequest)
	_, err := json.MarshalIndent(jsonRequest, "", "\t")
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}

	// log.Print(requestJson)
}
