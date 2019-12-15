package gadio

import (
	"github.com/luoyayu/goutils/gcore"
	myRedis "github.com/luoyayu/goutils/redis"
)

// ********************** Gadio *********************** //
type RadiosQuery struct {
	Data []RadioData `json:"data"`
}

type RadioData struct {
	ID    string      `json:"id"`
	Attrs *Attributes `json:"attributes"`
}

type Attributes struct {
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Cover       string `json:"cover"`
	PublishedAt string `json:"published-at"`
}

// ********************** Toml Config *********************** //
type GadioConfig struct {
	Gcore      *gcore.Struct      `toml:"gcore"`
	Gadio      *ConfigGadioStruct `toml:"gcore_gadio"`
	GadioRedis *ConfigRedisStruct `toml:"gcore_redis"`
}

type ConfigGadioStruct struct {
	ImageApi   string `toml:"image_api"`
	UploadsApi string `toml:"uploads_api"`
}

type ConfigRedisStruct struct {
	Addr   string `toml:"addr"`
	Enable bool   `toml:"enable"`
	Db     int    `toml:"db"`
}

// ********************** GadioRedis *********************** //
type Redis struct {
	*myRedis.Client
}
