package tools

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
)

func HttpReport(m string, u string, d string) {
	var res *http.Response
	var err error

	if m == "POST" {
		res, err = http.Post(u, "application/json", bytes.NewBufferString(d))
	} else {
		res, err = http.Get(u + url.QueryEscape(d))
	}

	if err != nil {
		PushToLog(3, err)
	}

	if res.StatusCode != 200 {
		PushToLog(3, errors.New("Request to `"+u+url.QueryEscape(d)+"` has failed: "+res.Status))
	}
}
