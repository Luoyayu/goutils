package biliAPI

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserInfo(t *testing.T) {
	t.Run("", func(t *testing.T) {
		if ret, err := GetUserInfo(383768376); err == nil {
			assert.Equal(t, int64(383768376), ret.Data.Mid)
			fmt.Println(ret.Data.Name)
		} else {
			panic(err)
		}
	})
}
