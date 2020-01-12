package biliAPI

import (
	"fmt"
	"log"
	"math"
	"testing"
)

const _test_uid = 9321359

func TestGetFollows(t *testing.T) {
	t.Run("", func(t *testing.T) {
		stat, _ := GetRelationStat(_test_uid)
		times := int(math.Ceil(float64(stat.Data.Following) / 50))
		log.Println(times)
		for i := 0; i < times; i++ {
			ret, err := GetRelationFollowings(_test_uid, i+1, "")
			if err != nil {
				panic(err)
			}
			for j, user := range ret.Data.List {
				fmt.Println(i*50+j+1, user.Mid, user.Uname)
			}
		}
	})
}

func TestGetRelationStat(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ret, err := GetRelationStat(_test_uid)
		if err == nil {
			log.Println(ret.Data.Following)
		} else {
			panic(err)
		}
	})
}
