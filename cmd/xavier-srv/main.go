package main

import (
    "fmt"
    "errors"
    "strconv"
    "net/http"
)

func main() {
    fmt.Println(cerebro("https://test.iglou.eu", 200))
    fmt.Println(cerebro("https://iglou.eu", 200))
    fmt.Println(cerebro("https://iglou.eu", 301))
    fmt.Println(cerebro("https://git.iglou.eu", 200))
    fmt.Println(cerebro("https://media.iglou.eu", 200))
    fmt.Println(cerebro("https://51.15.11.133", 200))
}

func httpStatus(url string, need int) (bool, error) {
    client := &http.Client{
	    CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
    }}

    resp, err := client.Head(url)

    if err == nil {
        statusCode := (need == resp.StatusCode)

        if statusCode {
            return true, nil
        } else {
            return false, errors.New("Head " + url + ": Response " + strconv.Itoa(resp.StatusCode) + ": Expected status code " + strconv.Itoa(need))
        }
    }

    return false, err
}

func cerebro(url string, need int) (bool) {
    xPsyAgatha, err := httpStatus(url, need)
    xPsyArthur, _ := httpStatus(url, need)
    xPsyDash  , _ := httpStatus(url, need)

    if (xPsyAgatha && xPsyArthur) ||
       (xPsyAgatha && xPsyDash)   ||
       (xPsyArthur && xPsyDash)   {
        return true
    }

    fmt.Println(err)
    return false
}
