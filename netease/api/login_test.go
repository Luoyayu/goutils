package api

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"testing"
)

func TestNetEaseClient_LoginByCellphone(t *testing.T) {
	_ = afero.NewOsFs().Remove("cookies")
	if err := _client.LoginByCellphone(os.Getenv("PHONE"), os.Getenv("PASSWD")); err != nil {
		log.Println(err)
	}
}
