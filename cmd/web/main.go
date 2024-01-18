package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	config "github.com/thepetk/snippetbox/cmd/web/config"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	snippets     *Snippet
	cfg          *config.Config
	queryTimeout int
}

func main() {
	cfg := cpnfig.Config{}
	cfgManager := config.ConfigManager{}
	flag.StringVar(&cfg.Addr, "addr", "", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "", "Path to static assets")
	flag.StringVar(&cfg.DB.DBPass, "db-pass", "", "Database password")
	flag.StringVar(&cfg.DB.DBUser, "db-user", "", "Database username")
	flag.StringVar(&cfg.DB.DBUser, "db-name", "", "Database name")
	flag.BoolVar(&cfg.DB.DBSSLDisabled, "db-ssl-disabled", true, "Database ssl config")
	flag.BoolVar(&cfg.SetLimits, "set-limits", false, "Set database pool limits")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		queryTimeout: 15,
	}

	flag.Parse()

	addr := cfgManager.GetConfigVar(cfg.Addr, ":4000", "SNIPPETBOX_ADDR")
	staticAddr := cfgManager.GetConfigVar(cfg.Addr, "./ui/static/", "SNIPPETBOX_STATIC")
	dbPass := cfgManager.GetConfigVar(cfg.DB.DBPass, "", "SNIPPETBOX_DB_PASSWORD")
	dbUser := cfgManager.GetConfigVar(cfg.DB.DBUser, "", "SNIPPETBOX_DB_USERNAME")
	dbName := cfgManager.GetConfigVar(cfg.DB.DBUser, "snippetbox", "SNIPPETBOX_DB_USERNAME")
	if dbPass == "" || dbUser == "" {
		app.errorLog.Fatal("Name or Password not set for database connection")
	}

	db, err := cfgManager.OpenDB(cfgManager.GetDSN(dbUser, dbPass, dbName, cfg.DB.DBSSLDisabled), cfg.SetLimits)
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
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  app.routes(staticAddr),
	}

	infoLog.Printf("Starting server on %s", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
