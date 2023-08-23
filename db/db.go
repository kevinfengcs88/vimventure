package main

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	connectionString := fmt.Sprintf("user=postgres password=%s dbname=person sslmode=disable", DB_PASSWORD)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("first error encountered")
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM person")
	if err != nil {
		fmt.Println("second error encountered")
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Printf("Type of rows is %T\n", rows)
	for rows.Next() {
		var id int
		var name string
		var country string
		err := rows.Scan(&id, &name, &country)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Country: %s\n", id, name, country)
	}
}
