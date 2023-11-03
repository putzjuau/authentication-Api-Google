package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
    googleOauthConfig = &oauth2.Config{
        RedirectURL:  "http://localhost:8080/callback",
        ClientID:     "389812348545-fhn4sn60n0650f42rg3b9g24cohnii26.apps.googleusercontent.com",
        ClientSecret: "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",
        Scopes:       []string{"profile", "email"},
        Endpoint:     google.Endpoint,
    }

    // Variável de estado gerada aleatoriamente para segurança.
	oauthStateString, err = generateRandomString(32)
)

func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/callback", handleCallback)

    fmt.Println("Servidor executando na porta :8080")
    http.ListenAndServe(":8080", nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
    var htmlIndex = `<html><body><a href="/login">Login with Google</a></body></html>`
    fmt.Fprintf(w, htmlIndex)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    url := googleOauthConfig.AuthCodeURL(oauthStateString)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
    if state != oauthStateString {
        fmt.Println("Estado inválido")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    code := r.FormValue("code")
    token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
    if err != nil {
        fmt.Println("Erro ao trocar o código por token: ", err)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    // Agora você pode usar o token para fazer solicitações à API do Google em nome do usuário autenticado.
    // Por exemplo, para obter informações do usuário:
    // client := googleOauthConfig.Client(oauth2.NoContext, token)
    // plusService, err := plus.New(client)
    // user, err := plusService.People.Get("me").Do()

    fmt.Fprintf(w, "Token de acesso: %s\n", token.AccessToken)
}

func generateRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
