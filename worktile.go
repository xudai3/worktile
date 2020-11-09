package worktile

import (
	"github.com/xudai3/worktile/utils"
	"net/http"
)

type Worktile struct {
	Token *AuthToken `json:"token"`
	Client *utils.HttpClient
}

func NewWorktile() *Worktile {
	w := &Worktile{
		Client: utils.NewHttpClient(&http.Client{}),
	}
	return w
}