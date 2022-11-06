package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type shortyType string

const (
	redirect shortyType = "r"
	file                = "f"
)

type shortyReq struct {
	sType  shortyType `json:"sType"`
	target string     `json:"target"`
}

var shortyFilesDir = "./shortyFilesDir/"

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	handlerMap := make(map[shortyType]func(http.ResponseWriter, *http.Request, string))
	handlerMap[redirect] = redirectRequest
	handlerMap[file] = sendFile

	reqShorty := a.store.getShorty(r.URL.Path[1:])
	f := func(w http.ResponseWriter, r *http.Request, s string) { notFound(w, r, s) }

	if reqShorty.target != "" {
		f = handlerMap[reqShorty.sType]
	}
	f(w, r, reqShorty.target)
}

func redirectRequest(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 301)
}

func notFound(w http.ResponseWriter, r *http.Request, some string) {
	http.NotFound(w, r)
}

func sendFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(shortyFilesDir + filePath)
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)
	w.Header().Set("Content-Type", "application/octetstream")

	io.Copy(w, file)
	return
}

func (a *app) addShorty(w http.ResponseWriter, r *http.Request) {
	// /add/{f|r}/{URL}/{target}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	d := json.NewDecoder(r.Body)
	newShorty := make(map[string]string)
	err := d.Decode(&newShorty)
	if err != nil {
		a.logFatal(err.Error())
	}
	sType := newShorty["sType"]
	shortcut := newShorty["shortcut"]
	target := newShorty["target"]
	a.store.addShorty(shortcut, &shortyReq{sType: shortyType(sType), target: target})
	w.Write([]byte("Shorty Added!"))
	return
}

func main() {
	addr := flag.String("addr", ":4444", "HTTP Network Address")
	flag.Parse()

	a := &app{}
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/add", a.addShorty)
	fmt.Printf("Starting service on %s\n", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
