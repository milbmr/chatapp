package main

import "net/http"

type auth struct {
	next http.Handler
}

func (a *auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("userdata")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.next.ServeHTTP(w, r)
}

func AuthMiddlware(handler http.Handler) http.Handler {
	return &auth{next: handler}
}
