package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	datastorekey "github.com/Deleplace/datastore-key"
	"google.golang.org/appengine/datastore"
)

func init() {
	http.HandleFunc("/decode", ajaxDecode)
}

func ajaxDecode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	keyString := trimmedFormValue(r, "keystring")
	c := r.Context()
	logf(c, "INFO", "Decoding %v\n", keyString)

	key, err := datastore.DecodeKey(keyString)
	if err != nil {
		logf(c, "ERROR", "Failed: %v\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := datastorekey.RecursiveJson(key).String()
	logf(c, "INFO", "%v\n", response)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, response)
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}
