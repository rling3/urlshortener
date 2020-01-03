package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var hashes = make(map[string]string)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/health", healthHandler)
	router.HandleFunc("/hash", hashHandler).Methods(http.MethodPost, http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/unhash/{hash}", unhashHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
	router.Use(mux.CORSMethodMiddleware(router))

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		fmt.Println(err)
	}
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside handler")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("got request: " + string(body))
	hashed := string(hash(body))
	fmt.Printf("Putting hash: %s in map with value: %s", string(hashed), string(body))
	hashes[hashed] = string(body)
	w.Write([]byte(hashed))
}

func unhashHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	hash := mux.Vars(r)["hash"]

	fmt.Println("Got hash, extracting URL: " + hash)
	fmt.Printf("Returning value: %s", hashes[hash])
	http.Redirect(w, r, hashes[hash], 302)
}


func hash(url []byte) string {
	h := md5.New()
	h.Write([]byte(url))
	return hex.EncodeToString(h.Sum(nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}