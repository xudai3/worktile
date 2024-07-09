package worktile

import (
	"encoding/json"
	"time"

	"github.com/xudai3/worktile/utils"
)

type GetAuthCodeReq struct {
	ResponseType string `json:"response_type"`
	ClientId     string `json:"client_id"`
	RedirectUri  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

type GetAuthCodeRsp struct {
	Code  string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}

type GetAuthTokenReq struct {
	GrantType    string `json:"grant_type"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	Code         string `json:"code"`
}

type GetAuthTokenRsp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshAuthTokenReq struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RefreshAuthTokenRsp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type GetTenantReq struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Tenant struct {
	TenantAccessToken string `json:"tenant_access_token"`
	ExpiresIn         int64  `json:"expires_in"`
}

const (
	GrantTypeAuthCode     = "authorization_code"
	GrantTypeRefreshToken = "refresh_token"
)

func (w *Worktile) GetAuthCode(responseType string, scope string) string {
	req := GetAuthCodeReq{
		ResponseType: responseType,
		ClientId:     w.ClientId,
		RedirectUri:  w.RedirectUri,
		Scope:        scope,
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
	data, _ := w.Client.Post(ApiGetTenant, "", req)
	json.Unmarshal(data, rsp)
	return rsp
}

func (w *Worktile) GetTenant() (string, error) {
	tenentKey := "tenant_access_token"
	item, found := w.Cache.Get(tenentKey)
	if found {
		tenent, _ := item.(*Tenant)
		now := time.Now().Unix()
		// fmt.Printf("GetTenant>get cache:%v %v\n", tenent.TenantAccessToken, tenent.ExpiresIn)
		// 检查token是否过期
		if now < tenent.ExpiresIn {
			return tenent.TenantAccessToken, nil
		}
	}
	req := GetTenantReq{ClientId: w.ClientId, ClientSecret: w.ClientSecret}
	rsp := &Tenant{}
	data, err := w.Client.Post(ApiGetTenant, "", req)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, rsp)
	if err != nil {

		return "", err
	}
	expireTime := rsp.ExpiresIn
	rsp.ExpiresIn += time.Now().Unix()
	// fmt.Printf("GetTenant>set cache:%v %v %v\n", rsp.TenantAccessToken, rsp.ExpiresIn, expireTime)
	w.Cache.SetWithTTL(tenentKey, rsp, 1, time.Duration(expireTime)*time.Second)

	return rsp.TenantAccessToken, nil
}

func (w *Worktile) RefreshAuthToken(refreshToken string) *AuthToken {
	req := RefreshAuthTokenReq{
		ClientId:     w.ClientId,
		ClientSecret: w.ClientSecret,
		GrantType:    GrantTypeRefreshToken,
		RefreshToken: refreshToken,
	}
	rsp := &AuthToken{}
	data, _ := w.Client.Post(ApiGetTenant, "", req)
	json.Unmarshal(data, rsp)
	return rsp
}
