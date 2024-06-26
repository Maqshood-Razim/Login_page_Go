package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

var (
	templates     = template.Must(template.ParseFiles("login.html", "home.html"))
	validUsername = "user"
	validPassword = "pass"
	sessionName   = "session-id"
)

func main() {
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/logout", logoutHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
	
		cookie, err := r.Cookie(sessionName)
		if err == nil && cookie.Value == "authenticated" {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		templates.ExecuteTemplate(w, "login.html", nil)
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == validUsername && password == validPassword {
			
			expiration := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{
				Name:     sessionName,
				Value:    "authenticated",
				Expires:  expiration,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			
			templates.ExecuteTemplate(w, "login.html", map[string]string{
				"Error": "Incorrect username or password",
			})
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionName)
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	w.Header().Set("Cache-Control", "no-store, must-revalidate")

	templates.ExecuteTemplate(w, "home.html", nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	
	cookie := http.Cookie{
		Name:     sessionName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
