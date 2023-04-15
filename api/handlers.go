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
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := app.Models.UserDao.GetByUsername(input.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if ok, _ := user.Password.Matches(input.Password); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := app.Auth.GenerateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch token.Kind {
	case auth.Cookie:
		http.SetCookie(w, &http.Cookie{Name: auth.Name, Value: token.Value, HttpOnly: true, MaxAge: int(auth.TTL.Seconds())})
		w.WriteHeader(http.StatusAccepted)
		return
	case auth.Body:
		app.writeJSON(w, http.StatusAccepted, envelope{auth.Name: token.Value}, nil)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello " + r.Context().Value(auth.UserContextKey).(*data.User).Username))
}
