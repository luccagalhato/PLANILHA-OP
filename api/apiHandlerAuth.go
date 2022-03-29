package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	c "op/config"
	"op/models"
	"strings"

	auth "github.com/korylprince/go-ad-auth/v3"
)

// Tokens ...
var Tokens map[string]byte = map[string]byte{}

//Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	domain := c.Yml.AUTH.Server
	baseDN := c.Yml.AUTH.BaseDN
	client := models.Login{}
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Bad Request",
			Error:  "",
			Data:   err.Error(),
		})
		return
	}

	user := strings.Split(client.Username, `\`)
	username := client.Username
	if len(user) > 1 {
		username = strings.Join(user[1:], "\\")
		domain = user[0] + "." + domain
		baseDN = fmt.Sprintf("DC=%s,%s", user[0], baseDN)
	}
	config := &auth.Config{
		// Server: "alfa.local",
		Server: domain,
		Port:   c.Yml.AUTH.Port,
		BaseDN: baseDN,
		// BaseDN:   "DC=alfa,DC=local",
		Security: auth.SecurityNone,
	}

	_, entries, _, err := auth.AuthenticateExtended(config, username, client.Userpassword, []string{"memberOf"}, nil)
	block := true
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Unauthorized",
			Error:  "",
			Data:   err.Error(),
		})
		return
	}
	fmt.Println(entries)
	if entries != nil {
		for _, value := range entries.GetAttributeValues("memberOf") {
			if strings.Contains(value, c.Yml.AUTH.Grupo) {
				block = false
				break
			}
		}
	}

	if block {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Response{
			Status: "Unauthorized",
			Error:  "",
			Data:   "Senha ou username inválidos",
		})
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	delete(Tokens, "1")
	token := genKey()
	Tokens[token] = 1
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Token", token)
	w.Header().Set("Username", client.Username)
	json.NewEncoder(w).Encode(models.Response{
		Status: "OK",
		Error:  "",
		Data:   "Liberado",
	})
}

//Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {

	if data, err := r.Cookie("Token"); err == nil {
		fmt.Println(data.Value)
		delete(Tokens, data.Value)
	} else {
		log.Println(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("access-control-expose-headers", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	fmt.Println(Tokens)
	json.NewEncoder(w).Encode(models.Response{
		Status: "OK",
		Error:  "",
	})
}
