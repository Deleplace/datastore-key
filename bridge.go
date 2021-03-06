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

// ConvertKeyForward transform an old-libs *Key into a new-libs *Key.
// The new-libs *Key does not contain an AppID.
func ConvertKeyForward(old *oldds.Key) (new_ *newds.Key, appID string) {
	if old == nil {
		return nil, ""
	}
	appID = old.AppID()
	newParent, _ := ConvertKeyForward(old.Parent())
	new_ = newds.IncompleteKey(old.Kind(), newParent)
	new_.ID = old.IntID()
	new_.Name = old.StringID()
	new_.Namespace = old.Namespace()
	return new_, appID
}

func ConvertKeyBackward(new_ *newds.Key, appID string) (old *oldds.Key) {
	if new_ == nil {
		return nil
	}
	oldParent := ConvertKeyBackward(new_.Parent, appID)
	old = CreateOldKey(appID, new_.Namespace, new_.Kind, new_.Name, new_.ID, oldParent)
	return old
}

func ConvertKeyStringForward(oldstr string) (newstr string, appID string, err error) {
	oldKey, err := oldds.DecodeKey(oldstr)
	if err != nil {
		return "", "", err
	}
	newKey, appID := ConvertKeyForward(oldKey)
	newstr = newKey.Encode()
	return newstr, appID, nil
}

func ConvertKeyStringBackward(newstr string, appID string) (oldstr string, err error) {
	newKey, err := newds.DecodeKey(newstr)
	if err != nil {
		return "", err
	}
	oldKey := ConvertKeyBackward(newKey, appID)
	oldstr = oldKey.Encode()
	return oldstr, nil
}
