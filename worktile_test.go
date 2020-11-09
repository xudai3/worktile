package worktile

import (
	"fmt"
	"github.com/xudai3/worktile/utils"
	"testing"
)

func TestNewWorktile(t *testing.T) {
	req := GetAuthCodeReq{
		ResponseType: "types",
		ClientId: "ids",
		RedirectUri: "uris",
		Scope: "scopes",
	}
	res := utils.ConvertStructToMap(req)
	fmt.Println(res)
}