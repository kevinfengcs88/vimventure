package db

import (
	"fmt"
	"log"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func dbConnect(user string, password string, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("first error encountered")
		log.Fatal(err)
	}
	return db
	// defer db.Close()
	//
	// rows, err := db.Query("SELECT * FROM person")
	// if err != nil {
	// 	fmt.Println("second error encountered")
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// fmt.Printf("Type of rows is %T\n", rows)
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	var country string
	// 	err := rows.Scan(&id, &name, &country)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("ID: %d, Name: %s, Country: %s\n", id, name, country)
	// }
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	db := dbConnect("postgres", DB_PASSWORD, "person")

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
