package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	port       = os.Getenv("PORT")
	bitlyURL   = os.Getenv("BITLY_API_ENDPOINT")
	bitlyToken = os.Getenv("BITLY_TOKEN")
	botToken   = os.Getenv("TG_BOT_TOKEN")
)

const (
	urlRegexp = `^https?:\/\/[\S]*`
)

type shortURLObject struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Link      string `json:"link"`
	LongURL   string `json:"long_url"`
}

func getURLs(bot *tgbotapi.BotAPI, newEvents tgbotapi.UpdatesChannel) error {

	urlPattern, _ := regexp.Compile(urlRegexp)
	for update := range newEvents {

		// check regexp
		if urlPattern.MatchString(update.Message.Text) {

			longURL := update.Message.Text
			fmt.Println(longURL)

			shortURL, err := shortenURL(longURL)
			if err != nil {
				return err
			}
			fmt.Println(shortURL.ID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, shortURL.ID)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
	return nil
}

func shortenURL(longURL string) (shortURLObject, error) {

	values := map[string]string{"long_url": longURL}
	jsonPayload, _ := json.Marshal(values)

	client := &http.Client{}
	req, err := http.NewRequest("POST", bitlyURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return shortURLObject{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+bitlyToken)
	res, err := client.Do(req)
	if err != nil {
		return shortURLObject{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return shortURLObject{}, err
	}

	var shortURL shortURLObject
	err = json.Unmarshal(body, &shortURL)
	if err != nil {
		return shortURLObject{}, err
	}

	fmt.Println(shortURL.Link)

	return shortURL, nil
}

func main() {

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/updates")
	portToListen := ":" + port
	go http.ListenAndServe(portToListen, nil)

	err = getURLs(bot, updates)
	if err != nil {
		log.Fatal(err)
	}
}
