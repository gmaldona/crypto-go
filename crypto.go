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

type jsonData struct {
	Data      []Currency_t `json:"data"`
	Timestamp int64        `json:"timestamp"`
}

/**
@param:
	- msg: error message to log to loggly
*/
func throwLogError(msg string) {
	client := loggly.New("CryptoApi")
	logErr := client.EchoSend("error", msg)
	if logErr != nil {
		os.Exit(1)
	}
}

/**
@param:
	- body: string json response from the endpoint
@return:
	- []Currency_t: an array of currency_t structs representing the json response info

*/
func extractJsonData(body string) []Currency_t {
	var jsonData jsonData
	err := json.Unmarshal([]byte(body), &jsonData)

	if err != nil {
		throwLogError("Could not parse String into Json format.")
		log.Fatalln(err)
	}

	return jsonData.Data
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		fmt.Println("-----==== Starting HTTP worker ====-----")

		// GET request on Endpoint
		resp, err := http.Get("https://api.coincap.io/v2/assets")

		if err != nil {
			throwLogError("Could not pull data from API.")
			log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			throwLogError("Could not read data from the API.")
			log.Fatalln(err)
		}

		bodyStr := string(body)

		jsonData := extractJsonData(bodyStr)

		for i := 0; i < len(jsonData); i++ {
			fmt.Printf("Id:\t\t\t%s\n", jsonData[i].Id)
			fmt.Printf("Rank:\t\t\t%s\n", jsonData[i].Rank)
			fmt.Printf("Symbol:\t\t\t%s\n", jsonData[i].Symbol)
			fmt.Printf("Name:\t\t\t%s\n", jsonData[i].Name)
			fmt.Printf("Supply:\t\t\t%s\n", jsonData[i].Supply)
			fmt.Printf("Max Supply:\t\t%s\n", jsonData[i].MaxSupply)
			fmt.Printf("Market Cap (USD):\t%s\n", jsonData[i].MarketCapUsd)
			fmt.Printf("Volume 24 Hours (USD):\t%s\n", jsonData[i].VolumeUsd24Hr)
			fmt.Printf("Price (USD):\t\t%s\n", jsonData[i].PriceUsd)
			fmt.Printf("Change Percent 24 Hr: \t%s\n", jsonData[i].ChangePercent24Hr)
			fmt.Printf("VWAP 24 Hours:\t\t%s\n\n", jsonData[i].Vwap24Hr)
		}

		// Send a Success messgae to Loggly
		client := loggly.New("CryptoApi")
		logErr := client.EchoSend("info", "Data polled. "+strconv.Itoa(len(jsonData)))
		if logErr != nil {
			println(logErr)
			return
		}

		time.Sleep(1 * time.Minute)
	}

}
