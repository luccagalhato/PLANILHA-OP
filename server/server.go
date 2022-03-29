package server

import (
	"embed"
	_ "embed" //haha
	"log"
	"net/http"
	"op/api"
	c "op/config"
	"strings"
	"time"
)

//go:embed html
var fs embed.FS

//Controllers ...
func Controllers() {
	// f, _ := fs.Open("login.html")
	// io.Copy(os.Stdout, f)
	// fmt.Println(fs.Open("html/login.html"))
	log.Printf("starting server at port: %s", c.Yml.API.Port)
	// http.Handle("/id", &Auth{api.DownloadExcell, true})
	http.Handle("/selectop", &Auth{api.SelectOP, true})
	http.HandleFunc("/login", api.Login)
	http.Handle("/logout", &Auth{api.Logout, true})
	fss := http.FileServer(http.FS(fs))
	http.Handle("/html/", http.StripPrefix("/", &Auth{fss.ServeHTTP, false}))
	http.HandleFunc("/", redirect)
	log.Fatal(http.ListenAndServe(":"+c.Yml.API.Port, nil))
}

func redirect(w http.ResponseWriter, r *http.Request) {

	// log.Println(r.URL.Path)
	paths := strings.Split(r.URL.Path, ".")
	if r.URL.Path != "/" && paths[len(paths)-1] != "html" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// remove/add not default ports from req.Host
	target := "http://" + r.Host + "/html/login.html"
	// log.Println(target)
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}
	// log.Printf("redirect to: %s", target)
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r.URL.Path)
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
	if a.all || (urlS[len(urlS)-1] == "html" && r.URL.String() != "html/login.html") {
		redirect(w, r)
		return
	}
	a.handler(w, r)
}
