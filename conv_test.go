package datastorekey

import (
	"testing"

	newds "cloud.google.com/go/datastore"
	oldds "google.golang.org/appengine/datastore"
)

func TestOldKeyString(t *testing.T) {
	// Valid keys
	for _, oldWebsafeKey := range []string{
		"ahRzfnByb2dyYW1taW5nLWlkaW9tc3IMCxIFSWRpb20YmgEM",
	} {
		oldKey, err := oldds.DecodeKey(oldWebsafeKey)
		if err != nil {
			t.Errorf("Could not decode %q: %v", oldWebsafeKey, err)
			continue
		}
		oldKeyStr := oldKey.Encode()
		if oldKeyStr != oldWebsafeKey {
			t.Errorf("Expected %q, got %q", oldWebsafeKey, oldKeyStr)
		}

		newKey, appID := ConvertKeyForward(oldKey)
		if appID != oldKey.AppID() {
			t.Errorf("Expected %q, got %q", oldKey.AppID(), appID)
		}
		if !Compare(oldKey, newKey) {
			t.Errorf("Old and New keys should be the same, but are not")
		}
		newKeyStr := newKey.Encode()
		newKey2, err := newds.DecodeKey(newKeyStr) // newKey2 ~ newKey
		if err != nil {
			t.Errorf("Error decoding newKeyStr: %v", err)
		}
		oldKey2 := ConvertKeyBackward(newKey2, appID)
		oldKeyStr2 := oldKey2.Encode()
		if oldKeyStr2 != oldKeyStr {
			t.Errorf("Expected %q, got %q", oldKeyStr, oldKeyStr2)
		}
	}
}
