package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"gopkg.in/unrolled/render.v1"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	twitterKey := os.Getenv("TWITTER_KEY")
	twitterSecret := os.Getenv("TWITTER_SECRET")

	goth.UseProviders(
		twitter.New(twitterKey, twitterSecret, "http://localhost:5000/auth/twitter/callback"),
	)

	v := render.New(render.Options{
		Directory:     "templates",
		Extensions:    []string{".html"},
		IsDevelopment: true,
	})

	p := pat.New()

	p.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		v.HTML(w, http.StatusOK, "user", user)
	})

	p.Get("/auth/{provider}", gothic.BeginAuthHandler)
	p.Get("/", func(w http.ResponseWriter, r *http.Request) {
		v.HTML(w, 200, "index", nil)
	})
	http.ListenAndServe(":5000", p)
}
