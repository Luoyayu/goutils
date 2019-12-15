package d4j

import myRedis "github.com/luoyayu/goutils/redis"

type Struct struct {
	ID        string
	Name      string
	Author    string
	ShareLink *share
	Tags      []string
	KeyWords  []string
	Category  string
	Title     string // not for tu_shu_zi_yuan
}

type share struct {
	Url string
	Key string
}

type Redis struct {
	*myRedis.Client
}

type task struct {
	BookId int64
}
