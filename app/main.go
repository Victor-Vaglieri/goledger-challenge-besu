package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL
	connStr := "host=localhost port=5432 user=admin password=admin dbname=goledger sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("FALHA - banco: %v", err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS contract_state (
			id SERIAL PRIMARY KEY,
			contract_value VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("FALHA - tabela: %v", err)
	}

	// Besu
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("FALHA - Besu: %v", err)
	}
	defer client.Close()

	// Definição das rotas da API REST
	http.HandleFunc("GET /get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Endpoint GET.")
	})

	fmt.Println("Servidor HTTP: 8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
