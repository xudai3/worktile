package worktile

import (
	"testing"

	"github.com/xudai3/worktile/utils"
)

func TestConvertStructToMap(t *testing.T) {
	req := GetAuthCodeReq{
		ResponseType: "types",
		ClientId:     "ids",
		RedirectUri:  "uris",
		Scope:        "scopes",
	}
	res := utils.ConvertStructToMap(req)
	t.Log(res)
}
