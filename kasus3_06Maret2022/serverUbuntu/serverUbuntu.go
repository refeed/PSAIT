package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/getVersion", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Access the local db
	db, err := sql.Open("mysql", "refeed:wikimedia@/")
	defer db.Close()

	if err != nil {
		panic(err)
	}

	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)

	fmt.Fprintf(w, version)
}
