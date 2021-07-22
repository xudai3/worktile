package worktile

import (
	"encoding/json"
	"errors"

	"github.com/xudai3/worktile/logger"
	"github.com/xudai3/worktile/utils"
)

func (w *Worktile) GetUserByUid(accessToken string, uid string) (*UserDetail, error) {
	req := GetUserDetailReq{AccessToken: accessToken, Uids: uid}
	var rsp []*UserDetail
	bytes, err := w.Client.Get(ApiGetUserByUid, utils.ConvertStructToMap(req))
	if err != nil {
		logger.Debugf("get user by uid:%s failed:%v\n", uid, err)
		return nil, err
	}
	err = json.Unmarshal(bytes, &rsp)
	if err != nil {
		logger.Errorf("unmarshal user details error:%v", err)
		return nil, err
	}
	if len(rsp) == 0 {
		return nil, errors.New("result empty")
	}
	return rsp[0], nil
}

func (w *Worktile) GetUsersByUids(accessToken string, uids []string) ([]*UserDetail, error) {
	var rsp []*UserDetail
	for _, uid := range uids {
		user, err := w.GetUserByUid(accessToken, uid)
		if err != nil {
			return nil, err
		}
		if user != nil {
			rsp = append(rsp, user)
		}
	}
	if len(rsp) == 0 {
		return nil, errors.New("result empty")
	}
	return rsp, nil
}
