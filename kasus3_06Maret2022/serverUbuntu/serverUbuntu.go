package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DBPATH = "refeeds:12345678@/student_kasus3"
)

func main() {
	http.HandleFunc("/getVersion", handler)
	http.HandleFunc("/student", studentHandler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	rows, _ := queryDB("SELECT VERSION()")
	rows.Next()
	var version string
	rows.Scan(&version)
	fmt.Fprintf(w, version)
	rows.Close()
}

func studentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, _ := queryDB("SELECT id, name FROM student")
		students := map[int]string{}
		for rows.Next() {
			var (
				id   int
				name string
			)
			if err := rows.Scan(&id, &name); err != nil {
				break
			}

			students[id] = name
		}
		jsonStr, _ := json.Marshal(students)
		fmt.Fprintf(w, string(jsonStr))
		return
	case "POST":
		type PostRequest struct {
			Name string
		}

		var postRequest PostRequest

		err := json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			http.Error(w, "JSON structure is not valid", 400)
			return
		}

		db, err := getDB()
		defer db.Close()
		_, err = db.Exec("INSERT INTO student (name) VALUES (?)", postRequest.Name)
		if err != nil {
			http.Error(w, "Insert data failed", 501)
			return
		}
		replyHttpSuccessJson(w)
	case "DELETE":
		type DeleteRequest struct {
			Id uint
		}

		var delRequest DeleteRequest

		err := json.NewDecoder(r.Body).Decode(&delRequest)
		if err != nil {
			http.Error(w, "JSON structure is not valid", 400)
			return
		}
		db, err := getDB()
		defer db.Close()

		_, err = db.Exec("DELETE FROM student WHERE id = ?", delRequest.Id)
		if err != nil {
			http.Error(w, "Delete failed", 500)
			fmt.Println("%v", err)
			return
		}
		replyHttpSuccessJson(w)
	case "PATCH":
		type PatchRequest struct {
			Id      uint
			NewName string
		}
		var patchRequest PatchRequest
		err := json.NewDecoder(r.Body).Decode(&patchRequest)
		if err != nil {
			http.Error(w, "JSON structure is not valid", 400)
			return
		}

		db, err := getDB()
		defer db.Close()
		_, err = db.Exec("UPDATE student SET name=? WHERE id=?",
			patchRequest.NewName, patchRequest.Id)
		if err != nil {
			http.Error(w, "Update data failed", 501)
			return
		}
		replyHttpSuccessJson(w)
	default:
		fmt.Fprintf(w, "Only GET, POST, PUT, and DELETE are supported")
	}
}

func queryDB(query string) (*sql.Rows, error) {
	db, err := getDB()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	rows, err := db.Query(query)

	return rows, err
}

func getDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", DBPATH)
	return db, err
}

func replyHttpSuccessJson(w http.ResponseWriter) {
	fmt.Fprintf(w, `{"status": "success"}`)
}
