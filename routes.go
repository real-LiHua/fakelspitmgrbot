package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var indexTmpl = template.Must(template.ParseFiles("static/index.html"))

func index(webappURL, namespace string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		challengeCode := rand.Text()
		err := indexTmpl.ExecuteTemplate(writer, "index.html", struct {
			challengecode string
			namespace     string
		}{
			challengecode: challengeCode,
			namespace:     namespace,
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}

func validate(token string) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authQuery, err := url.ParseQuery(r.Header.Get("X-Auth"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("validation failed; failed to parse auth query: " + err.Error()))
		}

		ok, err := ext.ValidateWebAppQuery(authQuery, token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("validation failed; error: " + err.Error()))
			return
		}
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("validation failed; data cannot be trusted."))
			return
		}

		var u gotgbot.User
		err = json.Unmarshal([]byte(authQuery.Get("user")), &u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("validation failed; failed to unmarshal user: " + err.Error()))
			return
		}

		w.Write([]byte(fmt.Sprintf("validation success; user '%s' is authenticated (id: %d).", u.FirstName, u.Id)))
	}
}
