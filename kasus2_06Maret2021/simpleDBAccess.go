package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "refeed:wikimedia@tcp(192.168.56.2)/")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)

	fmt.Fprint(w, "Version: %v", version)
}
