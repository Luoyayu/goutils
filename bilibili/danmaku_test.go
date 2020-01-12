package biliAPI

import (
	"context"
	"testing"
)

func TestDanmakuClient_NewDanmakuClient(t *testing.T) {
	var c = &DanmakuClient{}
	c.NewDanmakuClient(context.Background(), 21452505)
}
