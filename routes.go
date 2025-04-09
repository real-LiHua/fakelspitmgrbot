package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"text/template"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var indexTmpl = template.Must(template.ParseFiles("static/index.html"))

func (bot *Bot) webappValidateAndGetUserID(writer http.ResponseWriter, request *http.Request) int64 {
	response := map[string]string{}
	errorMessage := "Validation failed, please try again later\n验证失败，请稍后重试"
	writer.Header().Set("Content-Type", "application/json")
	authQuery, err := url.ParseQuery(request.Header.Get("X-Auth"))
	if err != nil {
		response["message"] = errorMessage
		json.NewEncoder(writer).Encode(response)
		return 0
	}
	ok, err := ext.ValidateWebAppQuery(authQuery, bot.self.Token)
	if err != nil || !ok {
		response["message"] = errorMessage
		json.NewEncoder(writer).Encode(response)
		return 0
	}
	var u gotgbot.User
	err = json.Unmarshal([]byte(authQuery.Get("user")), &u)
	if err != nil {
		response["message"] = errorMessage
		json.NewEncoder(writer).Encode(response)
		return 0
	}
	return u.Id
}

func (bot *Bot) webappIndex(writer http.ResponseWriter, request *http.Request) {
	userID := bot.webappValidateAndGetUserID(writer, request)
	challengeCode := bot.db.UpdateChallengeCode(userID)
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
	bot.webappValidateAndGetUserID(writer, request)
}
