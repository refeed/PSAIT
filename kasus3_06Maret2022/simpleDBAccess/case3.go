package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/studentProxy", studentProxyHandler)
	http.ListenAndServe(":8081", nil)
}

const (
	REMOTE_ADDRESS = "http://127.0.0.1:8080"
)

func studentProxyHandler(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	queryType := urlQuery.Get("type")
	queryId, _ := strconv.Atoi(urlQuery.Get("id"))
	queryNewName := ""

	switch queryType {
	case "update":
		queryNewName = urlQuery.Get("newName")
		jsonReq, _ := json.Marshal(map[string]interface{}{"id": queryId, "newName": queryNewName})
		req, _ := http.NewRequest(http.MethodPatch, REMOTE_ADDRESS+"/student", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		_, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
	case "post":
		queryNewName = urlQuery.Get("newName")
		jsonReq, _ := json.Marshal(map[string]string{"name": queryNewName})
		http.Post(REMOTE_ADDRESS+"/student", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	case "delete":
		jsonReq, _ := json.Marshal(map[string]int{"id": queryId})
		req, _ := http.NewRequest(http.MethodDelete, REMOTE_ADDRESS+"/student", bytes.NewBuffer(jsonReq))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		_, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
	}

	http.Redirect(w, r, "/", 300)
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

	renderHtml(w, localVersion, remoteVersion, getRemoteStudents())
}

func getUbuntuMariaDBVersion(ret *string) {
	resp, err := http.Get(REMOTE_ADDRESS + "/getVersion")
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

type RemoteStudent struct {
	Id   uint
	Name string
}

func getRemoteStudents() []RemoteStudent {
	var remoteStudents []RemoteStudent
	resp, err := http.Get(REMOTE_ADDRESS + "/student")
	if err != nil {
		log.Fatalln(err)
		return remoteStudents
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var studentsMap map[uint]string
	json.Unmarshal(body, &studentsMap)

	for id, name := range studentsMap {
		remoteStudents = append(remoteStudents, RemoteStudent{id, name})
	}

	return remoteStudents
}

func renderHtml(w io.Writer, hostSqlVer string, remoteSqlVer string, students []RemoteStudent) {
	fmt.Fprintf(w, `<html>
    <head>
        <title>
            Show database in native Golang
        </title>
    </head>
    <body>
        <h1>Simple DB Access Golang Native</h1>
        <p>Host db version: %v</p>
        <p>Remote db version: %v</p>
        <form action="/studentProxy" method="get">
			<input type="hidden" name="type" value="post" />
            <input type="text" name="newName" placeholder="Add student name">
            <input type="submit" value="Add">
        </form>
        <table>
            <tr>
                <th>
                    ID
                </th>
                <th>
                    Name
                </th>
                <th>
                    Action
                </th>
            </tr>`, hostSqlVer, remoteSqlVer)

	for _, student := range students {
		fmt.Fprintf(w, `<tr>
                <td>%[1]v</td>
                <td>%[2]v</td>
                <td><form action="/studentProxy" method="get">
                        <input type="hidden" name="id" value="%[1]v" />
                        <input type="hidden" name="type" value="update" />
                        <input type="text" name="newName" placeholder="New name">
                        <input type="submit" value="Update">
                    </form>|<a href="/studentProxy?type=delete&id=%[1]v">Delete</a></td>
            </tr>`, student.Id, student.Name)
	}
	fmt.Fprintf(w, `</table>
					</body>
				</html>`)
}
