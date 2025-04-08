package main

import (
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func check(b *gotgbot.Bot, ctx *ext.Context) (int64, error) {
	c := strings.SplitN(ctx.Message.Text, " ", 3)
	if len(c) < 2 {
		return 0, nil
	}

	sender, err := ctx.EffectiveChat.GetMember(b, ctx.EffectiveUser.Id, nil)

	if err != nil {
		return 0, err
	}

	if !sender.MergeChatMember().CanRestrictMembers {
		return 0, err
	}

	userID, err := strconv.ParseInt(c[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
