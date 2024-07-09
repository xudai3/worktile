package worktile

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xudai3/worktile/utils"
)

func (w *Worktile) GetUserByUid(uid string) (*UserDetail, error) {
	accessToken, err := w.GetTenant()
	if err != nil {
		return nil, err
	}
	req := GetUserDetailReq{AccessToken: accessToken, Uids: uid}
	var rsp []*UserDetail
	bytes, err := w.Client.Get(ApiGetUserByUid, utils.ConvertStructToMap(req))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUid>unmarshal user details error:%v", err)
	}
	if len(rsp) == 0 {
		return nil, errors.New("GetUserByUid>result empty")
	}
	return rsp[0], nil
}

func (w *Worktile) GetUsersByUids(uids []string) ([]*UserDetail, error) {
	var rsp []*UserDetail
	for _, uid := range uids {
		user, err := w.GetUserByUid(uid)
		if err != nil {
			return nil, err
		}
		if user != nil {
			rsp = append(rsp, user)
		}
	}
	if len(rsp) == 0 {
		return nil, errors.New("GetUsersByUids>result empty")
	}
	return rsp, nil
}
