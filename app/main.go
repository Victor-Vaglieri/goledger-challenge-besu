package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"desafio/blockchain"
	"desafio/database"
	"desafio/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("FALHA - .env")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db := database.InitDB(connStr)
	defer db.Close()

	client, contract := blockchain.InitContract(os.Getenv("BESU_NODE_URL"), os.Getenv("CONTRACT_ADDRESS"))
	defer client.Close()

	api := &handlers.API{
		DB:       db,
		Contract: contract,
		Client:   client,
	}

	http.HandleFunc("/get", api.GetHandler)
	http.HandleFunc("/set", api.SetHandler)
	http.HandleFunc("/sync", api.SyncHandler)
	http.HandleFunc("/check", api.CheckHandler)

	fmt.Println("HTTP WORKING - 8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
