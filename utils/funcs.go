package utils

import (
	"strings"
	"text/template"
)

var FuncMaps = template.FuncMap{
	"split": func(s any, sep string, i int) string {
		switch v := s.(type) {
		case string:
			parts := strings.Split(v, sep)
			if len(parts) > i {
				return parts[i]
			}
			return parts[0]
		default:
			return ""
		}
	},
}
