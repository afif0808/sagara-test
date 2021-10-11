package structs

import "reflect"

// Utility package for stuffs related to struct
// inspired by : https://dave.cheney.net/practical-go/presentations/qcon-china.html#_avoid_package_names_like_base_common_or_util

func GetStructTagValues(data interface{}, tagName string) []string {
	var values []string
	v := reflect.ValueOf(data)
	n := v.Type().NumField()
	for i := 0; i < n; i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" {
			continue
		}
		values = append(values, tag)
	}
	return values
}
