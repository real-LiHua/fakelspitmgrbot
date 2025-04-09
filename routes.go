package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var indexTmpl = template.Must(template.ParseFiles("static/index.html"))

func (bot *Bot) webappGetUserID(writer http.ResponseWriter, request *http.Request) int64 {
	authQuery, err := url.ParseQuery(request.Header.Get("X-Auth"))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("failed to parse auth query: " + err.Error()))
		return 0
	}
	ok, err := ext.ValidateWebAppQuery(authQuery, bot.self.Token)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("validation failed; error: " + err.Error()))
		return 0
	}
	if !ok {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("validation failed; data cannot be trusted."))
		return 0
	}
	var u gotgbot.User
	err = json.Unmarshal([]byte(authQuery.Get("user")), &u)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("failed to unmarshal user: " + err.Error()))
		return 0
	}
	return u.Id
}

func (bot *Bot) webappIndex(writer http.ResponseWriter, request *http.Request) {
	userID := bot.webappGetUserID(writer, request)
	challengeCode := bot.db.GetChallengeCode(userID)
	err := indexTmpl.ExecuteTemplate(writer, "index.html", struct {
		challengecode string
		namespace     string
	}{
		challengecode: challengeCode,
		namespace:     bot.namespace,
	})
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
	}
}

func (bot *Bot) webappValidate(writer http.ResponseWriter, request *http.Request) {
	userID := bot.webappGetUserID(writer, request)
	if userID == 0 {
		return
	}
	writer.Write([]byte("validation success"))
}
