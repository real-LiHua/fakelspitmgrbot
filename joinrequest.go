package main

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func JoinRequest(b *gotgbot.Bot, ctx *ext.Context, webappURL string) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Hello, I'm @%s.\nYou can use me to run a (very) simple telegram webapp demo!", b.User.Username), &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "Press me", WebApp: &gotgbot.WebAppInfo{Url: webappURL}},
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}
