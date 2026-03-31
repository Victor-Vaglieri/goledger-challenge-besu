package database

import (
        "database/sql"
        "log"

        _ "github.com/lib/pq"
)

func InitDB(connStr string) *sql.DB {
        db, err := sql.Open("postgres", connStr)
        if err != nil {
                log.Fatalf("FALHA - banco: %v", err)
        }

        _, err = db.Exec(`
                CREATE TABLE IF NOT EXISTS contract_state (
                        id SERIAL PRIMARY KEY,
                        contract_value VARCHAR(255) NOT NULL
                );
        `)
        if err != nil {
                log.Fatalf("FALHA - tabela: %v", err)
        }

        return db
}


