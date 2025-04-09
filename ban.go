package main

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) commandBan(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, reason, err := check(b, ctx)
	if err != nil {
		return nil
	}

	_, err = b.BanChatMember(bot.chatID, userID, nil)
	if err != nil {
		return err
	}
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("User %d banned for reason:\n%s", userID, reason), nil)
	bot.db.BanUser(bot.db.GetUserByTelegramID(userID))
	return nil
}

func (bot *Bot) commandBanGitHub(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, reason, err := check(b, ctx)
	if err != nil {
		return nil
	}

	user := bot.db.GetUserByGithubID(userID)
	_, err = b.BanChatMember(bot.chatID, user.TelegramID, nil)
	if err != nil {
		return err
	}
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Github user %d banned for reason:\n%s", userID, reason), nil)
	bot.db.BanUser(user)
	return nil
}
