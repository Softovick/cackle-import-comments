package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Channel struct {
	Id      uint64 `json:"id"`
	Channel string `json:"channel"`
	Url     string `json:"url"`
	Title   string `json:"title"`
	Created uint64 `json:"created"`
	Count   int    `json:"count"` //вообще там всегда вроде 0, но добавил для совместимости
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

	apiUrlChannels := fmt.Sprintf("http://cackle.me/api/3.0/comment/chan/list.json?id=%s&siteApiKey=%s&accountApiKey=%s",
		os.Getenv("ID"), os.Getenv("SITE_API_KEY"), os.Getenv("ACCOUNT_API_KEY"))
	fmt.Println(apiUrlChannels)
	page := 0
	var allChannels []Channel
	for {
		count, channels, err := getChannels(apiUrlChannels, page)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Get channels, iteration %d, count %d\n", page+1, len(channels))
		page++
		allChannels = append(allChannels, channels...)
		if count < 100 {
			break
		}
		time.Sleep(time.Duration(timeout) * time.Second)
	}
	fmt.Printf("All channels count %d\n", len(allChannels))

}

func getChannels(baseApiUrl string, page int) (count int, channels []Channel, err error) {
	if page > 0 {
		baseApiUrl = baseApiUrl + fmt.Sprintf("&page=%d", page)
	}
	response, err := http.Get(baseApiUrl)
	if err != nil {
		return 0, []Channel{}, err
	}
	var result map[string][]Channel
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return 0, []Channel{}, err
	}
	count = len(result["chans"])
	channels = result["chans"]
	err = nil
	return
}
