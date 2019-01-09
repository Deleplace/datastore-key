package datastorekey

import (
	"encoding/json"

	"google.golang.org/appengine/datastore"
)

func JsonifyKey(key *datastore.Key) (s string) {
	b, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func RecursiveJson(key *datastore.Key) Response {
	var parentJson Response
	if key.Parent() != nil {
		parentJson = RecursiveJson(key.Parent())
	}
	return Response{
		"stringID":  key.StringID(),
		"intID":     key.IntID(),
		"kind":      key.Kind(),
		"appID":     key.AppID(),
		"namespace": key.Namespace(),
		"parent":    parentJson,
	}
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
