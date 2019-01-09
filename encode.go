package datastorekey

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/appengine/datastore"
)

// See https://developers.google.com/appengine/docs/go/datastore/entities#Go_Kinds_and_identifiers
func CreateKey(appID string, namespace string, kind string, stringID string, intID int64, parent *datastore.Key) (*datastore.Key, error) {
	// c is the true context of the current request
	// forged is a wrapper context with our custom appID
	// forged := &ForgedContext{c, appID}
	// cc is a wrapper context with our custom namespace
	// cc, err := appengine.Namespace(forged, namespace)
	// if err != nil {
	// 	return nil, err
	// }
	// cc := c // TODO what about Namespace?
	// var cc context.Context
	cc := context.Background()

	os.Setenv("GAE_LONG_APP_ID", appID)
	os.Setenv("GAE_APPLICATION", appID)

	key := datastore.NewKey(
		cc,       // appengine.Context.
		kind,     // Kind.
		stringID, // String ID; empty means no string ID.
		intID,    // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		parent,   // Parent Key; nil means no parent.
	)
	return key, nil
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
