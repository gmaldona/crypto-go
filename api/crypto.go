package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jamespearly/loggly"

	crypto "cryptocurrency/structs"
)

func tableExists(tableName string, tableNames []*string) bool {
	for _, name := range tableNames {
		if tableName == *name {
			return true
		}
	}
	return false
}

func createTable(db *dynamodb.DynamoDB) {
	_, err := db.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("Maldonado-CryptoBro"),
	})

	if err != nil {
		log.Println(err)
		throwLogError("Could not create table in DynamoDB.")
	}
}

func PutItem(currency crypto.Currency_t, tableName string, db *dynamodb.DynamoDB) {
	_, err := db.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(currency.Id),
			},
			"Rank": {
				S: aws.String(currency.Rank),
			},
			"Symbol": {
				S: aws.String(currency.Symbol),
			},
			"Name": {
				S: aws.String(currency.Name),
			},
			"Supply": {
				S: aws.String(currency.Supply),
			},
			"MaxSupply": {
				S: aws.String(currency.MaxSupply),
			},
			"MarketCapUsd": {
				S: aws.String(currency.MarketCapUsd),
			},
			"VolumeUsd24Hr": {
				S: aws.String(currency.VolumeUsd24Hr),
			},
			"PriceUsd": {
				S: aws.String(currency.PriceUsd),
			},
			"ChangePercent24Hr": {
				S: aws.String(currency.ChangePercent24Hr),
			},
			"Vwap24Hr": {
				S: aws.String(currency.Vwap24Hr),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Println(err)
		throwLogError("Could not make a DynamoDB table entry.")
	}
}

type jsonData struct {
	Data      []crypto.Currency_t `json:"data"`
	Timestamp int64               `json:"timestamp"`
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
func extractJsonData(body string) []crypto.Currency_t {
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

		// Send a Success message to Loggly
		client := loggly.New("CryptoApi")
		logErr := client.EchoSend("info", "Data polled. "+strconv.Itoa(len(jsonData)))
		if logErr != nil {
			println(logErr)
			return
		}

		// Open a new DynamoDB session
		db := dynamodb.New(session.Must(session.NewSession(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String("https://dynamodb.us-east-1.amazonaws.com"),
		})))

		tables, _ := db.ListTables(&dynamodb.ListTablesInput{})
		tables.String()

		if !tableExists(crypto.DB_TABLE_NAME, tables.TableNames) {
			createTable(db)
		}

		for _, currency := range jsonData {
			PutItem(currency, crypto.DB_TABLE_NAME, db)
		}

		// Send a Success message to Loggly
		client = loggly.New("CryptoApi")
		logErr = client.EchoSend("info", "Entered all currencies into database.")
		if logErr != nil {
			println(logErr)
			return
		}

		time.Sleep(60 * time.Hour)
	}

}
