package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ConvertStructToMap(item interface{}) map[string]interface{} {
	data, _ := json.Marshal(item)
	resultMap := make(map[string]interface{})
	json.Unmarshal(data, &resultMap)
	return resultMap
}

// AggregateError 用于表示包含多个子错误的聚合错误。
type AggregateError struct {
	errors []error
}

// NewAggregateError 创建一个新的聚合错误实例。
func NewAggregateError(errors []error) *AggregateError {
	return &AggregateError{errors: errors}
}

// Error 实现error接口，返回所有子错误的字符串描述。
func (ae *AggregateError) Error() string {
	var errMsgs []string
	for _, err := range ae.errors {
		errMsgs = append(errMsgs, err.Error())
	}
	return fmt.Sprintf("Encountered %d errors: %s", len(ae.errors), strings.Join(errMsgs, "; "))
}

// Append 向聚合错误中追加新的错误。
func (ae *AggregateError) Append(err error) {
	if err != nil {
		ae.errors = append(ae.errors, err)
	}
}

// CollectErrors 收集错误并将它们转换为一个聚合错误。
func CollectErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return NewAggregateError(errs)
}
