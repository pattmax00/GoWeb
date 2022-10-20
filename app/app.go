package app

import (
	"GoWeb/config"
	"database/sql"
)

// App contains and supplies available configurations and connections
type App struct {
	Config config.Configuration
	Db     *sql.DB
}
