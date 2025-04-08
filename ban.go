package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) commandBan(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, err := check(b, ctx)
	if err != nil {
		return nil
	}

	ctx.EffectiveChat.BanMember(b, userID, nil)
	bot.db.BanUser(bot.db.GetUserByTelegramID(userID))

	return nil
}

func (bot *Bot) commandBanGitHub(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, err := check(b, ctx)
	if err != nil {
		return nil
	}

	user := bot.db.GetUserByGithubID(userID)
	ctx.EffectiveChat.BanMember(b, user.TelegramID, nil)
	bot.db.BanUser(user)

	return nil
}
