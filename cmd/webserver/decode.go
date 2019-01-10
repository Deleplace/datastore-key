package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	newds "cloud.google.com/go/datastore"
	datastorekey "github.com/Deleplace/datastore-key"
	oldds "google.golang.org/appengine/datastore"
)

func init() {
	http.HandleFunc("/decode", ajaxDecode)
}

func ajaxDecode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	oldKeyString := trimmedFormValue(r, "oldkeystring")
	newKeyString := trimmedFormValue(r, "newkeystring")

	if (oldKeyString == "") == (newKeyString == "") {
		http.Error(w, "Please provide either oldkeystring or newkeystring to decode", http.StatusBadRequest)
		return
	}

	c := r.Context()
	logf(c, "INFO", "Decoding %v\n", oldKeyString)
	var response string

	if oldKeyString != "" {
		oldKey, err := oldds.DecodeKey(oldKeyString)
		if err != nil {
			logf(c, "ERROR", "Failed: %v\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		responseMap := datastorekey.RecursiveJsonOld(oldKey)

		newKeyString, _, err = datastorekey.ConvertKeyStringForward(oldKeyString)
		if err == nil {
			responseMap["newkeystring"] = newKeyString
		} else {
			logf(c, "ERROR", "Failed to convert old key to new key: %v\n", err.Error())
		}

		response = responseMap.String()
	} else if newKeyString != "" {
		newKey, err := newds.DecodeKey(newKeyString)
		if err != nil {
			logf(c, "ERROR", "Failed: %v\n", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response = datastorekey.RecursiveJsonNew(newKey).String()
	}

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
