package server

import (
	"log"
	"olx/pkg"

	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"), tgbotapi.NewInlineKeyboardButtonData("2", "2"), tgbotapi.NewInlineKeyboardButtonData("3", "3")), tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("4", "4"), tgbotapi.NewInlineKeyboardButtonData("5", "5"), tgbotapi.NewInlineKeyboardButtonData("6", "6")))

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мониторы"),
		tgbotapi.NewKeyboardButton("Моноблоки"),
		tgbotapi.NewKeyboardButton("Системники"),
		tgbotapi.NewKeyboardButton("Оптом электороника"),
	),
)

func Run() {

	bot, err := tgbotapi.NewBotAPI(pkg.BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var res string
		if update.Message != nil { // If we got a message
			switch update.Message.Text {
			case "/start":
				res = `Добро пожаловать на наш сервис уведомления объявлений.
				В появившемся меню вы можете выбрать тип объявлений который вы хотите мониторить
				`
			case "Мониторы":
				res = "БОТ начал мониторинг"
				go OlxCheck("Мониторы", update.Message.Chat.ID)
			case "Моноблоки":
				res = "БОТ начал мониторинг"
				go OlxCheck("Моноблоки", update.Message.Chat.ID)
			case "Системники":
				res = "БОТ начал мониторинг"
				go OlxCheck("Системники", update.Message.Chat.ID)
			case "Оптом электороника":
				res = "БОТ начал мониторинг"
				go OlxCheck("Оптом электороника", update.Message.Chat.ID)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
			if update.Message.Text == "/start" {
				msg.ReplyMarkup = mainMenu

			}

			bot.Send(msg)
		}
	}

}

func OlxCheck(s string, chatId int64) {

	for {
		OlxParser(s, chatId)
		time.Sleep(15 * time.Minute)
	}
}

// if _, err = bot.Send(msg); err != nil {
// 	panic(err)
// }
// } else if update.CallbackQuery != nil {
// // Respond to the callback query, telling Telegram to show the user
// // a message with the data received.
// callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
// if _, err := bot.Request(callback); err != nil {
// 	panic(err)
// }
// var res string
// // And finally, send a message containing the data received.
// if update.CallbackQuery.Data == "2" {
// 	res = "yyyyy"
// } else {
// 	res = update.CallbackQuery.Data
// }

// msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, res)
// if _, err := bot.Send(msg); err != nil {
// 	panic(err)
// }
