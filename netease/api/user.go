package api

type User struct {
	Code    int64         `json:"code"`
	Token   string        `json:"token"`
	Account accountStruct `json:"account"`
	Profile profileStruct `json:"profile"`
}

type accountStruct struct {
	Salt           string `json:"salt"`
	UserName       string `json:"userName"`
	ViptypeVersion int64  `json:"viptypeVersion"`
}
type profileStruct struct {
	UserId        int64  `json:"userId"`
	VipType       int64  `json:"vipType"`
	Nickname      string `json:"nickname"`
	AtarUrl       string `json:"atarUrl"`
	BackgroundUrl string `json:"backgroundUrl"`
	Signature     string `json:"signature"`
}
