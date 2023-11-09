package api

import (
	"fmt"
	"reflect"
	"strings"
)

func buildQueryParams(params map[string]interface{}) string {
	var queryParams []string
	for key, value := range params {
		if value == nil || value == "" {
			continue
		}
		// Check if the value is a slice and if it's empty
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			if reflect.ValueOf(value).Len() == 0 {
				continue
			}
		}
		queryParams = append(queryParams, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(queryParams, "&")
}
