package utils

import (
	"fmt"
	"testing"
)

func TestConvertToQueryParams(t *testing.T) {
	params := make(map[string]interface{})
	params["name"] = "xudai"
	params["password"] = "123456"
	str := ConvertToQueryParams(nil)
	fmt.Println(str)
}