package d4j

import (
	"github.com/BurntSushi/toml"
	"github.com/anaskhan96/soup"
	"github.com/luoyayu/goutils/net"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Config *SanqiuConfig
var DEBUG = false

func init() {
	configPath := "d4j_config.toml"
	Config = &SanqiuConfig{}
	if _, err := toml.DecodeFile(configPath, &Config); err != nil {
		log.Println("not found d4j_config.toml, loading config from main thread!")
	} else {
		DEBUG = Config.Sanqiu.Debug
	}
}

// split `bookName【authorName】`
func splitTitle(title string, book *Struct) {
	nameWAuthorRunes := []rune(title)

	var tokenLeft, tokenRight int

	for i := len(nameWAuthorRunes) - 1; i >= 0; i-- {
		if nameWAuthorRunes[i] == '【' {
			tokenLeft = i
			book.Name = strings.TrimSpace(string(nameWAuthorRunes[:tokenLeft]))
			break
		}
	}

	for i := len(nameWAuthorRunes) - 1; i >= 0; i-- {
		if nameWAuthorRunes[i] == '】' {
			tokenRight = i
			book.Author = strings.TrimSpace(string(nameWAuthorRunes[tokenLeft+1 : tokenRight]))
			break
		}
	}

	if book.Name != "" || book.Author != "" {
		return

	}

	for i := len(nameWAuthorRunes) - 1; i >= 0; i-- {
		if nameWAuthorRunes[i] == '（' {
			tokenLeft = i
			book.Name = strings.TrimSpace(string(nameWAuthorRunes[:tokenLeft]))
			break
		}
	}

	for i := len(nameWAuthorRunes) - 1; i >= 0; i-- {
		if nameWAuthorRunes[i] == '）' {
			tokenRight = i
			book.Author = strings.TrimSpace(string(nameWAuthorRunes[tokenLeft+1 : tokenRight]))
			book.Author = strings.TrimLeft(book.Author, "作者-")
			book.Author = strings.TrimLeft(book.Author, "作者:")
			book.Author = strings.TrimLeft(book.Author, "作者：")
			break
		}
	}

}

func getCategoryFromKeyWords(kw []string) string {
	kws := strings.Join(kw, ",")
	if strings.Contains(kws, "动漫资讯") {
		return "dong_man_zi_xun"
	} else if strings.Contains(kws, "影评影讯") {
		return "ying_ping_ying_xun"
	} else if strings.Contains(kws, "文化漫谈") {
		return "wen_hua_man_tan"
	} else {
		return "tu_shu_zi_yuan"
	}
}

func GetBookInfoFromRss(justGetMaxId bool) (books []*Struct, maxId string, err error) {
	defer func() {
		err = errors.Wrap(err, "GetBookInfoFromRss ->")
	}()

	var feed *gofeed.Feed

	rssUrl, _ := net.JoinPath(Config.Sanqiu.BaseUrl, []interface{}{Config.Sanqiu.RssPath}) // rss page
	feed, err = gofeed.NewParser().ParseURL(rssUrl.String())
	if err != nil {
		return
	}
	books = make([]*Struct, len(feed.Items))

	for i, item := range feed.Items {
		linksSlice := strings.Split(item.Link, "/")
		bookPath := linksSlice[len(linksSlice)-1]
		bookId := strings.Split(bookPath, ".")[0]
		_, books[i], _ = GetBookInfo(bookId)
		books[i].ID = bookId
		if maxId == "" && justGetMaxId == true {
			maxId = books[i].ID
			return nil, maxId, nil
		}
		splitTitle(item.Title, books[i])
		time.Sleep(50 * time.Microsecond)
	}
	return books, "", err
}

func GetBookShareFromDLPage(bookId string, book *Struct) (err error) {
	var qUrl *url.URL
	var resp string

	params := url.Values{}
	qUrl, _ = net.JoinPath(Config.Sanqiu.BaseUrl, []interface{}{Config.Sanqiu.DownloadPath}) // book download page
	params.Add("id", book.ID)

	if resp, err = soup.Get(qUrl.String() + "?" + params.Encode()); err == nil {
		doc := soup.HTMLParse(resp)
		keyDiv := doc.
			Find("div", "class", "plus_l").
			FindAll("li")

		for _, li := range keyDiv {
			if li.Find("font").Error == nil && strings.Contains(li.Text(), "百度") {
				book.ShareLink.Key = li.Find("font").Text()
			} else if strings.Contains(li.Text(), "名称") {
				a := strings.Split(li.Text(), "：")
				if len(a) <= 1 {
					continue
				}
				if a[1] == "" {
					continue
				}
				book.Name = strings.TrimSpace(a[1])
			} else if strings.Contains(li.Text(), "作者") {
				a := strings.Split(li.Text(), "：")
				if len(a) <= 1 {
					continue
				}
				if a[1] == "" {
					continue
				}
				authorRunes := []rune(strings.TrimSpace(a[1]))
				if authorRunes[0] == '【' && authorRunes[len(authorRunes)-1] == '】' {
					authorRunes = authorRunes[1 : len(authorRunes)-1]
				}
				book.Author = string(authorRunes)
			}
		}

		linkDiv := doc.
			Find("div", "class", "panel").
			Find("div", "class", "panel-body").
			FindAll("a")
		for _, link := range linkDiv {
			if strings.Contains(link.Text(), "百度") {
				book.ShareLink.Url = link.Attrs()["href"]
				break
			}
		}
	}
	return
}

func GetBookInfo(bookId string) (int64, *Struct, error) {
	if bookId == "" {
		panic("book ID is None!")
	}
	var err error

	book := &Struct{ID: bookId, ShareLink: &share{}}

	httpC := &http.Client{}
	qurl, _ := net.JoinPath(Config.Sanqiu.BaseUrl, []interface{}{bookId + ".html"}) // book main page
	req, _ := http.NewRequest("GET", qurl.String(), nil)

	if resp, err := httpC.Do(req); err == nil {
		if resp.StatusCode == 404 { // 404
			return 0, nil, errors.New("not found: " + qurl.String())
		} else if strings.TrimRight(qurl.String(), "/") !=
			strings.TrimRight(resp.Request.URL.String(), "/") { // 301
			// FIXME: ignore 301 response's header location
			return 0, nil, errors.New("not found: " + qurl.String())
			rePath := resp.Request.URL.Path[1:] // trim left `/`
			bookIdWithExt := strings.Split(rePath, "/")[0]
			if nextBookId, err := strconv.ParseInt(strings.Split(bookIdWithExt, ".")[0], 10, 64); err != nil { // error when parse new bookid
				log.Println("oUrl:", qurl.String())
				log.Println("reUrl:", resp.Request.URL)
				log.Println("rePath:", rePath)
				return 0, nil, errors.New("redirect to " + rePath + " error!")
			} else {
				return nextBookId, nil, nil
			}
		} else { // 200
			defer func() {
				v := recover()
				if v != nil {
					log.Println("recover: ", qurl)
					panic(v)
				}
			}()

			b, _ := ioutil.ReadAll(resp.Body)
			mainPageDoc := soup.HTMLParse(string(b))
			head := mainPageDoc.Find("head")

			title, ogTitle, keywords := "", "", ""

			if title_ := head.Find("title"); title_.Error == nil {
				title = title_.Text()
			} else {
				return 0, nil, errors.New("hit 404.html")
			}
			if ogTitle_ := head.Find("meta", "property", "og:title"); ogTitle_.Error == nil {
				ogTitle = ogTitle_.Attrs()["content"]
			} else {
				return 0, nil, errors.New("hit 404.html")
			}
			if keywords_ := head.Find("meta", "name", "keywords"); keywords_.Error == nil {
				keywords = keywords_.Attrs()["content"]
			} else {
				return 0, nil, errors.New("hit 404.html")
			}

			// get title
			splitTitle(title, book)
			if book.Title == "" {
				book.Title = title
				book.Name = title
			}

			// get keyword
			keywordsList := strings.Split(keywords, ",")
			if len(keywordsList) > 1 {
				book.KeyWords = keywordsList[1 : len(keywordsList)-1]
			}

			// get Category
			book.Category = getCategoryFromKeyWords(book.KeyWords)

			if book.Category != "tu_shu_zi_yuan" {
				book.Title = title
				reg, _ := regexp.Compile("《([^》]+)》")
				result := reg.FindAll([]byte(ogTitle), -1)
				if len(result) == 1 {
					resultRunes := []rune(string(result[0]))
					book.Name = string(resultRunes[1 : len(resultRunes)-1])
				} else if len(result) > 1 {
					book.Name = ""
					for i, re := range result {
						rr := []rune(string(re))
						tmp := string(rr[1 : len(rr)-1])
						book.Name += tmp
						if i != len(result)-1 {
							book.Name += "-"
						}
					}

				}
			}

			// get tags
			tagsDiv := mainPageDoc.
				Find("footer").
				Find("div", "class", "footer-tag").
				Find("div", "class", "pull-left").
				FindAll("a")
			for _, tag := range tagsDiv {
				book.Tags = append(book.Tags, tag.Text())
			}

			// get share link
			// 1) for old page

			postContent := mainPageDoc.Find("div", "class", "kratos-post-content")
			postContentLinks := postContent.FindAll("a")
			ps := postContent.FindAll("p")

			for _, p := range ps {
				if len(p.FindAll("a")) != 0 && p.FindAll("a")[0].Attrs()["class"] == "downbtn downcloud" {
					tmp := strings.Split(p.Text(), ":")
					book.ShareLink.Key = strings.TrimSpace(tmp[len(tmp)-1])
				}
			}

			for _, a := range postContentLinks {
				if a.Attrs()["class"] == "downbtn downcloud" {
					book.ShareLink.Url = a.Attrs()["href"]
				}
			}

			// 2) for new page
			if strings.Contains(book.ShareLink.Url, "pan") == false {
				_ = GetBookShareFromDLPage(book.ID, book)
			}
		}
	}
	return 0, book, err
}
