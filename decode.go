package datastorekey

import (
	"encoding/json"

	newds "cloud.google.com/go/datastore"
	oldds "google.golang.org/appengine/datastore"
)

func JsonifyKey(key *oldds.Key) (s string) {
	b, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func RecursiveJsonOld(key *oldds.Key) Response {
	var parentJson Response
	if key.Parent() != nil {
		parentJson = RecursiveJsonOld(key.Parent())
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

func RecursiveJsonNew(key *newds.Key) Response {
	var parentJson Response
	if key.Parent != nil {
		parentJson = RecursiveJsonNew(key.Parent)
	}
	return Response{
		"stringID":  key.Name,
		"intID":     key.ID,
		"kind":      key.Kind,
		"namespace": key.Namespace,
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
