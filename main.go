package main

import (
	"GoWeb/app"
	"GoWeb/config"
	"GoWeb/database"
	"GoWeb/models"
	"GoWeb/routes"
	"embed"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed templates static
var res embed.FS

func main() {
	// Create instance of App
	appLoaded := app.App{}

	// Load config file to application
	appLoaded.Config = config.LoadConfig()

	// Load templates
	appLoaded.Res = &res

	// Create logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			panic("Failed to create log directory")
		}
	}

	// Create log file and set output
	file, err := os.OpenFile("logs/"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(file)

	// Connect to database and run migrations
	appLoaded.Db = database.ConnectDB(&appLoaded)
	if appLoaded.Config.Db.AutoMigrate {
		err = models.RunAllMigrations(&appLoaded)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// Define Routes
	routes.GetRoutes(&appLoaded)
	routes.PostRoutes(&appLoaded)

	// Start server
	log.Println("Starting server and listening on " + appLoaded.Config.Listen.Ip + ":" + appLoaded.Config.Listen.Port)
	err = http.ListenAndServe(appLoaded.Config.Listen.Ip+":"+appLoaded.Config.Listen.Port, nil)
	if err != nil {
		log.Println(err)
		return
	}
}
