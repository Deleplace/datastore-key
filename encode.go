package datastorekey

import (
	"context"
	"log"
	"os"

	newds "cloud.google.com/go/datastore"
	"google.golang.org/appengine"
	oldds "google.golang.org/appengine/datastore"
)

// See https://developers.google.com/appengine/docs/go/datastore/entities#Go_Kinds_and_identifiers
func CreateOldKey(appID string, namespace string, kind string, stringID string, intID int64, parent *oldds.Key) *oldds.Key {
	c := context.Background()

	if namespace != "" {
		cc, err := appengine.Namespace(c, namespace)
		if err == nil {
			c = cc
		} else {
			log.Printf("Couldn't switch to namespace %q \n", namespace)
		}
	}

	os.Setenv("GAE_LONG_APP_ID", appID)
	os.Setenv("GAE_APPLICATION", appID)

	key := oldds.NewKey(
		c,        // appengine.Context.
		kind,     // Kind.
		stringID, // String ID; empty means no string ID.
		intID,    // Integer ID; if 0, generate automatically. Ignored if string ID specified.
		parent,   // Parent Key; nil means no parent.
	)

	if key.AppID() == "" {
		log.Println("Couldn't find appID :(")
	}

	if namespace != "" && key.Namespace() == "" {
		log.Printf("Couldn't set namespace %q \n", namespace)
	}

	return key
}

func CreateNewKey(namespace string, kind string, stringID string, intID int64, parent *newds.Key) *newds.Key {
	newKey := newds.IncompleteKey(kind, parent)
	newKey.ID = intID
	newKey.Name = stringID
	newKey.Namespace = namespace
	return newKey
}
