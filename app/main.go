package main

import (
	"fmt"
	"database/sql"
	"log"
	"encoding/json"
	// "net/http"
	// "time"
	_ "github.com/go-sql-driver/mysql"
)

type Recipe struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	MakingTime  string `json:"making_time"`
	Serves		  string `json:"serves"`
	Ingredients string `json:"ingredients"`
	Cost        string `json:"cost"`
}

type OutPutAll struct {
	Recipes []Recipe
}

const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	DataSourceName = "root:golang@tcp(mysql-container:3306)/rest_crud"
)

var rcp = make(map[int]Recipe)

func main() {
	
	getAll()
}

func getAll() {
	db, dbErr := sql.Open(DriverName, DataSourceName)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()
	rows, queryErr := db.Query("SELECT id, title, making_time, serves, ingredients, cost FROM recipes")
	if queryErr != nil {
			log.Print("query error :", queryErr)
	}
	defer rows.Close()

	var outPut OutPutAll
	
	for rows.Next() {
		var recipe Recipe

    if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.MakingTime, &recipe.Serves, &recipe.Ingredients, &recipe.Cost); err != nil {
      log.Fatal(err)
    }

		outPut.Recipes = append(outPut.Recipes, Recipe{recipe.Id, recipe.Title, recipe.MakingTime, recipe.Serves, recipe.Ingredients, recipe.Cost})
	}

	outPutJson, err := json.Marshal(outPut)
	if err != nil {
		fmt.Printf("json変換エラー", err)
	}

	fmt.Println(string(outPutJson))
}
