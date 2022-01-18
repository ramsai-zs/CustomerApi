package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

var db *sql.DB //db variable

//fetching data from database to the struct
type Customer struct {
	ID      int
	Name    string
	Phone   int
	Address string
}

func main() {
	fmt.Println("Go MySQL Tutorial")
	var err error
	db, err = sql.Open("mysql", "sai:password@tcp(127.0.0.1:3306)/Customers")
	if err != nil {
		fmt.Println("error:", err)
	}
	defer db.Close()

	http.HandleFunc("/customer", handler)
	log.Fatal(http.ListenAndServe(":8084", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		get(w, r)
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		post(w, r)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// post reads the JSON body and inserts in the database
func post(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &customer)
	_, err := db.Exec("INSERT INTO customers (ID,Name, PhoneNo,Address) VALUES(?,?,?,?)", customer.ID, customer.Name, customer.Phone, customer.Address)
	if err != nil {
		log.Println("error:", err)

		w.WriteHeader(500)
		_, _ = w.Write([]byte("something unexpected happened"))
		return
	}
	_, _ = w.Write([]byte("success"))
}

// get retrieves the data from database and writes data as a JSON.
func get(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * from customers;")
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(500)
	}
	defer rows.Close()
	var customers []Customer
	for rows.Next() {
		var a Customer
		_ = rows.Scan(&a.ID, &a.Name, &a.Phone, &a.Address)
		customers = append(customers, a)
	}

	resp, _ := json.Marshal(customers)
	_, _ = w.Write(resp)
}
