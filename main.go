package main

import (
	"GoWeb/app"
	"GoWeb/config"
	"GoWeb/database"
	"GoWeb/models"
	"GoWeb/routes"
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
			log.Println("Failed to create log directory")
			log.Println(err)
			return
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

	// Assign and run scheduled tasks
	appLoaded.ScheduledTasks = app.Scheduled{
		EveryReboot: []func(app *app.App){models.ScheduledSessionCleanup},
		EveryMinute: []func(app *app.App){models.ScheduledSessionCleanup},
	}

	// Define Routes
	routes.GetRoutes(&appLoaded)
	routes.PostRoutes(&appLoaded)

	// Start server
	server := &http.Server{Addr: appLoaded.Config.Listen.Ip + ":" + appLoaded.Config.Listen.Port}
	go func() {
		log.Println("Starting server and listening on " + appLoaded.Config.Listen.Ip + ":" + appLoaded.Config.Listen.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", appLoaded.Config.Listen.Ip+":"+appLoaded.Config.Listen.Port, err)
		}
	}()

	// Wait for interrupt signal and shut down the server
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	stop := make(chan struct{})
	go app.RunScheduledTasks(&appLoaded, 100, stop)

	<-interrupt
	log.Println("Interrupt signal received. Shutting down server...")

	err = server.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
}
