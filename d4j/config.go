package d4j

type SanqiuConfig struct {
	Sanqiu      *Sanqiu            `toml:"sanqiu"`
	SanqiuRedis *ConfigRedisStruct `toml:"sanqiu_redis"`
}

type Sanqiu struct {
	Debug bool `toml:"debug"`
	targets
}
type targets struct {
	BaseUrl      string `toml:"base_url"`
	RssPath      string `toml:"rss_path"`
	InfoPath     string `toml:"info_path"`
	DownloadPath string `toml:"download_path"`
}

type ConfigRedisStruct struct {
	Addr   string `toml:"addr"`
	Enable bool   `toml:"enable"`
	Db     int    `toml:"db"`
}
