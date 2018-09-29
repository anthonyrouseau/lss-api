package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("testdb"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := Routes()
	fmt.Println("server starting up")
	http.ListenAndServe(":1337", r)
}
