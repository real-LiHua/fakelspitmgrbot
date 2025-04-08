package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/chatjoinrequest"
	_ "modernc.org/sqlite"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("TOKEN environment variable is empty")
	}
	webappURL := os.Getenv("URL")
	if webappURL == "" {
		panic("URL environment variable is empty")
	}

	webhookSecret := os.Getenv("WEBHOOK_SECRET")
	if webhookSecret == "" {
		panic("WEBHOOK_SECRET environment variable is empty")
	}

	chatID := os.Getenv("CHAT_ID")
	if chatID == "" {
		panic("CHAT_ID environment variable is empty")
	}
	parsedChatID, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		panic("CHAT_ID environment variable is not a valid integer: " + err.Error())
	}

	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		panic("NAMESPACE environment variable is empty")
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)
	dispatcher.AddHandler(handlers.NewCommand("dl", dl))
	dispatcher.AddHandler(handlers.NewCommand("dl_debug", dl_debug))
	dispatcher.AddHandler(handlers.NewCommand("ban", ban))
	dispatcher.AddHandler(handlers.NewCommand("unban", unban))
	dispatcher.AddHandler(handlers.NewCommand("ban_github", ban_github))
	dispatcher.AddHandler(handlers.NewCommand("unban_github", unban_github))
	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		c := strings.SplitN(ctx.Message.Text, " ", 3)
		if len(c) == 2 {
			switch c[1] {
			case "dl":
				return dl(b, ctx)
			case "dl_debug":
				return dl_debug(b, ctx)
			}
		}
		return nil
	}))
	dispatcher.AddHandler(handlers.NewChatJoinRequest(chatjoinrequest.ChatID(parsedChatID),
		func(b *gotgbot.Bot, ctx *ext.Context) error {
			return JoinRequest(b, ctx, webappURL)
		}))

	err = updater.AddWebhook(b, b.Token, &ext.AddWebhookOpts{SecretToken: webhookSecret})
	if err != nil {
		panic("Failed to add bot webhooks to updater: " + err.Error())
	}

	updaterSubpath := "/" + rand.Text() + "/"
	err = updater.SetAllBotWebhooks(webappURL+updaterSubpath, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        webhookSecret,
	})
	if err != nil {
		panic("Failed to set bot webhooks: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", index(webappURL, namespace))
	mux.HandleFunc("/validate", validate(token))
	mux.HandleFunc("/submit", submit(namespace))
	mux.HandleFunc(updaterSubpath, updater.GetHandlerFunc(updaterSubpath))

	server := http.Server{
		Handler: mux,
		Addr:    listenAddr,
	}

	log.Printf("%s has been started...\n", b.User.Username)
	if err := server.ListenAndServe(); err != nil {
		panic("failed to listen and serve: " + err.Error())
	}
}
