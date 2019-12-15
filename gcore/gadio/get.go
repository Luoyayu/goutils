package gadio

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/url"
)

var Config *GadioConfig
var DEBUG = false

func init() {
	tomlPath := "gcore_config.toml"
	Config := &GadioConfig{}
	if _, err := toml.DecodeFile(tomlPath, &Config); err != nil {
		log.Println("not found gcore_config.toml, loading config from main thread!")
	} else {
		DEBUG = Config.Gcore.Debug
		if DEBUG {
			log.Printf("Gadio Config: %+v\n", Config)
		}
	}
}

func GetLatestN(radiosNum int) (resultMap *RadiosQuery, err error) {
	defer func() {
		err = errors.Wrap(err, "GetGadioLatestN ->")
	}()
	params := url.Values{}
	params.Add("page[limit]", fmt.Sprint(radiosNum))
	params.Add("page[offset]", "0")
	params.Add("sort", "-published-at")
	params.Add("include", "")
	params.Add("fields[radios]", "title,desc,cover,published-at,")
	params.Add("from-app", "1")

	radioUrl, _ := url.Parse(Config.Gcore.GApi + "radios?" + params.Encode())

	var resp *http.Response
	if resp, err = http.Get(radioUrl.String()); err == nil {
		if err = json.NewDecoder(resp.Body).Decode(&resultMap); err == nil {
			if DEBUG {
				log.Println(radioUrl)
				log.Println("decode body of response ok!")
			}
			return
		}
	}
	return nil, err
}
