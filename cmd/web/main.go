package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getConfigVar(flagAddr string, fallback string, envVar string) string {
	if len(flagAddr) == 0 {
		return getEnv(envVar, fallback)
	}
	return flagAddr
}

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.Addr, "addr", "", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "", "Path to static assets")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(app.home))
	mux.Handle("/snippet", http.HandlerFunc(app.showSnippet))
	mux.Handle("/snippet/create", http.HandlerFunc(app.createSnippet))

	flag.Parse()

	addr := getConfigVar(cfg.Addr, ":4000", "SNIPPETBOX_ADDR")
	staticAddr := getConfigVar(cfg.Addr, "./ui/static/", "SNIPPETBOX_STATIC")

	fileServer := http.FileServer(http.Dir(staticAddr))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
