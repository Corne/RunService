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

func test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", "henk")
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	p := Person{
		Name: "Henk",
		Age:  50,
	}

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Print("error: ", err)
	}
	fmt.Fprintf(w, "%s", b)
}

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
	//Logger.Println("db: ", &db == nil)
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
		Logger.Println("creating table, ")
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
	test := []Run{}
	db.Find(&test)
	Logger.Println("Hello, ", test)

	value := new(bool)
	db.Table("runs").HasTable(&value)
	Logger.Println("Hello, ", *value)
	return nil
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	err := setupRunDb()
	if err != nil {
		respond(w, 500, err.Error())
	}

	if r.Method == "GET" {
		runs, err := getRuns()
		if err != nil {
			respond(w, 500, err.Error())
		}

		b, err := json.Marshal(runs)
		if err != nil {
			fmt.Print("error: ", err)
		}

		respond(w, 200, b)
	}
}

func respond(w http.ResponseWriter, statuscode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	fmt.Fprintf(w, "%s", body)
}

func main() {

	http.HandleFunc("/test/", test)
	http.HandleFunc("/hello/", hello)

	http.HandleFunc("/runs/", runHandler)

	http.ListenAndServe(":8080", nil)
}
