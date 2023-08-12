package database

import (
	"GoWeb/app"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// Connect returns a new database connection
func Connect(app *app.App) *sql.DB {
	postgresConfig := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		app.Config.Db.Ip, app.Config.Db.Port, app.Config.Db.User, app.Config.Db.Password, app.Config.Db.Name)

	db, err := sql.Open("postgres", postgresConfig)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Connected to database successfully on " + app.Config.Db.Ip + ":" + app.Config.Db.Port + " using database " + app.Config.Db.Name)

	return db
}
