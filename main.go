package main

import (
	"database/sql"
	"example/web-service-gin/pkg/database"
	"example/web-service-gin/pkg/routes"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var err error

	// Get a database handle
	// "root:password1@tcp(127.0.0.1:3306)/test"
	// db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/recordings")
	database.DB, err = sql.Open("mysql", database.DbURL(database.BuildDBConfig()))

	if err != nil {
		log.Fatal(err)
	}
	// See "Important settings" section.
	database.DB.SetConnMaxLifetime(time.Minute * 3)
	database.DB.SetMaxOpenConns(10)
	database.DB.SetMaxIdleConns(10)

	defer database.DB.Close()

	// pingErr := db.Ping()
	// if pingErr != nil {
	// 	log.Fatal(pingErr)
	// }

	fmt.Println("Connected!")

	router := routes.Routes()

	router.Run("localhost:8080")
}
