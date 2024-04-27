package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.template.Execute(w, r)
}

func main() {
	conf := &oauth2.Config{
		ClientID:     "978105866416-d3t4r3ffju8oav74riahbpabj2ennndu.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:8080/auth/callback/google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	addr := flag.String("addr", ":8080", "the port of the app")
	flag.Parse()
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	http.Handle("/auth/", &authUser{conf: conf})

	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("listenandserve:", err)
	}
}
