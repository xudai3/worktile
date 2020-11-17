package worktile

import (
	"encoding/json"
	"fmt"
	"github.com/xudai3/worktile/utils"
)

type UserDetailReq struct {
	AccessToken string `json:"access_token"`
	Uids string `json:"uids"`
}

type UserDetail struct {
	Uid string `json:"uid"`
	Name string `json:"name"`
	DisplayName string `json:"display_name"`
	Avatar string `json:"avatar"`
	Status int `json:"status"`
	DisplayNamePinyin string `json:"display_name_pinyin"`
}

func (w *Worktile) GetUserByUid(accessToken string, uid string) UserDetail {
	req := UserDetailReq{AccessToken:accessToken, Uids:uid}
	var rsp []UserDetail
	bytes, err := w.Client.Get(ApiGetUserByUid, utils.ConvertStructToMap(req))
	if err != nil {
		fmt.Printf("get user by uid:%s failed:%v\n", uid, err)
		return UserDetail{}
	}
	json.Unmarshal(bytes, &rsp)
	return rsp[0]
}