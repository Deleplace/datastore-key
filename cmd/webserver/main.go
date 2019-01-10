package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	newds "cloud.google.com/go/datastore"
	datastorekey "github.com/Deleplace/datastore-key"
	oldds "google.golang.org/appengine/datastore"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

var templates *template.Template

func init() {
	var err error
	templates, err = template.New("datastore-keys").ParseGlob("template/*.html")
	check(err)

	http.HandleFunc("/", index)

	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static/default")))
	http.Handle("/static/", static)
}

func index(w http.ResponseWriter, r *http.Request) {
	c := r.Context()
	data := extractGetParameters(r)

	if oldKeyString := data["oldkeystring"].(string); oldKeyString != "" {
		logf(c, "INFO", "Decoding oldkeystring %q\n", oldKeyString)
		err := autodecodeOld(oldKeyString, data)
		if err == nil {
			logf(c, "INFO", "Decoded %q\n", data)
		} else {
			logf(c, "ERROR", "Failed: %v\n", err.Error())
			// If autodecode failed, render the page with key not decoded
		}

		newKeyString, _, err := datastorekey.ConvertKeyStringForward(oldKeyString)
		if err == nil {
			data["newkeystring"] = newKeyString
		} else {
			logf(c, "ERROR", "Failed to convert old key to new key: %v\n", err.Error())
		}
	} else if newKeyString := data["newkeystring"].(string); newKeyString != "" {
		logf(c, "INFO", "Decoding newkeystring %q\n", newKeyString)
		err := autodecodeNew(newKeyString, data)
		if err == nil {
			logf(c, "INFO", "Decoded %q\n", data)
		} else {
			logf(c, "ERROR", "Failed: %v\n", err.Error())
			// If autodecode failed, render the page with key not decoded
		}
	}

	templates.ExecuteTemplate(w, "index", data)
}

func extractGetParameters(r *http.Request) map[string]interface{} {
	data := map[string]interface{}{
		"kind":         trimmedFormValue(r, "kind"),
		"stringid":     trimmedFormValue(r, "stringid"),
		"intid":        trimmedFormValue(r, "intid"),
		"appid":        trimmedFormValue(r, "appid"),
		"namespace":    trimmedFormValue(r, "namespace"),
		"kind2":        trimmedFormValue(r, "kind2"),
		"stringid2":    trimmedFormValue(r, "stringid2"),
		"intid2":       trimmedFormValue(r, "intid2"),
		"kind3":        trimmedFormValue(r, "kind3"),
		"stringid3":    trimmedFormValue(r, "stringid3"),
		"intid3":       trimmedFormValue(r, "intid3"),
		"oldkeystring": trimmedFormValue(r, "oldkeystring"),
		"newkeystring": trimmedFormValue(r, "newkeystring"),
	}
	return data
}

// IF oldkeystring was given as GET parameter
// THEN it is nice that all decoded values are directly served in the html
func autodecodeOld(oldkeystring string, data map[string]interface{}) error {
	if oldkeystring == "" {
		// Nothing to decode
		return nil
	}
	if data["appid"] != "" || data["kind"] != "" || data["intid"] != "" || data["stringid"] != "" {
		// Don't overwrite user-provided values
		return nil
	}

	key, err := oldds.DecodeKey(oldkeystring)
	if err != nil {
		return err
	}
	fillFieldsOld(key, data)
	return nil
}

func fillFieldsOld(key *oldds.Key, data map[string]interface{}) {
	data["kind"] = key.Kind()
	data["stringid"] = key.StringID()
	data["intid"] = key.IntID()
	data["appid"] = key.AppID()
	data["namespace"] = key.Namespace()
	if key.Parent() != nil {
		data["kind2"] = key.Parent().Kind()
		data["stringid2"] = key.Parent().StringID()
		data["intid2"] = key.Parent().IntID()
		if key.Parent().Parent() != nil {
			data["kind3"] = key.Parent().Parent().Kind()
			data["stringid3"] = key.Parent().Parent().StringID()
			data["intid3"] = key.Parent().Parent().IntID()
		}
	}
}

// IF newkeystring was given as GET parameter
// THEN it is nice that all decoded values are directly served in the html
func autodecodeNew(newkeystring string, data map[string]interface{}) error {
	if newkeystring == "" {
		// Nothing to decode
		return nil
	}
	if data["appid"] != "" || data["kind"] != "" || data["intid"] != "" || data["stringid"] != "" {
		// Don't overwrite user-provided values
		return nil
	}

	key, err := newds.DecodeKey(newkeystring)
	if err != nil {
		return err
	}
	fillFieldsNew(key, data)
	return nil
}

func fillFieldsNew(key *newds.Key, data map[string]interface{}) {
	data["kind"] = key.Kind
	data["stringid"] = key.Name
	data["intid"] = key.ID
	data["namespace"] = key.Namespace
	if key.Parent != nil {
		data["kind2"] = key.Parent.Kind
		data["stringid2"] = key.Parent.Name
		data["intid2"] = key.Parent.ID
		if key.Parent.Parent != nil {
			data["kind3"] = key.Parent.Parent.Kind
			data["stringid3"] = key.Parent.Parent.Name
			data["intid3"] = key.Parent.Parent.ID
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
