package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/chatjoinrequest"
	"github.com/joho/godotenv"
)

type Bot struct {
	self      *gotgbot.Bot
	db        *Database
	namespace string
	chatID    int64
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
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

	db, err := NewDatabase()
	if err != nil {
		panic("failed to create new database: " + err.Error())
	}

	b, err := gotgbot.NewBot(token, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	bot := Bot{
		self:      b,
		db:        db,
		namespace: namespace,
		chatID:    parsedChatID,
	}
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)
	dispatcher.AddHandler(handlers.NewCommand("dl", bot.commandDownload))
	dispatcher.AddHandler(handlers.NewCommand("dl_debug", bot.commandDownloadDebug))
	dispatcher.AddHandler(handlers.NewCommand("ban", bot.commandBan))
	dispatcher.AddHandler(handlers.NewCommand("unban", bot.commandUnban))
	dispatcher.AddHandler(handlers.NewCommand("ban_github", bot.commandBanGitHub))
	dispatcher.AddHandler(handlers.NewCommand("unban_github", bot.commandUnbanGitHub))
	dispatcher.AddHandler(handlers.NewCommand("start", func(b *gotgbot.Bot, ctx *ext.Context) error {
		c := strings.SplitN(ctx.Message.Text, " ", 3)
		if len(c) == 2 {
			switch c[1] {
			case "dl":
				return bot.commandDownload(b, ctx)
			case "dl_debug":
				return bot.commandDownloadDebug(b, ctx)
			}
		}
		return nil
	}))
	dispatcher.AddHandler(handlers.NewChatJoinRequest(chatjoinrequest.ChatID(bot.chatID),
		func(b *gotgbot.Bot, ctx *ext.Context) error {
			_, err = b.SendMessage(ctx.ChatJoinRequest.UserChatId,
				"Please complete the verification\n请完成验证",
				&gotgbot.SendMessageOpts{
					ReplyMarkup: gotgbot.InlineKeyboardMarkup{
						InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
							{Text: "Verify 验证", WebApp: &gotgbot.WebAppInfo{Url: webappURL}},
						}},
					},
				})
			if err != nil {
				return fmt.Errorf("failed to send start message: %w", err)
			}
			return nil
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
	mux.HandleFunc("/", bot.webappIndex)
	mux.HandleFunc("/validate", bot.webappValidate)
	mux.HandleFunc("/submit", bot.webappSubmit)
	mux.HandleFunc("/css/main.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/css/main.css")
	})
	mux.HandleFunc("/js/telegram-web-app.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/js/telegram-web-app.js")
	})
	mux.HandleFunc("/js/main.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/js/main.js")
	})
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
