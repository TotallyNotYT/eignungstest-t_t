package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/atselvan/pokeapi"
	_ "github.com/lib/pq"
)

var connStr string

func init() {
	var host = os.Getenv("DB_HOST")

	if host == "" {
		host = "localhost"
	}

	connStr = fmt.Sprintf("user=postgres password=admin host=%v port=5432 dbname=postgres sslmode=disable", host)
}

func main() {

	run()

	http.HandleFunc("/", handler)
	http.HandleFunc("/load", loadyes)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func loadyes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	rows, err := db.Query("SELECT * FROM pokemon")
	if err != nil {
		log.Fatal("Oopsie", err)
	}

	defer rows.Close()

	var stringArray []string

	type Pokemon struct {
		Name string
	}

	for rows.Next() {
		var pokemon Pokemon
		if err := rows.Scan(&pokemon.Name); err != nil {
			return
		}
		stringArray = append(stringArray, pokemon.Name)
	}

	jsonData, err := json.Marshal(stringArray)
	if err != nil {
		http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

func run() {
	client := pokeapi.NewClient()

	result, restErr := client.Resource.Berries()

	if restErr != nil {
		log.Fatal("Fatal")
		log.Fatal(restErr)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	db.Exec("DROP TABLE IF EXISTS pokemon")
	db.Exec("CREATE TABLE IF NOT EXISTS pokemon (name VARCHAR(255))")

	slicedPokemon := (*result)[0:15]

	stmt, err := db.Prepare("INSERT INTO pokemon (name) VALUES ($1)")
	if err != nil {
		log.Fatal(err)
	}
	for _, pokemon := range slicedPokemon {
		_, err = stmt.Exec(pokemon.Name)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err != nil {
		log.Fatal("oh yeah", err)
	}

	for i := 0; i < len(*result); i++ {
		fmt.Printf("%+v\n", (*result)[i].Name)
	}
}
