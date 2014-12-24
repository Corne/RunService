package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

type Run struct {
	Id       int     `json:"id"`
	Distance float64 `json:"distance"`
	Result   int64   `json:"result"`
	Date     int64   `json:"date"`
}

var db gorm.DB
var dberror error
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func getRuns() ([]Run, error) {
	runs := []Run{}
	db.Find(&runs)
	return runs, nil
}

func setupRunDb() error {
	db, dberror = gorm.Open("sqlite3", "data/run_db.sqlite") //todo replcae db location with config value
	db.LogMode(false)
	db.SetLogger(Logger)

	if dberror != nil {
		Logger.Println(dberror)
		return dberror
	}

	_, err := os.Stat("data/run_db.sqlite")
	if err != nil { //db doesn't exists
		db.CreateTable(Run{})

		example_run := Run{
			Distance: 5.0,
			Result:   60 * 23,
			Date:     time.Now().UnixNano(),
		}

		err = db.Save(&example_run).Error
		if err != nil {
			Logger.Println(err)
		}
	}
	return nil
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	err := setupRunDb()
	if err != nil {
		respond(w, 500, err)
	}

	if r.Method == "GET" {
		runs, err := getRuns()
		if err != nil {
			respond(w, 500, err)
		}

		respond(w, 200, runs)
	}
}

/*
 *	Marshals response body to json, and writes away response
 */
func respond(w http.ResponseWriter, statuscode int, body interface{}) {
	message, err := json.Marshal(body)
	if err != nil {
		fmt.Print("error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	fmt.Fprintf(w, "%s", message)
}

func main() {
	http.HandleFunc("/runs/", runHandler)

	http.ListenAndServe(":8080", nil)
}
