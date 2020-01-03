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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	hashed := string(hash(body))
	hashes[hashed] = string(body)
	w.Write([]byte(hashed))
}

func unhashHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	hash := mux.Vars(r)["hash"]

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