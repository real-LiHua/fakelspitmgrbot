package main

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (bot *Bot) commandUnban(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, _, err := check(b, ctx)
	if err != nil {
		return nil
	}
	_, err = b.UnbanChatMember(bot.chatID, userID, nil)
	if err != nil {
		return err
	}
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("User %d unbanned", userID), nil)
	bot.db.UnbanUser(bot.db.GetUserByTelegramID(userID))
	return nil
}

func (bot *Bot) commandUnbanGitHub(b *gotgbot.Bot, ctx *ext.Context) error {
	userID, _, err := check(b, ctx)
	if err != nil {
		return nil
	}

	user := bot.db.GetUserByGithubID(userID)
	_, err = b.UnbanChatMember(bot.chatID, user.TelegramID, nil)
	if err != nil {
		return err
	}
	ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Github user %d unbanned", userID), nil)
	bot.db.UnbanUser(user)
	return nil
}
