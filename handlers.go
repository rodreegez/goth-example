package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/markbates/goth/gothic"
)

type User struct {
	Email       string
	Name        string
	NickName    string
	Description string
	UserID      string
	AvatarURL   string
	Location    string
	AccessToken string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	view.HTML(w, 200, "index", getCurrentUser)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	newUser := User{
		user.Email,
		user.Name,
		user.NickName,
		user.Description,
		user.UserID,
		user.AvatarURL,
		user.Location,
		user.AccessToken,
	}

	setCurrentUser(w, r, user.UserID)

	view.JSON(w, http.StatusOK, newUser)
}

func setCurrentUser(w http.ResponseWriter, r *http.Request, id string) {
	session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
	session.Values["current_user"] = id
	session.Save(r, w)
}
func getCurrentUser(r *http.Request) string {
	session, _ := store.Get(r, os.Getenv("SESSION_NAME"))
	return session.Values["current_user"].(string)
}
