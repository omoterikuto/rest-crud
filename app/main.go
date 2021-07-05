package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
)

type Recipe struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	MakingTime  string `json:"making_time"`
	Serves      string `json:"serves"`
	Ingredients string `json:"ingredients"`
	Cost        string `json:"cost"`
}

const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	DataSourceName       = "root:golang@tcp(mysql-container:3306)/rest_crud"
	HerokuDataSourceName = "b441499201432f:a46600a6@tcp(us-cdbr-east-04.cleardb.com:3306)/rest_crud?parseTime=true"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()
	e.GET("/recipes", GetAllRecipes)
	e.GET("/recipes/:id", GetRecipe)
	e.POST("/recipes", CreateRecipe)
	e.PATCH("/recipes/:id", UpdateRecipe)
	e.DELETE("/recipes/:id", DeleteRecipe)
	e.Logger.Fatal(e.Start(":" + port))
}

var Update struct {
	Title       string
	MakingTime  string `json:"making_time"`
	Serves      string
	Ingredients string
	Cost        string
}

func UpdateRecipe(c echo.Context) error {
	id := c.Param("id")

	recipe := Update
	if bindErr := c.Bind(&recipe); bindErr != nil {
		log.Print("error bind to struct:", bindErr)
	}

	items := reflect.TypeOf(recipe)
	values := reflect.ValueOf(recipe)
	num := values.NumField()

	var datasource string
	if os.Getenv("DATABASE_URL") != "" {
		// Heroku用
		datasource = HerokuDataSourceName
	} else {
		// ローカル用
		datasource = DataSourceName
	}
	db, dbErr := sql.Open("mysql", datasource)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()

	for i := 0; i < num; i++ {
		field := items.Field(i)
		value := values.FieldByName(field.Name).String()
		if value == "" {
			continue
		}
		if field.Name == "MakingTime" {
			field.Name = "Making_Time"
		}

		sql := "UPDATE recipes SET " + strings.ToLower(field.Name) + " = ? WHERE id = ?"
		update, updateErr := db.Prepare(sql)
		if updateErr != nil {
			log.Fatal(updateErr)
		}
		update.Exec(value, id)
	}

	dr := new(Recipe)
	if selectErr := db.QueryRow("SELECT id, title, making_time, serves, ingredients, cost FROM recipes ORDER BY updated_at DESC").Scan(&dr.Id, &dr.Title, &dr.MakingTime, &dr.Serves, &dr.Ingredients, &dr.Cost); selectErr != nil {
		log.Print("query error :", selectErr)
	}

	outPut := new(Get)

	outPut.Message = "Recipe details by id"
	outPut.Recipe = append(outPut.Recipe, Recipe{dr.Id, dr.Title, dr.MakingTime, dr.Serves, dr.Ingredients, dr.Cost})

	return c.JSONPretty(
		http.StatusOK,
		outPut,
		"  ",
	)
}

type DeleteMessage struct {
	Message string `json:"message"`
}

func DeleteRecipe(c echo.Context) error {
	id := c.Param("id")
	dlMessage := new(DeleteMessage)

	var datasource string
	if os.Getenv("DATABASE_URL") != "" {
		// Heroku用
		datasource = HerokuDataSourceName
	} else {
		// ローカル用
		datasource = DataSourceName
	}
	db, dbErr := sql.Open("mysql", datasource)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()

	err := db.QueryRow("SELECT * FROM recipes WHERE id=?", id).Scan()
	if err == sql.ErrNoRows {
		dlMessage.Message = "No Recipe found"
		return c.JSONPretty(
			http.StatusExpectationFailed,
			dlMessage,
			"  ",
		)
	}

	db.QueryRow("DELETE FROM recipes WHERE id=?", id)

	dlMessage.Message = "Recipe successfully removed!"

	return c.JSONPretty(
		http.StatusOK,
		dlMessage,
		"  ",
	)
}

var ReceiveJson struct {
	Title       string `json:"title"`
	MakingTime  string `json:"making_time"`
	Serves      string `json:"serves"`
	Ingredients string `json:"ingredients"`
	Cost        string `json:"cost"`
}

type CreateFailed struct {
	Message  string `json:"message"`
	Required string `json:"required"`
}

func CreateRecipe(c echo.Context) error {
	recipe := ReceiveJson
	fail := new(CreateFailed)
	fail.Message = "Recipe creation failed!"
	fail.Required = "title, making_time, serves, ingredients, cost"

	if bindErr := c.Bind(&recipe); bindErr != nil {
		log.Print("error bind to struct:", bindErr)
	}

	items := reflect.TypeOf(recipe)
	values := reflect.ValueOf(recipe)
	num := values.NumField()

	for i := 0; i < num; i++ {
		field := items.Field(i)
		value := values.FieldByName(field.Name).String()
		if value == "" {
			return c.JSONPretty(
				http.StatusExpectationFailed,
				fail,
				"  ",
			)
		}
	}

	var datasource string
	if os.Getenv("DATABASE_URL") != "" {
		// Heroku用
		datasource = HerokuDataSourceName
	} else {
		// ローカル用
		datasource = DataSourceName
	}
	db, dbErr := sql.Open("mysql", datasource)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()

	insert, insertErr := db.Prepare("INSERT INTO recipes(title, making_time, serves, ingredients, cost) VALUES(?, ?, ?, ?, ?)")
	if insertErr != nil {
		log.Fatal(insertErr)
	}
	insert.Exec(recipe.Title, recipe.MakingTime, recipe.Serves, recipe.Ingredients, recipe.Cost)

	dr := new(Recipe)
	if selectErr := db.QueryRow("SELECT id, title, making_time, serves, ingredients, cost FROM recipes ORDER BY created_at DESC").Scan(&dr.Id, &dr.Title, &dr.MakingTime, &dr.Serves, &dr.Ingredients, &dr.Cost); selectErr != nil {
		log.Print("query error :", selectErr)
	}

	outPut := new(Get)

	outPut.Message = "Recipe successfully created!"
	outPut.Recipe = append(outPut.Recipe, Recipe{dr.Id, dr.Title, dr.MakingTime, dr.Serves, dr.Ingredients, dr.Cost})

	return c.JSONPretty(
		http.StatusOK,
		outPut,
		"  ",
	)
}

type Get struct {
	Message string   `json:"message"`
	Recipe  []Recipe `json:"recipe"`
}

func GetRecipe(c echo.Context) error {
	id := c.Param("id")

	var datasource string
	if os.Getenv("DATABASE_URL") != "" {
		// Heroku用
		datasource = HerokuDataSourceName
	} else {
		// ローカル用
		datasource = DataSourceName
	}
	db, dbErr := sql.Open("mysql", datasource)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()

	rows, queryErr := db.Query("SELECT id, title, making_time, serves, ingredients, cost FROM recipes where id = ?", id)
	if queryErr != nil {
		log.Print("query error :", queryErr)
	}
	defer rows.Close()

	var outPut Get

	for rows.Next() {
		var recipe Recipe

		if scanErr := rows.Scan(&recipe.Id, &recipe.Title, &recipe.MakingTime, &recipe.Serves, &recipe.Ingredients, &recipe.Cost); scanErr != nil {
			log.Fatal(scanErr)
		}
		outPut.Message = "Recipe details by id"
		outPut.Recipe = append(outPut.Recipe, Recipe{recipe.Id, recipe.Title, recipe.MakingTime, recipe.Serves, recipe.Ingredients, recipe.Cost})
	}

	return c.JSONPretty(
		http.StatusOK,
		outPut,
		"  ",
	)
}

type GetAll struct {
	Recipes []Recipe `json:"recipes"`
}

func GetAllRecipes(c echo.Context) error {
	var datasource string
	if os.Getenv("DATABASE_URL") != "" {
		// Heroku用
		datasource = HerokuDataSourceName
	} else {
		// ローカル用
		datasource = DataSourceName
	}
	db, dbErr := sql.Open("mysql", datasource)
	if dbErr != nil {
		log.Print("error connecting to database:", dbErr)
	}
	defer db.Close()

	rows, queryErr := db.Query("SELECT id, title, making_time, serves, ingredients, cost FROM recipes")
	if queryErr != nil {
		log.Print("query error :", queryErr)
	}

	defer rows.Close()

	var outPut GetAll

	for rows.Next() {
		var recipe Recipe

		if err := rows.Scan(&recipe.Id, &recipe.Title, &recipe.MakingTime, &recipe.Serves, &recipe.Ingredients, &recipe.Cost); err != nil {
			log.Fatal(err)
		}

		outPut.Recipes = append(outPut.Recipes, Recipe{recipe.Id, recipe.Title, recipe.MakingTime, recipe.Serves, recipe.Ingredients, recipe.Cost})
	}

	return c.JSONPretty(
		http.StatusOK,
		outPut,
		"  ",
	)
}
