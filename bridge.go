package datastorekey

import (
	newds "cloud.google.com/go/datastore"
	oldds "google.golang.org/appengine/datastore"
)

// Compare all contents EXCEPT a.AppID() which is ignored.
func Compare(a *oldds.Key, b *newds.Key) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Note the AppID is NOT part of b

	if a.Kind() != b.Kind {
		return false
	}
	if a.IntID() != b.ID {
		return false
	}
	if a.StringID() != b.Name {
		return false
	}
	if a.Namespace() != b.Namespace {
		return false
	}

	return Compare(a.Parent(), b.Parent)
}

// ConvertForward transform an old-libs *Key into a new-libs *Key.
// The new-libs *Key does not contain an AppID.
func ConvertForward(old *oldds.Key) (new_ *newds.Key, appID string) {
	if old == nil {
		return nil, ""
	}
	appID = old.AppID()
	newParent, _ := ConvertForward(old.Parent())
	new_ = newds.IncompleteKey(old.Kind(), newParent)
	new_.ID = old.IntID()
	new_.Name = old.StringID()
	new_.Namespace = old.Namespace()
	return new_, appID
}

func ConvertBackward(new_ *newds.Key, appID string) (old *oldds.Key) {
	if new_ == nil {
		return nil
	}
	oldParent := ConvertBackward(new_.Parent, appID)
	old = CreateKey(appID, new_.Namespace, new_.Kind, new_.Name, new_.ID, oldParent)
	return old
}
