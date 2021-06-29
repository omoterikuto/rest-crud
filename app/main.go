package main

import (
	"fmt"
	"database/sql"
	"log"
	// "net/http"
	// "time"
	_ "github.com/go-sql-driver/mysql"
)

type Recipe struct {
	ID          int
	Title       string
	CreatedAt   string
}
const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	DataSourceName = "root:golang@tcp(mysql-container:3306)/rest_crud"
)

var rcp = make(map[int]Recipe)

func main() {
	db, dbErr := sql.Open(DriverName, DataSourceName)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()
	rows, queryErr := db.Query("SELECT * FROM recipes")
	if queryErr != nil {
			log.Print("query error :", queryErr)
	}
	defer rows.Close()
	fmt.Println(rows)
}
