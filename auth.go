package worktile

import (
	"encoding/json"
	"github.com/xudai3/worktile/utils"
)

type GetAuthCodeReq struct {
	ResponseType string `json:"response_type"`
	ClientId string `json:"client_id"`
	RedirectUri string `json:"redirect_uri"`
	Scope string `json:"scope"`
}

type GetAuthCodeRsp struct {
	Code string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}

type GetAuthTokenReq struct {
	GrantType string `json:"grant_type"`
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri string `json:"redirect_uri"`
	Code string `json:"code"`
}

type GetAuthTokenRsp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshAuthTokenReq struct {
	GrantType string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RefreshAuthTokenRsp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

const (
	GrantTypeAuthCode = "authorization_code"
	GrantTypeRefreshToken = "refresh_token"
)

func (w *Worktile) GetAuthCode(responseType string, scope string) string {
	req := GetAuthCodeReq{
		ResponseType: responseType,
		ClientId: w.ClientId,
		RedirectUri: w.RedirectUri,
		Scope: scope,
	}
	rsp := &GetAuthCodeRsp{}
	data, _ := w.Client.Get(ApiGetAuthCode, utils.ConvertStructToMap(req))
	json.Unmarshal(data, rsp)
	return rsp.Code
}

func (w *Worktile) GetAuthToken(code string) *AuthToken {
	req := GetAuthTokenReq{
		ClientId:     w.ClientId,
		ClientSecret: w.ClientSecret,
		RedirectUri:  w.RedirectUri,
		Code:         code,
		GrantType:    GrantTypeAuthCode,
	}
	rsp := &AuthToken{}
	data, _ := w.Client.Post(ApiGetAuthToken, "", req)
	json.Unmarshal(data, rsp)
	return rsp
}

func (w *Worktile) RefreshAuthToken() *AuthToken {
	req := RefreshAuthTokenReq{
		ClientId: w.ClientId,
		ClientSecret: w.ClientSecret,
		GrantType: GrantTypeRefreshToken,
		RefreshToken: w.Token.RefreshToken,
	}
	rsp := &AuthToken{}
	data, _ := w.Client.Post(ApiGetAuthToken, "", req)
	json.Unmarshal(data, rsp)
	return rsp
}