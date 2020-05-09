package tools

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.iglou.eu/xavierSrv/config"
)

func JsonEscapeString(d string) string {
	r := strings.NewReplacer(
		"\\", "\\\\",
		"/", "\\/",
		"\"", "\\\"",
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\x08", "\\f",
		"\x0c", "\\b",
	)
	return r.Replace(d)
}

func HttpStatus(url string, need int) (bool, error) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(config.Global.HTTP.MaxWait)))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	resp, err := client.Head(url)

	if err == nil {
		statusCode := (need == resp.StatusCode)

		if statusCode {
			return true, nil
		}

		return false, errors.New("Head " + url + ": Response " + strconv.Itoa(resp.StatusCode) + ": Expected status code " + strconv.Itoa(need))
	}

	return false, err
}
