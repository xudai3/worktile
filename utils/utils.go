package utils

import (
	"encoding/json"
)

func ConvertStructToMap(item interface{}) map[string]interface{} {
	data, _ := json.Marshal(item)
	resultMap := make(map[string]interface{})
	json.Unmarshal(data, &resultMap)
	return resultMap
}
