package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) commandUnban(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, err := check(b, ctx)
	if err != nil {
		return nil
	}
	ctx.EffectiveChat.UnbanMember(b, userID, nil)
	bot.db.UnbanUser(bot.db.GetUserByTelegramID(userID))
	return nil
}

func (bot *Bot) commandUnbanGitHub(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, err := check(b, ctx)
	if err != nil {
		return nil
	}

	user := bot.db.GetUserByGithubID(userID)
	ctx.EffectiveChat.UnbanMember(b, user.TelegramID, nil)
	bot.db.UnbanUser(user)

	return nil
}
