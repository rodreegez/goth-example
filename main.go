package main

import (
	"log"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"gopkg.in/unrolled/render.v1"
)

var view = render.New(render.Options{
	Directory:     "templates",
	Extensions:    []string{".html"},
	IsDevelopment: true,
	IndentJSON:    true,
})

var port = "5000"
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port = os.Getenv("PORT")
	baseURL := "http://127.0.0.1:" + port

	goth.UseProviders(
		twitter.New(
			os.Getenv("TWITTER_KEY"),
			os.Getenv("TWITTER_SECRET"),
			baseURL+"/auth/twitter/callback"),
	)
}

func main() {
	p := pat.New()

	p.Get("/auth/{provider}/callback", CallbackHandler)
	p.Get("/auth/{provider}", gothic.BeginAuthHandler)
	p.Get("/signout", SignoutHandler)
	p.Get("/", IndexHandler)

	n := negroni.Classic()
	n.UseHandler(p)
	n.Run(":" + port)
}
