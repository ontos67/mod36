// Сервер GoNews.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"Agrigator/pkg/api"
	"Agrigator/pkg/rss"
	storage "Agrigator/pkg/storage/pstg"
)

type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"cicle"`
}

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	chPosts := make(chan []storage.Article)
	chErrs := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	go func() {
		for posts := range chPosts {
			db.SaveArticles(posts)
		}
	}()

	go func() {
		for err := range chErrs {
			log.Println("ошибка:", err)
		}
	}()

	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// Асинхронное чтение потока RSS. Раскодированные
// новости и ошибки пишутся в каналы.
func parseURL(url string, db *storage.DB, posts chan<- []storage.Article, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
		fmt.Println("прошла минута")
	}
}
