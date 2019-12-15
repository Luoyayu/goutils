package gadio

import (
	"testing"
)

func TestGetLatestN(t *testing.T) {

	t.Run("GetLatestN", func(t *testing.T) {
		if ret, err := GetLatestN(5); err == nil {
			if err := Query(ret); err == nil {
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	})
}
