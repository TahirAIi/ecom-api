package main

import (
	"database/sql"
	data "ecom-api/inernal/data/models"
	"errors"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(int64(app.config.multipartFormSize))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(name) < 1 {
		app.sendResponse(w, response{"message": "Name is required"}, http.StatusUnprocessableEntity)
		return
	}

	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9.!#$%&'*+/=?^_\x60\{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`, email); !ok {
		app.sendResponse(w, response{"message": "Email is invalid"}, http.StatusUnprocessableEntity)
		return
	}

	if len(password) < 8 || len(password) > 72 {
		app.sendResponse(w, response{"message": "Password should be between 8 and 72 characters"}, http.StatusUnprocessableEntity)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}

	user, err := app.models.User.GetByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}
	if user != nil {
		app.sendResponse(w, response{"message": "User already exists with this email"}, http.StatusConflict)
		return
	}
	user = &data.User{
		FullName: name,
		Email:    email,
		Password: string(hash),
	}
	err = app.models.User.Insert(user)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}
	app.sendResponse(w, response{"user": user}, http.StatusCreated)
	return

}


//swagger:route POST /authorize Auth authUser
//Creates a token.
//
//Consumes:
//	- multipart/form-data
//
//Produces:
//	- application/json
//Parameters:
//	AuthBody
//responses:
//	200: AuthResponse
//swagger:response
func (app *application) authUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(int64(app.config.multipartFormSize))
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := app.models.User.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.sendResponse(w, response{"message": "Invalid email or password"}, http.StatusUnauthorized)
			return
		default:
			app.log(err)
			app.sendInternalServerErrorResponse(w)
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		app.sendResponse(w, response{"message": "Invalid email or password"}, http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.Email, user.FullName)
	if err != nil {
		app.log(err)
		app.sendInternalServerErrorResponse(w)
		return
	}
	app.sendResponse(w, response{"token": token}, http.StatusCreated)
}
