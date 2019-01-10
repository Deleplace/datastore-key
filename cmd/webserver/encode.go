package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	datastorekey "github.com/Deleplace/datastore-key"
	"google.golang.org/appengine/datastore"
)

func init() {
	http.HandleFunc("/encode", ajaxEncode)
}

func ajaxEncode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	c := r.Context()

	kind := trimmedFormValue(r, "kind")
	stringID := trimmedFormValue(r, "stringid")
	intIDstr := trimmedFormValue(r, "intid")
	appID := trimmedFormValue(r, "appid")
	namespace := trimmedFormValue(r, "namespace")

	// Parent (optional)
	kind2 := trimmedFormValue(r, "kind2")
	stringID2 := trimmedFormValue(r, "stringid2")
	intIDstr2 := trimmedFormValue(r, "intid2")

	// Grand-parent (optional)
	kind3 := trimmedFormValue(r, "kind3")
	stringID3 := trimmedFormValue(r, "stringid3")
	intIDstr3 := trimmedFormValue(r, "intid3")

	logf(c, "INFO", "Encoding %v\n", []string{
		appID, namespace,
		kind, stringID, intIDstr,
		kind2, stringID2, intIDstr2,
		kind3, stringID3, intIDstr3,
	})

	var key, parent, grandparent *datastore.Key

	if kind2 != "" {
		if kind3 != "" {
			grandparent = datastorekey.CreateKey(appID, namespace, kind3, stringID3, intID64(intIDstr3), nil)
		}
		parent = datastorekey.CreateKey(appID, namespace, kind2, stringID2, intID64(intIDstr2), grandparent)
	}

	key = datastorekey.CreateKey(appID, namespace, kind, stringID, intID64(intIDstr), parent)
	//fmt.Fprint(w, keyString)
	w.Header().Set("Content-Type", "application/json")
	keyString := key.Encode()
	fmt.Fprint(w, Response{
		"keystring": keyString,
	})
	logf(c, "INFO", "Encoded %v\n", keyString)
}

func intID64(intIDstr string) int64 {
	if intIDstr == "" {
		return 0
	}
	intID64, _ := strconv.ParseInt(intIDstr, 10, 64)
	return intID64
}

func trimmedFormValue(r *http.Request, paramName string) string {
	return strings.TrimSpace(r.FormValue(paramName))
}
