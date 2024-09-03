package main

import (
	"fmt"
	"game-mining-server/app"
	"game-mining-server/handlers"
	"game-mining-server/routers"
	"log"
	"os"
)

func main() {
	// os.Args[0] always is the executable file, so our args starts with os.Args[1]
	args := os.Args
	if len(args) != 2 && len(args) != 3 {
		log.Println("Please specify env and botToken, Usage: ./entry dev botToken")
		return
	}
	env := args[1]
	botToken := ""
	if len(args) == 3 {
		botToken = args[2]
	}

	// create app and config
	if e0 := app.CreateApp(fmt.Sprintf("config.%s.json", env), botToken); e0 != nil {
		panic(fmt.Errorf("create app failed: %s", e0))
	}

	if e2 := handlers.RegisterBotAndRun(app.Bot(), app.Config().Bot); e2 != nil {
		panic(fmt.Errorf("bot server run failed: %s", e2))
	}

	// start http server
	if e3 := routers.InitAndRun(app.Config().Basic); e3 != nil {
		panic(fmt.Errorf("http server run failed: %s", e3))
	}
}
