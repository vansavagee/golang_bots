package botqueries

import (
	"log"
	"talk_bot/internal/sqlqueries"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func ShowResultInfo(message *tgbotapi.Message, bot *tgbotapi.BotAPI, company, model string) {

	data, err := sqlqueries.SelectAll()
	if err != nil {
		log.Fatalf("query error: %v", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, data)

	bot.Send(msg)
}
func MakeRowsCompany(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите производителя:")

	data, err := sqlqueries.SelectUniqCompanies()
	if err != nil {
		log.Fatalf("query error: %v", err)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, v := range data {
		button := tgbotapi.NewInlineKeyboardButtonData(v, v)
		row := []tgbotapi.InlineKeyboardButton{button}
		rows = append(rows, row)
	}
	button := tgbotapi.NewInlineKeyboardButtonData("/start", "/start")
	row := []tgbotapi.InlineKeyboardButton{button}
	rows = append(rows, row)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg.ReplyMarkup = menu

	bot.Send(msg)

}
func MakeRowsDevices(message *tgbotapi.Message, bot *tgbotapi.BotAPI, company string) bool {

	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите устройство:")

	data, err := sqlqueries.SelectUniqDevices(company)
	if err != nil {
		log.Fatalf("query error2: %v", err)
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, v := range data {
		button := tgbotapi.NewInlineKeyboardButtonData(v, v)
		row := []tgbotapi.InlineKeyboardButton{button}
		rows = append(rows, row)
	}
	button := tgbotapi.NewInlineKeyboardButtonData("/start", "/start")
	row := []tgbotapi.InlineKeyboardButton{button}
	rows = append(rows, row)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg.ReplyMarkup = menu
	bot.Send(msg)
	return true
}
