package worktile

import (
	"net/http"

	"github.com/dgraph-io/ristretto"
	"github.com/xudai3/worktile/utils"
)

type Worktile struct {
	Cache        *ristretto.Cache  `json:"token"`
	Client       *utils.HttpClient `json:"client"`
	ClientId     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
	RedirectUri  string            `json:"redirect_uri"`
}

func NewWorktile(clientId string, clientSecret string, redirectUri string) *Worktile {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000, // number of keys to track frequency of (10M).
		MaxCost:     1000,  // maximum cost of cache (1GB).
		BufferItems: 64,    // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	w := &Worktile{
		Cache:        cache,
		Client:       utils.NewHttpClient(&http.Client{}, 3, 2),
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUri:  redirectUri,
	}
	return w
}
