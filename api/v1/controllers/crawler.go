package controllers

import (
	"crawler/utils"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/html"
	"log"
	"net/http"
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

	for i := range urls {
		wg.Add(1)
		go func(i int) {
			var url = urls[i].(string)
			page, err := parse(url)
			if err != nil {
				log.Printf("cannot get page info url:%s err:%s\n", url, err)
			} else {
				title[i] = getTitle(page)
			}
			wg.Done()
		}(i)

	}
	wg.Wait()
	msg["data"] = title

	return c.JSON(http.StatusOK, msg)
}

func getTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = getTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

func parse(url string) (*html.Node, error) {

	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get page: %v", err)
	}
	defer r.Body.Close()
	b, err := html.Parse(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot parse page")
	}
	return b, err
}
