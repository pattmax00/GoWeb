package main

import (
	"GoWeb/app"
	"GoWeb/config"
	"GoWeb/database"
	"GoWeb/routes"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Create instance of App
	app := app.App{}

	// Load config file to application
	app.Config = config.LoadConfig()

	// Set log file
	file, _ := os.Create("logs/log-" + time.Now().String() + ".log")
	log.SetOutput(file)

	// Connect to database
	app.Db = database.ConnectDB(&app)

	// Define Routes
	routes.GetRoutes(&app)
	routes.PostRoutes(&app)

	// Start server
	log.Println("Starting server and listening on " + app.Config.Listen.Ip + ":" + app.Config.Listen.Port)
	err := http.ListenAndServe(app.Config.Listen.Ip+":"+app.Config.Listen.Port, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
