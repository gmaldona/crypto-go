package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jamespearly/loggly"
)

type Currency_t struct {
	Id                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
	Explorer          string `json:"explorer"`
}

type Json_t struct {
	Data      []Currency_t `json:"data"`
	Timestamp int64        `json:"timestamp"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		fmt.Println(loggly.New("CryptoApi"))
		fmt.Printf("Token: %s\n", os.Getenv("LOGGLY_TOKEN"))

		fmt.Println("-----==== Starting HTTP worker ====-----")

		resp, err := http.Get("https://api.coincap.io/v2/assets")
		if err != nil {
			client := loggly.New("CryptoApi")
			logErr := client.EchoSend("error", "Could not pull data from API.")
			if logErr != nil {
				return
			}
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			client := loggly.New("CryptoApi")
			logErr := client.EchoSend("error", "Could not read data from the API.")
			if logErr != nil {
				return
			}
			log.Fatalln(err)
		}

		bodyStr := string(body)

		var json_t Json_t
		err = json.Unmarshal([]byte(bodyStr), &json_t)

		if err != nil {
			client := loggly.New("CryptoApi")
			logErr := client.EchoSend("error", "Could not parse String into Json format.")
			if logErr != nil {
				return
			}
			log.Fatalln(err)
		}

		for i := 0; i < len(json_t.Data); i++ {
			fmt.Printf("Id : %s\n", json_t.Data[i].Id)
			fmt.Printf("Rank : %s\n", json_t.Data[i].Rank)
			fmt.Printf("Symbol : %s\n", json_t.Data[i].Symbol)
			fmt.Printf("Name : %s\n", json_t.Data[i].Name)
			fmt.Printf("Supply : %s\n", json_t.Data[i].Supply)
			fmt.Printf("Max Supply : %s\n", json_t.Data[i].MaxSupply)
			fmt.Printf("Market Cap (USD) : %s\n", json_t.Data[i].MarketCapUsd)
			fmt.Printf("Volume 24 Hours (USD) : %s\n", json_t.Data[i].VolumeUsd24Hr)
			fmt.Printf("Price (USD) : %s\n", json_t.Data[i].PriceUsd)
			fmt.Printf("Change Percent 24 Hours : %s\n", json_t.Data[i].ChangePercent24Hr)
			fmt.Printf("VWAP 24 Hours : %s\n\n", json_t.Data[i].Vwap24Hr)
		}

		client := loggly.New("CryptoApi")
		logErr := client.EchoSend("info", "Data polled. "+strconv.Itoa(len(json_t.Data)))
		if logErr != nil {
			println(logErr)
			return
		}

		fmt.Println("Loggly Message Success.")
		time.Sleep(1 * time.Minute)
	}

}
