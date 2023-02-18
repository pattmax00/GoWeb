package models

import (
	"GoWeb/app"
	"GoWeb/database"
	"time"
)

// RunAllMigrations defines the structs that should be represented in the database
func RunAllMigrations(app *app.App) error {
	// Declare new dummy user for reflection
	user := User{
		Id:        1, // Id is handled automatically, but it is added here to show it will be skipped during column creation
		Username:  "migrate",
		Password:  "migrate",
		AuthToken: "migrate",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := database.Migrate(app, user)
	if err != nil {
		return err
	}

	return nil
}
