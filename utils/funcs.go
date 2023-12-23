package utils

import (
	"fmt"
	"math/rand"
	"strconv"
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
	"randInt": func(sep string, n ...int) string {
		result := make([]string, len(n))
		for i := range n {
			result[i] = strconv.Itoa(rand.Intn(n[i]))
		}

		return strings.Join(result, sep)
	},
	"join": func(sep string, s ...any) string {
		result := make([]string, len(s))
		for i := range s {
			result[i] = fmt.Sprintf("%v", s[i])
		}
		return strings.Join(result, sep)
	},
}
