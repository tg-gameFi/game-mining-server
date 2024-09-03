package cmd

import (
	"game-mining-server/app"
	"game-mining-server/utils"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleCmdStart(bot *telego.Bot, update telego.Update) {

	photo := &telego.SendPhotoParams{
		ChatID:  tu.ID(update.Message.Chat.ID),
		Photo:   telego.InputFile{File: utils.MustOpenFile("static/start.png")},
		Caption: "Welcome to CodexField Wallet! 🚀\n\nYour gateway to the world of EVM is now at your fingertips. Easily manage your crypto assets, interact with DeFi protocols, and explore the exciting world of Web3, all within Telegram.\n\nStay tuned for more features and updates!",
		ReplyMarkup: tu.InlineKeyboard(
			tu.InlineKeyboardRow(tu.InlineKeyboardButton("Open Wallet").WithWebApp(&telego.WebAppInfo{URL: app.Config().Bot.WebUrl})),
		),
	}
	_, _ = bot.SendPhoto(photo)
}
