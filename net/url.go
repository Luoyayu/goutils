package net

import (
	"fmt"
	"net/url"
	"strings"
)

func JoinPath(base string, path []interface{}) (*url.URL, error) {
	a := []string{strings.TrimSuffix(base, "/")}
	for i, p := range path {
		sp := fmt.Sprint(p)
		if i == 0 {
			a = append(a, strings.TrimPrefix(sp, "/"))
		} else {
			a = append(a, sp)
		}
	}
	b := strings.Join(a, "/")
	b = strings.TrimRight(b, "/")
	return url.Parse(b)
}
