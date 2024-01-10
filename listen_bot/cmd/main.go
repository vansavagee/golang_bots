package main

import (
	"encoding/json"
	"fmt"
	"listen_bot/internal/gpt"
	"listen_bot/internal/model"
	"listen_bot/internal/sqlqueries"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// info, err := GetResponseFromGPT("how are you?")
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Print(info)
	// }
	godotenv.Load()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TgBotToken"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	customer := ""
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		} else if update.Message.IsCommand() {
			commandWithArgs := strings.SplitN(update.Message.Text, " ", 2)
			customer = commandWithArgs[1]

		} else {

			log.Printf("[%s] %s", update.Message.Chat.UserName, update.Message.Text)
			s, err := gpt.GetResponseFromGPT(os.Getenv("QueryPrefixGPT") + update.Message.Text)
			if err != nil {
				log.Printf("error from gpt: %v", err)
			}

			// Создаем экземпляр структуры Vehicle
			var devices model.Devices

			// Декодируем JSON-строку в структуру
			log.Println(s)

			err = json.Unmarshal([]byte(s), &devices)
			if err != nil {
				log.Printf("Ошибка декодирования JSON, попытка 2: %v", err)
				devices = model.Devices{}
				s, err = gpt.GetResponseFromGPT(os.Getenv("QueryPrefixGPT") + update.Message.Text)
				if err != nil {
					log.Printf("error from gpt: %v", err)
				}
				err = json.Unmarshal([]byte(s), &devices)
			}
			if err != nil {
				log.Printf("Ошибка декодирования JSON: %v", err)

			} else {
				s := ""
				for _, v := range devices.Devices {
					s += fmt.Sprintf("%v %v\n", v, customer)
					v.Customer = customer
					err := sqlqueries.InsertData(v)
					if err != nil {
						log.Printf("error from insert: %v", err)
					}
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Устройства добавлены:\n%v", s))
				bot.Send(msg)
			}

		}

	}
}
