package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

const (
	oauthState     = "oauthState"
	oauthGoogleURL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type authUser struct {
	conf *oauth2.Config
}

func (a *authUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	seg := strings.Split(r.URL.Path, "/")
	action := seg[2]
	provider := seg[3]
	switch action {
	case "login":
		if provider != "google" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "auth action not supported %s", action)
			return
		}
		state := generateAuthState(w)
		fmt.Println(state)
		url := a.conf.AuthCodeURL(state)
		w.Header().Set("location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		stateCookie, err := r.Cookie(oauthState)
		if err != nil {
			log.Printf("error reading cookie %s", err.Error())
			return
		}
		if r.FormValue("state") != stateCookie.Value {
			log.Println("invalid google auth state")
			w.Header().Set("location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		token, err := a.conf.Exchange(context.Background(), r.FormValue("code"))
		if err != nil {
			log.Printf("error exchaning token %s", err.Error())
			w.Header().Set("location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		response, err := http.Get(oauthGoogleURL + token.AccessToken)
		if err != nil {
			log.Printf("error getting user data %s", err.Error())
			w.Header().Set("location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("error parsing response %s", err.Error())
			w.Header().Set("location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		fmt.Println(string(data))
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "auth action not supported %s", action)
	}
}

func generateAuthState(w http.ResponseWriter) string {
	expiration := time.Now().Add(20 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	http.SetCookie(w, &http.Cookie{Name: oauthState, Value: state, Path: "/", Expires: expiration})
	return state
}
