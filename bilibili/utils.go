package biliAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func GetDefault(tUrl string, params map[string]interface{}, in interface{}) (out interface{}, err error) {
	out = in
	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	var resp *http.Response
	req, _ := http.NewRequest("GET", tUrl+"?"+l.Encode(), nil)
	if resp, err = (&http.Client{}).Do(req); err == nil {
		defer resp.Body.Close()
		//b, _ := ioutil.ReadAll(resp.Body)
		//log.Println(string(b))
		err = json.NewDecoder(resp.Body).Decode(&in)
		if err != nil {
			return
		}
	}
	return
}

func GetWithReq(tUrl string, params map[string]interface{}, req *http.Request, in interface{}) (out interface{}, err error) {
	out = in
	l := url.Values{}
	for k, v := range params {
		l.Add(k, fmt.Sprint(v))
	}

	req.URL, _ = url.Parse(tUrl + "?" + l.Encode())

	var resp *http.Response
	if resp, err = (&http.Client{}).Do(req); err == nil {
		defer resp.Body.Close()

		//b, _ := ioutil.ReadAll(resp.Body)
		//log.Println(string(b))
		//err = json.NewDecoder(bytes.NewReader(b)).Decode(&in)

		err = json.NewDecoder(resp.Body).Decode(&in)
		if err != nil {
			return
		}
	}
	return
}
