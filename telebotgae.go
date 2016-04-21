// Package for creating a running a simple Telegram bot on Google
// App Engine. This bot is capable just to answer simple user/group
// messages, all the logic must be implemented inside a Responder func
package telebotgae

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Interface Repserent a generic telegram bot. Exported functions
// are just LoadSettings to load a configuration and Start to
// launch the bot.
type IBot interface {
	Start(conf Configuration, resp Responder)
	LoadSettings(filename string) (Configuration, error)
	telegramSendURL(conf Configuration) string
}

// Struct representing a telegram Bot (will implement IBot).
// Bot has no field (no state), it's just an empty bot
type Bot struct{}

// Responder function, responsible of handling to user commands.
// This function represent the logic of your bot, you must provide
// a couple (string, error) for every message. The returned string
// will be sent to the user. If you set the error, the user will
// see an informative message.
// TODO
type Responder func(string, *http.Request) (string, error)

// Configuration struct representing the configuration used from
// the bot to run properly. Configuration is usually loaded from file,
// or hardcoded inside the client code.
type Configuration struct {
	BotName string `json:"BotName"` // Name of the bot
	ApiKey  string `json:"ApiKey"`  // API Key of the bot (ask @BotFather)
	ProjId  string `json:"ProjId"`  // Project ID on GAE
}

// Starts the telegram bot. The parameter conf represent the running
// configuration. The conf is mandatory otherwise the bot can't authenticate.
// The parameter resp is the Responder function. Also this parameter is
// mandatory, otherwise the bot don't know how to anser to user questions.
// TODO
func (t Bot) Startgae(conf Configuration, resp Responder) {
	// Settings management
	if len(conf.BotName) == 0 {
		fmt.Println("FATAL: Bot Name not set. Please check your configuration")
		os.Exit(1)
	}
	if len(conf.ApiKey) == 0 {
		fmt.Println("FATAL: API Key not set. Please check your configuration")
		os.Exit(1)
	}
	if len(conf.ProjId) != 0 {
		fmt.Println("INFO: Don't forget to visit:")
		fmt.Println(t.telegramWebhookURL(conf))
		fmt.Println("INFO: Just once, you won't be able to receive messages")
	}

	fmt.Println("INFO: Settings loaded!")
	fmt.Println("INFO: Working as: " + conf.BotName)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.postHandler(w, r, conf, resp)
	})
}

// TODO
func (t Bot) postHandler(w http.ResponseWriter, r *http.Request, conf Configuration, resp Responder) {

	body, err := ioutil.ReadAll(r.Body)
	ctx := appengine.NewContext(r)
	if err != nil {
		log.Errorf(ctx, "WARN: Malformed body from Telegram!", err)
		return
	}

	var message teleResults
	if err = json.Unmarshal(body, &message); err != nil {
		log.Errorf(ctx, "WARN: Telegram JSON Error: ", err)
	} else {
		log.Infof(ctx, "INFO: ##### Received message")
		SEND_URL := t.telegramSendURL(conf)
		client := urlfetch.Client(ctx)

		log.Infof(ctx, "INFO: Message: '"+message.Message.Text+"' From: '"+message.Message.Chat.Uname+"'")
		// Answer message
		var err error
		answer := t.getResponse(message.Message.Text, conf, resp, r)
		log.Infof(ctx, "INFO: Response: '"+answer+"'")

		vals := url.Values{
			"chat_id": {strconv.FormatInt(message.Message.Chat.Chatid, 10)},
			"text":    {answer}}
		if _, err = client.PostForm(SEND_URL, vals); err != nil {
			log.Errorf(ctx, "WARN: Could not send post request: %v\n", err)
		} else {
			log.Infof(ctx, "INFO: Answer: '"+answer+"' To: '"+message.Message.From.Uname+"'")
		}
	}
	fmt.Fprint(w, "Telebot working :)")
}

// Load a configuration from a Json file and returns a configuration.
// See file `settings.json.sample` to see how settings should be formatted.
func (t Bot) LoadSettings(filename string) (Configuration, error) {
	configuration := Configuration{}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("FATAL: Unable to find file "+filename, err)
		return configuration, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("FATAL: Unable to read file "+filename+"! Please copy from settings.json.sample", err)
		return configuration, err
	}
	return configuration, nil
}

// Returns the telegram send URL, used to send messages.
// The URL is built using the loaded configuration.
func (t Bot) telegramSendURL(conf Configuration) string {
	sendurl := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   "bot" + conf.ApiKey + "/sendMessage"}
	return sendurl.String()
}

// Returns the telegram webhook URL, used to receive messages.
// The URL is built using the loaded configuration.
func (t Bot) telegramWebhookURL(conf Configuration) string {
	gaeurl := url.URL{
		Scheme: "https",
		Host:   conf.ProjId + ".appspot.com"}
	teleurl := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org"}
	teleurl.Path = "bot" + conf.ApiKey + "/setWebhook"
	teleurl.RawQuery = "url=" + gaeurl.String()
	return teleurl.String()
}

// Process a single user message and returns the answer.
// This method will remove the @BotName (e.g. /start@TestBot) from received message
// to allow a unique interpretation of messages
func (t Bot) getResponse(message string, conf Configuration,
	resp Responder, request *http.Request) string {

	var answer string
	var err error
	message = strings.Replace(message, "@"+conf.BotName, "", 1)

	answer, err = resp(message, request)
	if err != nil {
		answer = "I'm not able to answer :("
	}
	return answer
}
