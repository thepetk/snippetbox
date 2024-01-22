package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	models "github.com/thepetk/snippetbox/cmd/models"
)

type application struct {
	logger   *log.Logger
	snippets *models.Snippet
	cfg      *Config
}

func main() {
	// Initialize application
	app := &application{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	// Configure application
	app.cfg = &Config{}
	db, err := app.InitConfig()

	if err != nil {
		app.logger.Fatal("ERROR\t", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			app.logger.Fatalln("ERROR\t", err)
		}
	}(db)

	srv := &http.Server{
		Addr:     app.GetAddr(),
		ErrorLog: app.logger,
		Handler:  app.routes(),
	}

	app.Log("INFO", fmt.Sprintf("Starting server on %s", app.GetAddr()))
	err = srv.ListenAndServe()
	app.logger.Fatal("ERROR\t", err)
}
