package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

type Run struct {
	Id       int     `json:"id"`
	Distance float64 `json:"distance"`
	Result   int64   `json:"result"`
	Date     int64   `json:"date"`
}

//implementing sort interface ByDate
type ByDate []Run

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date < a[j].Date }

var db gorm.DB
var dberror error
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func getRuns() ([]Run, error) {
	runs := []Run{}
	db.Find(&runs)

	//todo move sorting to own method
	sort.Sort(ByDate(runs))
	return runs, nil
}

//save will create run if not existing
func (run *Run) save() error {
	return db.Save(run).Error
}

func (run *Run) isValid() bool {
	return run.Distance > 0 && run.Result > 0 && run.Distance > 0
}

//todo move all db stuff to own file(/package)
func setupRunDb() error {
	//todo replcae db location with config value
	db, dberror = gorm.Open("sqlite3", "data/run_db.sqlite")
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
			Date:     time.Now().Unix(),
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

	Logger.Println("Method: ", r.Method)
	switch r.Method {
	case "GET":
		runs, err := getRuns()
		if err != nil {
			respond(w, 500, err)
		} else {
			respond(w, 200, runs)
		}
	case "POST":
		run := Run{}
		//todo run validation
		err = json.NewDecoder(r.Body).Decode(&run)
		if err != nil {
			Logger.Println("error: ", err)
			respond(w, 500, err)
		} else if run.isValid() == false {
			respond(w, 422, "invalid values")
		} else {
			err = run.save()
			if err != nil {
				respond(w, 500, err)
			} else {
				respond(w, 200, run)
			}
		}
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
	w.Header().Add("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	//w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	w.WriteHeader(statuscode)
	fmt.Fprintf(w, "%s", message)
}

//https://github.com/jordan-wright/gophish/blob/master/controllers/api.go
func main() {
	http.HandleFunc("/runs/", runHandler)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		Logger.Println(err)
	}
}
