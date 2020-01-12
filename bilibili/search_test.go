package biliAPI

import (
	"fmt"
	"testing"
)

func TestGetSearchUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		if ret, err := GetSearchUser("bilibili", 1); err == nil {
			for _, result := range ret.Data.Result {
				fmt.Println(result.UName)
			}
		} else {
			panic(err)
		}
	})
}

func TestGetSearchSuggest(t *testing.T) {
	t.Run("", func(t *testing.T) {
		if ret, err := GetSearchSuggest("哔哩哔哩弹幕网"); err == nil {
			for _, tag := range ret.Result.Tags {
				fmt.Println(tag.Name)
			}
		} else {
			panic(err)
		}
	})
}
