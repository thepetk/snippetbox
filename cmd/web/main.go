package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	config "github.com/thepetk/snippetbox/cmd/config"
	models "github.com/thepetk/snippetbox/cmd/models"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	debugLog     *log.Logger
	snippets     *models.Snippet
	cfg          *config.Config
	queryTimeout int
}

func main() {
	// Initialize application
	app := &application{
		debugLog: log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime),
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Configure application
	app.cfg = &config.Config{}
	db, err := app.cfg.InitConfig()

	if err != nil {
		app.errorLog.Fatal(err)
	}

	if err != nil {
		log.Fatalln(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	srv := &http.Server{
		Addr:     app.cfg.GetAddr(),
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", app.cfg.GetAddr())
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
