package hello

import (
	"github.com/cortinico/telebotgae"
	"net/http"
)

func init() {
	conf := telebotgae.Configuration{
		BotName: "PLACE-HERE-YOUR-BOT-NAME",
		ApiKey:  "PLACE-HERE-YOUR-API-KEY"}

	var bot telebotgae.Bot

	bot.Startgae(conf, func(mess string, r *http.Request) (string, error) {
		return "You typed " + mess, nil
	})
}
