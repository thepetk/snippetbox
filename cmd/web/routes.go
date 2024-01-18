package main

import "net/http"

func (app *application) routes(staticAddr string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	fileServer := http.FileServer(http.Dir(staticAddr))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
