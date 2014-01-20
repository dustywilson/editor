package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var wwwDirPath = "www"
var wwwDir = http.Dir(wwwDirPath)
var wwwFileServer = http.FileServer(wwwDir)
var dataDirPath = "data"
var dataDir = http.Dir(dataDirPath)
var dataFileServer = http.FileServer(dataDir)

func main() {
	var router = mux.NewRouter()
	router.PathPrefix(`/`).Methods("GET").Handler(http.HandlerFunc(ObjectFetch))
	router.PathPrefix(`/`).Methods("POST").Handler(http.HandlerFunc(ObjectUpdate))
	http.Handle("/", router)
	http.ListenAndServe(":1337", nil)
}

func ObjectFetch(w http.ResponseWriter, r *http.Request) {
	_, err := dataDir.Open(r.URL.Path)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := wwwDir.Open(r.URL.Path)
			if err != nil {
				if os.IsNotExist(err) {
					MissingObjectFetch(w, r)
				} else {
					fmt.Fprintf(w, "Error: %s", err)
				}
			} else {
				wwwFileServer.ServeHTTP(w, r)
			}
		} else {
			fmt.Fprintf(w, "Error: %s", err)
		}
	} else {
		dataFileServer.ServeHTTP(w, r)
	}
}

func ObjectUpdate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// FIXME: the following path construction is very insecure!
	err := ioutil.WriteFile(path.Join(dataDirPath, r.URL.Path), []byte(r.PostFormValue("content")), 0644)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		return
	}
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func MissingObjectFetch(w http.ResponseWriter, r *http.Request) {
	// here is a good place to create the missing files (fetch then cache, render, generate, compile, or whatever)
	r.URL.Path = "/post.html"
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	wwwFileServer.ServeHTTP(w, r)
}
