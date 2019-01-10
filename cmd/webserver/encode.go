package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	newds "cloud.google.com/go/datastore"
	datastorekey "github.com/Deleplace/datastore-key"
	oldds "google.golang.org/appengine/datastore"
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

	// Old keys contain an AppID, new keys don't.
	// Producing an old key requires an AppID.
	hasAppID := appID != ""

	logf(c, "INFO", "Encoding %v\n", []string{
		appID, namespace,
		kind, stringID, intIDstr,
		kind2, stringID2, intIDstr2,
		kind3, stringID3, intIDstr3,
	})

	var oldKey, oldParent, oldGrandparent *oldds.Key
	oldKeyString := ""

	if hasAppID {
		if kind2 != "" {
			if kind3 != "" {
				oldGrandparent = datastorekey.CreateOldKey(appID, namespace, kind3, stringID3, intID64(intIDstr3), nil)
			}
			oldParent = datastorekey.CreateOldKey(appID, namespace, kind2, stringID2, intID64(intIDstr2), oldGrandparent)
		}
		oldKey = datastorekey.CreateOldKey(appID, namespace, kind, stringID, intID64(intIDstr), oldParent)
		oldKeyString = oldKey.Encode()
	}

	var newKey, newParent, newGrandparent *newds.Key
	if kind2 != "" {
		if kind3 != "" {
			newGrandparent = datastorekey.CreateNewKey(namespace, kind3, stringID3, intID64(intIDstr3), nil)
		}
		newParent = datastorekey.CreateNewKey(namespace, kind2, stringID2, intID64(intIDstr2), newGrandparent)
	}
	newKey = datastorekey.CreateNewKey(namespace, kind, stringID, intID64(intIDstr), newParent)
	newKeyString := newKey.Encode()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{
		"oldKeyString": oldKeyString,
		"newKeyString": newKeyString,
	})
	logf(c, "INFO", "Encoded %v\n", oldKeyString)
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
