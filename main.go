package main

import (
	"GoWeb/app"
	"GoWeb/app/models"
	"GoWeb/app/routes"
	"GoWeb/config"
	"GoWeb/database"
	"GoWeb/templating"
	"context"
	"embed"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed app/frontend/templates app/frontend/static
var res embed.FS

func main() {
	// Create instance of Deps
	appLoaded := app.Deps{}

	// Load config file to application
	appLoaded.Config = config.LoadConfig()

	// Load templates
	appLoaded.Res = &res

	// Create logs directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			panic("failed to create log directory: " + err.Error())
		}
	}

	// Create log file and set output
	file, err := os.OpenFile("logs/"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("error creating log file: " + err.Error())
	}

	logger := slog.New(slog.NewTextHandler(file, nil))
	slog.SetDefault(logger) // Set structured logger globally

	// Connect to database and run migrations
	appLoaded.Db = database.Connect(&appLoaded)
	if appLoaded.Config.Db.AutoMigrate {
		err = models.RunAllMigrations(&appLoaded)
		if err != nil {
			slog.Error("error running migrations: " + err.Error())
			os.Exit(1)
		}
	}

	// Assign and run scheduled tasks
	appLoaded.ScheduledTasks = app.Scheduled{
		EveryReboot: []func(app *app.Deps){models.ScheduledSessionCleanup},
		EveryMinute: []func(app *app.Deps){models.ScheduledSessionCleanup},
	}

	// Define Routes
	routes.Get(&appLoaded)
	routes.Post(&appLoaded)

	// Prepare templates
	err = templating.BuildPages(&appLoaded)
	if err != nil {
		slog.Error("error building templates: " + err.Error())
		os.Exit(1)
	}

	// Start server
	server := &http.Server{Addr: appLoaded.Config.Listen.Ip + ":" + appLoaded.Config.Listen.Port}
	go func() {
		slog.Info("starting server and listening on " + appLoaded.Config.Listen.Ip + ":" + appLoaded.Config.Listen.Port)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("could not listen on %s: %v\n", appLoaded.Config.Listen.Ip+":"+appLoaded.Config.Listen.Port, err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal and shut down the server
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	stop := make(chan struct{})
	go app.RunScheduledTasks(&appLoaded, 100, stop)

	<-interrupt
	slog.Info("interrupt signal received. Shutting down server...")

	err = server.Shutdown(context.Background())
	if err != nil {
		slog.Error("could not gracefully shutdown the server: %v\n", err)
		os.Exit(1)
	}
}
