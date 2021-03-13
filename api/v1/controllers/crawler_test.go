package controllers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	body = `{
    "urls" : [
        "http://google.com",
        "yandex.ru",
        "http://vk.com",
        "http://daryo.uz",
        "tes",
        "asd"
    ]
}`
	getTitleJson = `{
    "data": [
        "Google",
        "",
        "VKontaktening mobil versiyasi | VKontakte",
        "Daryo — yangiliklar daryosidan chetda qolib ketmang!",
        "",
        ""
    ],
    "ok": true
}`
)

func TestGetTitle(t *testing.T) {

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/get/title")

	if assert.NoError(t, GetTitle(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var expected, actual echo.Map
		if assert.NoError(t, json.Unmarshal([]byte(getTitleJson), &expected)) &&
			assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &actual)) {
			assert.Equal(t, expected, actual)
		}

	}

}
