package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/lib/pq"
	models "github.com/thepetk/snippetbox/cmd/models"
)

type application struct {
	logger        *log.Logger
	snippets      *models.SnippetModel
	cfg           *Config
	templateCache map[string]*template.Template
}

func main() {
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		return
	}
	app := &application{
		logger:        log.New(os.Stdout, "", log.Ldate|log.Ltime),
		snippets:      &models.SnippetModel{},
		templateCache: templateCache,
	}

	// Configure application
	app.cfg = &Config{}
	db, err := app.InitConfig()
	app.snippets.SetDB(db)

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
