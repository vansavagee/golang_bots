package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"talk_bot/internal/constants"
	"talk_bot/internal/gpt"
	sqlqueries "talk_bot/internal/sqlQueries"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TgBotToken"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updConfig := tgbotapi.NewUpdate(0)
	updConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updConfig)
	if err != nil {
		log.Fatalf("enable to connect with bot: %v", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		} else if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название устройства.")
			bot.Send(msg)
		} else {
			log.Printf("[%s] %s", update.Message.Chat.UserName, update.Message.Text)
			sqlqueries.DeleteData()
			deviceName := update.Message.Text
			allDevices, err := sqlqueries.SelectAll()
			log.Printf("ВСЕ ДЕВАЙСЫ: %v", allDevices)
			if err != nil {
				log.Printf("error from query(SelectAll): %v", err)
			}
			s, err := gpt.GetResponseFromGPT(constants.QueryGptPart_1 + deviceName + constants.QueryGptPart_2 + allDevices)
			if err != nil {
				log.Printf("error from gpt(GetResponseFromGPT): %v", err)
			}
			ids := convertStringsToInts(s)
			log.Printf("ids: %v", ids)
			if len(ids) == 0 || ids[0] == -1 {
				s, err = gpt.GetResponseFromGPT(constants.QueryGptPart_1 + deviceName + constants.QueryGptPart_2 + allDevices)
				if err != nil {
					log.Printf("error from gpt(GetResponseFromGPT): %v", err)
				}
				ids = convertStringsToInts(s)
			}
			if len(ids) == 0 || ids[0] == -1 {
				log.Printf("NOT FOUND %v", ids)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Устройство \"%v\" не найдено.", deviceName))
				bot.Send(msg)
			} else {
				s, err := sqlqueries.SelectAllByIdsArray(ids)
				if err != nil {
					log.Printf("error from query(SelectAllByIdsArray): %v", err)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Устройство \"%v\" найдено:\n%v", deviceName, s))
				bot.Send(msg)
			}
		}
	}
}
func convertStringsToInts(s string) []int {
	buffer := strings.Split(s, ",")
	res := make([]int, 0, len(buffer))
	for i := 0; i < len(buffer); i++ {
		v, err := strconv.Atoi(buffer[i])
		if err == nil {
			res = append(res, v)

		}
	}
	return res
}
