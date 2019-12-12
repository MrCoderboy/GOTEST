package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	host      = "localhost"
	port      = 5432
	user      = "postgres"
	password1 = "123"
	dbname    = "auzziodb"
)

func SuggestedAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "appliation/json")
	searchstring := vars["searchstring"]
	//w.Write([]byte(fmt.Sprintf("%s\n", searchstring)))
	fmt.Println(searchstring)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password1, dbname)
	db, _ := sql.Open("postgres", psqlInfo)
	defer db.Close()
	sqlStatement := `select jsonbdata from public."Address_Search" where document @@ to_tsquery($1);`
	rows, err := db.Query(sqlStatement, searchstring)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var full_address string
		err = rows.Scan(&full_address)
		if err != nil {
			panic(err)
		}
		//userJson, _ := json.Marshal(full_address)

		//fmt.Fprintf(w,userJson)
		w.Header().Set("Content-Type", "appliation/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, full_address+"\n")
		w.WriteHeader(http.StatusOK)
		//w.Write(userJson)

	}

}

func HandleAPI(w http.ResponseWriter, r *http.Request) {

	// Queries will automatically break down the &variables
	// you don't need to worry about the ampersand & in the
	// URL.

	vars := mux.Vars(r)
	searchstring := vars["searchstring"]
	w.Write([]byte(fmt.Sprintf("%s\n", searchstring)))
	fmt.Println(searchstring)

}

func main() {
	router := mux.NewRouter().StrictSlash(false)
	//handlerequests()
	router.HandleFunc("/address/search/{searchstring}", SuggestedAddress).Methods("GET")
	router.HandleFunc("/api/getDataAPI/{searchstring}", SuggestedAddress).Queries("searchstring", "{searchstring}")
	log.Fatal(http.ListenAndServe(":8081", router))
}
