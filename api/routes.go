package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", app.registerUser)
	mux.HandleFunc("/login", app.loginUser)
	mux.Handle("/", app.Auth.Authenticate(http.HandlerFunc(app.index)))

	return mux
}
