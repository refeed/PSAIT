package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "refeeds:12345678@/")
	defer db.Close()

	if err != nil {
		panic(err)
	}

	var localVersion string
	var remoteVersion string
	db.QueryRow("SELECT VERSION()").Scan(&localVersion)
	getUbuntuMariaDBVersion(&remoteVersion)

	fmt.Fprintln(w, "Host SQL Version: ", localVersion)
	fmt.Fprintln(w, "Ubuntu Server SQL Version: ", remoteVersion)
}

func getUbuntuMariaDBVersion(ret *string) {
	resp, err := http.Get("http://192.168.56.2:8080/getVersion")
	if err != nil {
		log.Fatalln(err)
		*ret = ""
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		*ret = ""
		return
	}

	sb := string(body)
	*ret = sb
}
