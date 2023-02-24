package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect() {
	_, err := sql.Open("postgres", "host=db port=5432 user=pickpin password=pickpinpswd dbname=pickpindb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("okay")
}
