package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/url"
	"os"
	"strconv"
	"strings"
)

type Posts []Post
type Post struct {
	Id int
	Content ContextInfo
	Type string
}
type ContextInfo struct {
	Rendered string
}

func main(){
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)


	var ucfg = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	chPosts := getPosts()
	chatID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"),10, 64)
		for _, msgImg := range chPosts {
			fmt.Println(msgImg)

			msg := tgbotapi.NewPhotoShare(chatID, msgImg)
			bot.Send(msg)

			fmt.Println(msg)
		}
}


func getPosts() []string{
	var out  []string
	var err error
	r, err := http.Get(os.Getenv("URL_API"))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	data := Posts{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	postHtml := strings.Join(strings.Split(data[0].Content.Rendered, "\n"),"")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(postHtml))
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		imagesUrl, _ := s.Find("img").Attr("src")
		if imagesUrl != "" {
			out = append(out, imagesUrl)
		}
	})
	return out
}


