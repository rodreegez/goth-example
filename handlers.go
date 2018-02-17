package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/markbates/goth/gothic"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(r)
	if user == "" {
		user = "(none)"
	}
	view.HTML(w, http.StatusOK, "index", user)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	setCurrentUser(w, r, user.Name)

	//view.JSON(w, http.StatusOK, user)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	setCurrentUser(w, r, "")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func setCurrentUser(w http.ResponseWriter, r *http.Request, id string) {
	session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
	if id == "" {
		delete(session.Values, "current_user")
	} else {
		session.Values["current_user"] = id
	}
	session.Save(r, w)
}

func getCurrentUser(r *http.Request) string {
	session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
	if id, ok := session.Values["current_user"]; ok {
		if sid, ok := id.(string); ok {
			return sid
		}
	}
	return ""
}
