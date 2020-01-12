package biliAPI

type respStruct struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		TokenInfo  *TokenInfoStruct  `json:"token_info"`
		CookieInfo *CookieInfoStruct `json:"cookie_info"`
		SSO        []string          `json:"sso"`
	} `json:"data"`
}

type TokenInfoStruct struct {
	Hash        string `json:"hash"`
	Key         string `json:"key"`
	AccessToken string `json:"access_token"`
	Mid         int64  `json:"mid"`
	ExpiresIn   int64  `json:"expires_in"` // 2592000s
}

type CookieInfoStruct struct {
	Cookies []*struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
		HttpOnly int    `json:"http_only"`
		Expires  int64  `json:"expires"`
	} `json:"cookies"`
	CookiesMap map[string]string
	Domains    []string `json:"domains"`
}
