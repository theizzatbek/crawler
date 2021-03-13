package controllers

import (
	"crawler/utils"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
)

// Get titles from urls
// Example body:
// {
//    "urls" : [
//        "http://google.com",
//        "yandex.ru",
//        "http://vk.com",
//        "http://daryo.uz",
//        "tes",
//        "asd"
//    ]
// }
// out:
// {
//    "data": [
//        "Google",
//        "",
//        "VKontaktening mobil versiyasi | VKontakte",
//        "Daryo — yangiliklar daryosidan chetda qolib ketmang!",
//        "",
//        ""
//    ],
//    "ok": true
// }
func GetTitle(c echo.Context) error {
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, utils.Error(err.Error()))
	}

	var msg = utils.Message()
	urls := m["urls"].([]interface{})
	var length = len(urls)
	var title = make([]string, length, length)
	var wg = sync.WaitGroup{}

	wg.Add(length)
	for i := range urls {

		go func(i int) {
			var url = urls[i].(string)
			page, err := parse(url)
			if err != nil {
				log.Errorf("cannot get page info url:%s err:%s\n", url, err)
			} else {
				title[i] = getTitle(&page)
			}
			wg.Done()
		}(i)

	}
	wg.Wait()
	msg["data"] = title

	return c.JSON(http.StatusOK, msg)
}

func parse(url string) (string, error) {

	r, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("cannot get page: %v", err)
	}
	defer func() {
		_ = r.Body.Close()
	}()

	b, err := ioutil.ReadAll(r.Body)

	return string(b), err
}

func getTitle(h *string) string {
	r := strings.NewReader(*h)
	tokenizer := html.NewTokenizer(r)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			return ""
		}
		if tt == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "title" {
				tokenizer.Next()
				return tokenizer.Token().Data
			}
		}

	}
}
