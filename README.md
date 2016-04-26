# Telebot 4 Google App Engine [![Build Status](https://travis-ci.org/cortinico/telebotgae.svg?branch=master)](https://travis-ci.org/cortinico/telebotgae) [![GoDoc](https://godoc.org/github.com/cortinico/telebotgae?status.svg)](https://godoc.org/github.com/cortinico/telebotgae)

A simple Telegram bot skeleton written in Go (GAE Capable)

This library is derived from [telebot](https://github.com/cortinico/telebot) and allows you to deploy your
bot on [Google App Engine](https://appengine.google.com). With this you can simply deploy your bot code on AppEngine and forget about server management and scaling.

## Setup

* [Download, Setup and Configure the Go App Engine SDK](https://cloud.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go). Don't forget to set the `$PATH` environment variable as explained in Google guide.

* We assume that **Go** and the `$GOPATH` variable are properly set up. If not, please follow [getting started with Go](https://golang.org/doc/install)

* Create a new project on [Google Cloud Dashboard](http://cloud.google.com) and remember the **project id** that Google gives to you, say `my-telegrambot-project-id`.

<img src="http://i.imgur.com/lvyjpCN.png" width=500px />

* Create a new bot interacting with [@BotFather](https://telegram.me/botfather) on Telegram. Open a chat with @BotFather and start asking with `/newbot`. He will guide you through the creation of a new bot, and he will give you an **API Key**, say `162227600:AAAAAAAAAAABBBBBBBBBBCCCCCCCCCDDDDD`.

<img src="http://i.imgur.com/fgnwUgE.png" width=500px />

* Copy the **example bot skeleton** inside the `bot_example` folder of this repository.
```bash
git clone https://github.com/cortinico/telebotgae.git && cd telebotgae/bot_example
```

* Edit the file `hello.go` adding your **bot name** (without @) and your **API Key**
```go
func init() {
	conf := telebotgae.Configuration{
		BotName: "MyNewSampleBot",
		ApiKey:  "162227600:AAAAAAAAAAABBBBBBBBBBCCCCCCCCCDDDDD"}
```

* Edit the file `app.yaml` adding your **project-id** from Google Cloud Dashboard
```yaml
application: my-telegrambot-project-id
version: 1
runtime: go
api_version: go1

handlers:
- url: /.*
  script: _go_app
  secure: always
```

* Grab this library and build the project
```bash
goapp get github.com/cortinico/telebotgae && goapp build
```

* Deploy your bot to App Engine. If it's the first time you deploy, you will be asked for **Google authentication**.
```bash
goapp deploy
```

* Visit the following web page in a web browser:
```
https://api.telegram.org/bot[API_KEY]/setWebhook?url=https://[PROJECT-ID].appspot.com
```

So your URL should look like this:
```
https://api.telegram.org/bot162227600:AAAAAAAAAAABBBBBBBBBBCCCCCCCCCDDDDD/setWebhook?url=https://my-telegrambot-project-id.appspot.com
```
**DON'T FORGET TO DO IT, AND DON'T MISPELL, OTHERWISE YOUR BOT WON'T WORK**

If you see this message:
```json
{"ok":true,"result":true,"description":"Webhook was set"}
```
![browser](http://i.imgur.com/TIwT19v.png)

Then your bot is working :D Have fun with telebotgae

## hello.go

Your bot code should look like this:
```go
package hello

import (
	"github.com/cortinico/telebotgae"
	"net/http"
)

func init() {
	conf := telebotgae.Configuration{
		BotName: "MyNewSampleBot",
		ApiKey:  "162227600:AAAAAAAAAAABBBBBBBBBBCCCCCCCCCDDDDD"}

	var bot telebotgae.Bot

	bot.Startgae(conf, func(mess string, r *http.Request) (string, error) {
		return "You typed " + mess, nil
	})
}
```

You can use the second parameter of `Startgae` to implement the logic of your bot.
The second parameter must be a function with the following type:
```go
type Responder func(string, *http.Request) (string, error)
```
You will receive in input a `string` with the message from the user, and an `http.Request`.
You can use the `http.Request` to get access to all the App Engine nice feature such as 
Datastore, Memcache, etc.

You have to provide a tuple made by the answer and the error. Set the error to nil if
nothing unexpected has occurred.

## Configuration

`Configuration` can also be loaded from a JSON file, using the `LoadSettings(filename string) (Configuration, error)`
function.

## Licence

The following software is released under the [MIT Licence](https://github.com/cortinico/telebot/blob/master/LICENSE)
