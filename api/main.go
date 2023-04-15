package main

import (
	"database/sql"
	"net/http"

	"github.com/OmarTariq612/web-auth/auth"
	"github.com/OmarTariq612/web-auth/data"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	Auth   auth.Auth
	Models data.Models
}

func (app *application) serve() error {
	srv := http.Server{
		Addr:    "localhost:5555",
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	models := data.NewModels(db)

	// app := &application{Auth: auth.NewBasicAuth(models.UserDao), Models: models}
	// app := &application{Auth: auth.NewStatefulAuth(models.UserDao, models.TokenDao), Models: models}
	app := &application{Auth: auth.NewStatelessAuth(models.UserDao), Models: models}

	if err := app.serve(); err != nil {
		panic(err)
	}
}
