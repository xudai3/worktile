package worktile

import (
	"github.com/xudai3/worktile/utils"
	"net/http"
)

type Worktile struct {
	//Token *AuthToken `json:"token"`
	Client *utils.HttpClient `json:"client"`
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri string `json:"redirect_uri"`
}

func NewWorktile(clientId string, clientSecret string, redirectUri string) *Worktile {
	w := &Worktile{
		Client: utils.NewHttpClient(&http.Client{}),
		ClientId: clientId,
		ClientSecret: clientSecret,
		RedirectUri: redirectUri,
	}
	return w
}