package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Channel struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Url     string `json:"url"`
	Title   string `json:"title"`
	Created int64  `json:"created"`
	Count   int    `json:"count"` //вообще там всегда вроде 0, но добавил для совместимости
}

type Author struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Www      string `json:"www"`
	Provider string `json:"provider"`
	OpenID   string `json:"openId"`
	Verify   bool   `json:"verify"`
	Notify   bool   `json:"notify"`
}

type Comment struct {
	ID       int     `json:"id"`
	SiteID   int     `json:"siteId"`
	ParentID int     `json:"parentId"`
	Path     []int   `json:"path"`
	Message  string  `json:"message"`
	Rating   int     `json:"rating"`
	Status   string  `json:"status"`
	Created  int64   `json:"created"`
	Author   Author  `json:"author"`
	Chan     Channel `json:"chan"`
	IP       string  `json:"ip"`
	Modified string  `json:"modified"`
}

func main() {
	log.Println("Start import comments from cackle.me...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("Get channels from widget ID: %s", os.Getenv("ID"))
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil || timeout <= 0 {
		timeout = 5
	}
	// читаем каналы (по сути это список страниц, к которым есть комментарии)
	apiUrlChannels := fmt.Sprintf("http://cackle.me/api/3.0/comment/chan/list.json?id=%s&siteApiKey=%s&accountApiKey=%s",
		os.Getenv("ID"), os.Getenv("SITE_API_KEY"), os.Getenv("ACCOUNT_API_KEY"))
	//log.Println(apiUrlChannels)
	page := 0
	var allChannels []Channel
	for { //цикл по получению списка каналов, следуя рекомендации сервиса запросы делаются с паузами и по 100 шт
		count, channels, err := getChannels(apiUrlChannels, page)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Get channels, iteration %d, count %d\n", page+1, len(channels))
		page++
		allChannels = append(allChannels, channels...)
		if count < 100 {
			break
		}
		time.Sleep(time.Duration(timeout) * time.Second)
	}
	log.Printf("All channels count %d\n", len(allChannels))
	log.Printf("Saving all channels to file...")

	file, err := json.MarshalIndent(allChannels, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("channels.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Done\n")
	//читаем комментарии
	log.Println("Getting comments...")
	apiUrlComments := fmt.Sprintf("http://cackle.me/api/3.0/comment/list.json?id=%s&siteApiKey=%s&accountApiKey=%s",
		os.Getenv("ID"), os.Getenv("SITE_API_KEY"), os.Getenv("ACCOUNT_API_KEY"))
	var allComments []Comment
	for index, channel := range allChannels {
		log.Printf("Get comments from channel id %d (%d from %d),\n URL: %s\n Title: %s\n", channel.Id, index, len(allChannels), channel.Url, channel.Title)
		commentId := 0
		currentCount := 0
		for {
			time.Sleep(time.Duration(timeout) * time.Second)
			count, comments, err := getComments(apiUrlComments, channel.Channel, commentId)
			log.Println("Get block with count: ", count)
			currentCount += count
			if err != nil {
				log.Println(err)
				break
			} else {
				allComments = append(allComments, comments...)
				if count < 100 {
					commentId = 0
				} else {
					commentId = comments[99].ID
				}
			}
			if commentId == 0 {
				break
			}
		}
		log.Printf("Count comments: %d", currentCount)
		currentCount = 0
	}
	log.Println("Done...")
	log.Printf("All comments count %d\n", len(allComments))
	log.Printf("Saving all comments to file...")
	file, err = json.MarshalIndent(allComments, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("comments.json", file, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done...")
}

func getChannels(baseApiUrl string, page int) (count int, channels []Channel, err error) {
	if page > 0 {
		baseApiUrl = baseApiUrl + fmt.Sprintf("&page=%d", page)
	}
	response, err := http.Get(baseApiUrl)
	if err != nil {
		return 0, nil, err
	}
	var result map[string][]Channel
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return 0, nil, err
	}
	count = len(result["chans"])
	channels = result["chans"]
	err = nil
	return
}

func getComments(baseApiUrl string, channel string, commentId int) (count int, comments []Comment, err error) {
	baseApiUrl = baseApiUrl + fmt.Sprintf("&chan=%s", channel) + fmt.Sprintf("&commentId=%s", commentId)

	response, err := http.Get(baseApiUrl)
	if err != nil {
		return 0, nil, err
	}
	var result map[string][]Comment
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return 0, nil, err
	}
	count = len(result["comments"])
	comments = result["comments"]
	err = nil
	return
}
