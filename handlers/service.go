package handlers

import (
	"fmt"
	"game-mining-server/configs"
	"game-mining-server/handlers/cmd"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"log"
)

func RegisterBotAndRun(bot *telego.Bot, config *configs.BotConfig) error {
	if bot == nil {
		log.Println("Not register bot handlers")
		return nil
	}
	if botUser, e1 := bot.GetMe(); e1 != nil {
		return fmt.Errorf("Register bot handlers get info failed: " + e1.Error())
	} else {
		log.Printf("Register bot(%s) handlers\n", botUser.Username)
	}
	updates, e0 := bot.UpdatesViaLongPolling(nil)
	if e0 != nil {
		return e0
	}
	bh, e1 := th.NewBotHandler(bot, updates)
	if e1 != nil {
		return e1
	}
	if e2 := createBotMenu(bot, config); e2 != nil {
		return e2
	}

	registerHandlers(bh)
	go bh.Start()
	return nil
}

func createBotMenu(bot *telego.Bot, config *configs.BotConfig) error {
	e0 := bot.SetChatMenuButton(&telego.SetChatMenuButtonParams{
		MenuButton: &telego.MenuButtonWebApp{
			Type: "web_app",
			Text: "💰Wallet",
			WebApp: telego.WebAppInfo{
				URL: config.WebUrl,
			},
		},
	})
	if e0 != nil {
		return fmt.Errorf("Create bot menu button failed: " + e0.Error())
	} else {
		log.Println("Create bot menu button success")
	}

	e1 := bot.SetMyCommands(&telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{
				Command:     "start",
				Description: "start",
			},
		},
	})
	if e1 != nil {
		return fmt.Errorf("Create bot set commands failed: " + e1.Error())
	}
	log.Println("Create bot set commands success")
	return nil
}

func registerHandlers(bh *th.BotHandler) {
	bh.Handle(cmd.HandleCmdStart, th.CommandEqual("start"))
}
