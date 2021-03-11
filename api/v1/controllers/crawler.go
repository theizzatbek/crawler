package controllers

import (
	"crawler/utils"
	"fmt"
	"github.com/labstack/echo"
	"golang.org/x/net/html"
	"net/http"
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
	var title = make([]string, len(urls), len(urls))
	for i := range urls {
		var url = urls[i].(string)
		page, err := parse(url)
		if err != nil {
			fmt.Printf("Error getting page %s %s\n", url, err)
			continue
		}
		title[i] = getTitle(page)
	}
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
		return nil, fmt.Errorf("cannot get page")
	}
	b, err := html.Parse(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot parse page")
	}
	return b, err
}
