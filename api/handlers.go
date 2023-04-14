package main

import (
	"net/http"

	"github.com/OmarTariq612/web-auth/auth"
	"github.com/OmarTariq612/web-auth/data"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := data.User{Username: input.Username}
	if err := user.Password.Set(input.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := app.Models.UserDao.Insert(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello " + r.Context().Value(auth.UserContextKey).(*data.User).Username))
}
