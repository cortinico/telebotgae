# Telebotgae [![Build Status](https://travis-ci.org/cortinico/telebotgae.svg?branch=master)](https://travis-ci.org/cortinico/telebotgae)

A simple Telegram bot skeleton written in Go (GAE Capable)

This library is derived from [telebot](https://github.com/cortinico/telebot) and allows you to deploy your
bot on [Google App Engine](https://appengine.google.com)

## Usage

You simply need a configuration (BotName + API Key + ProjectID) and a Response function.

Checkout this sample code:
```go
package hello

import (
	"github.com/cortinico/telebotgae"
	"net/http"
)

func init() {
    conf := telebotgae.Configuration{
        BotName: "SampleBot",
        ApiKey:  "162227600:AAAAAAAAAAABBBBBBBBBBCCCCCCCCCDDDDD",
        ProjID: "mysimple-telegram-bot"}

    var bot telebotgae.Bot

    bot.Startgae(conf, func(mess string, req *http.Request)
         (string, error) {
        var answer string
		switch mess {
		case "/test":
			answer = "Test command works :)"
		default:
			answer = "You typed " + mess
		}
		return answer, nil
    })
}
```

## Licence

The following software is released under the [MIT Licence](https://github.com/cortinico/telebot/blob/master/LICENSE)
