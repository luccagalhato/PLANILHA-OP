package server

import (
	"log"
	"net/http"
	"op/api"
	c "op/config"
	"strings"
	"time"
)

//Controllers ...
func Controllers() {
	log.Printf("starting server at port: %s", c.Yml.API.Port)
	http.HandleFunc("/", redirect)
	// http.Handle("/id", &Auth{api.DownloadExcell, true})
	http.Handle("/selectop", &Auth{api.SelectOP, true})
	http.HandleFunc("/login", api.Login)
	http.Handle("/logout", &Auth{api.Logout, true})
	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/html/", http.StripPrefix("/html/", &Auth{fs.ServeHTTP, false}))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	// remove/add not default ports from req.Host
	target := "http://" + r.Host + "/html/login.html"
	// log.Println(target)
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	// log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target,
		// see comments below and consider the codes 308, 302, or 301
		http.StatusTemporaryRedirect)
}
func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")

	// if s[len(s)-1] == "html" {
	// fmt.Println(r.URL.String())
	// }

	if token, err := r.Cookie("Token"); err == nil {
		if _, ok := api.Tokens[token.Value]; ok {
			if token.Expires.Before(time.Now()) {
				a.handler(w, r)
				return
			}
			//deletar token vencido- verificar se ja deleta automatico
		}
	}

	urlS := strings.Split(r.URL.String(), ".")
	if a.all || (urlS[len(urlS)-1] == "html" && r.URL.String() != "login.html") {
		redirect(w, r)
		return
	}
	a.handler(w, r)
}
