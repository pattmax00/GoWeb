package models

import (
	"GoWeb/app"
	"GoWeb/database"
)

// RunAllMigrations defines the structs that should be represented in the database
func RunAllMigrations(app *app.App) error {
	// Declare new dummy user for reflection
	user := User{
		Id:        1, // Id is handled automatically, but it is added here to show it will be skipped during column creation
		Username:  "migrate",
		Password:  "migrate",
		AuthToken: "migrate",
		CreatedAt: "2021-01-01 00:00:00",
		UpdatedAt: "2021-01-01 00:00:00",
	}

	return database.Migrate(app, user)
}
