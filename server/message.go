package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"olx/pkg"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	bot *tgbotapi.BotAPI
}

func New(apiKey string) *Client {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
	}

	return &Client{
		bot: bot,
	}
}

func (c *Client) SendMessage(text string, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"
	_, err := c.bot.Send(msg)
	return err
}

func OlxParser(s string, chatId int64) {
	var url string
	switch s {
	case "Мониторы":
		url = "https://www.olx.kz/d/elektronika/kompyutery-i-komplektuyuschie/monitory/astana/?search%5Border%5D=created_at:desc"
	case "Моноблоки":
		url = "https://www.olx.kz/d/elektronika/kompyutery-i-komplektuyuschie/nastolnye-kompyutery/astana/q-%D0%9C%D0%BE%D0%BD%D0%BE%D0%B1%D0%BB%D0%BE%D0%BA/?search%5Border%5D=created_at:desc"

	case "Системники":
		url = "https://www.olx.kz/d/elektronika/kompyutery-i-komplektuyuschie/nastolnye-kompyutery/astana/?search%5Border%5D=created_at:desc"
	}

	response, err := http.Get(url)
	Check(err)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	Check(err)

	massUrl := make([]string, 0)
	doc.Find("div.css-oukcj3").Find("div.css-1sw7q4x").Each(func(index int, item *goquery.Selection) {

		url, _ := item.Find("a").Attr("href")

		massUrl = append(massUrl, "https://www.olx.kz"+url)

	})
	// DataBase(title, "https://www.olx.kz"+url)

	for i := 0; i < len(massUrl); i++ {
		DataBase(massUrl[i], chatId)
		// fmt.Println(massUrl[i])
	}

	fmt.Println("Выполнено")

}
func Check(err error) {

	if err != nil {
		fmt.Println(err)
	}
}

func DataBase(s string, chatId int64) {
	c := New(pkg.BOT_TOKEN)

	database, err := sql.Open("sqlite3", "./listUrl"+strconv.Itoa(int(chatId))+".db")
	Check(err)
	defer database.Close()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS olx (id INTEGER PRIMARY KEY, url TEXT)")
	Check(err)
	statement.Exec()
	rows, _ := database.Query("SELECT id, url FROM olx")
	var id int

	var url1 string
	chek := false
	for rows.Next() {
		rows.Scan(&id, &url1)
		if s == url1 {
			// fmt.Println(url1)
			chek = true
		}

		// fmt.Printf("%d: %s \n", id, url1)

	}

	rows.Close()
	statement, err = database.Prepare("INSERT INTO olx ( url) VALUES ( ?)")
	Check(err)
	coin := 0

	if !chek {
		coin++
		fmt.Println("Новое объявление ", coin)
		c.SendMessage(s, int64(chatId))
		// c.SendMessage(s, int64(452639799)) //995356946
		// c.SendMessage(s, int64(995356946))
		statement.Exec(s)
	}
	// fmt.Println("Есть контакт")

}
