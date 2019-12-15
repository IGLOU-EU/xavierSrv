package main

import (
    _"os"
    "fmt"
    "time"
    "sync"
    "errors"
    "strconv"
    "net/http"
    "math/rand"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Overall struct {
		LoopWait      int    `toml:"loopWait"`
		CheckListFile string `toml:"checkListFile"`
		ErrorListFile string `toml:"errorListFile"`
	} `toml:"overall"`
	HTTP struct {
		MaxWait int `toml:"maxWait"`
	} `toml:"http"`
}

type CheckList struct {
	Team []struct {
		Enable bool     `toml:"enable"`
		Name   string   `toml:"name"`
		Error  []string `toml:"error"`
		App    []struct {
			Name     string `toml:"name"`
            URL      string `toml:"url"`
			Response int    `toml:"response"`
		} `toml:"app"`
    } `toml:"team"`
}

type ErrorList struct {
	Team map[string]map[string]struct {
		Process string `toml:"process"`
		SMTP    struct {
			TLS      bool   `toml:"tls"`
			Starttls bool   `toml:"starttls"`
			Port     string `toml:"port"`
			Host     string `toml:"host"`
			From     string `toml:"from"`
			User     string `toml:"user"`
			Passwd   string `toml:"passwd"`
		} `toml:"smtp"`
		Message struct {
			Recipients []string `toml:"recipients"`
			Subject    string   `toml:"subject"`
			Signing    string   `toml:"signing"`
			Body       string   `toml:"body"`
		} `toml:"message"`
		HTTP    struct {
			Methods     string `toml:"methods"`
			URL         string `toml:"url"`
			Body        string `toml:"body"`
			ContentType string `toml:"Content-Type"`
		} `toml:"http"`
	} `toml:"team"`
}

var config    Config
var checkList CheckList
var errorList ErrorList
var wg sync.WaitGroup

func main() {
//init:
    rand.Seed(42)

    if _, err := toml.DecodeFile("./examples/etc/xavier-srv/config.toml", &config); err != nil {
        panic(err)
	}
    if _, err := toml.DecodeFile(config.Overall.CheckListFile, &checkList); err != nil {
        panic(err)
	}
    if _, err := toml.DecodeFile(config.Overall.ErrorListFile, &errorList); err != nil {
        panic(err)
	}

loop:
    for index, team := range checkList.Team {
        if team.Enable {
            wg.Add(1)
            go xTeam(index)
        }
    }

    wg.Wait()
    time.Sleep(time.Duration(config.Overall.LoopWait))
    goto loop
}

func pushLog(lvlId int, message error) () {
    lvlName := "??"

    switch lvlId {
        case 6:
            lvlName = "Info"
        case 3:
            lvlName = "Error"
        case 0:
            lvlName = "Panic"
	}

    fmt.Printf("[%s] (%s) %s\n", time.Now().Format(time.RFC3339), lvlName, message);
}

func xTeam(teamId int) (bool) {
    team := &checkList.Team[teamId]

    for _, app := range team.App {
        pushLog(6, errors.New("Checking '" + app.Name + "' ..."))

        if _, err := cerebro(app.URL, app.Response); err != nil {
            pushLog(3, err)
        }
    }

    defer wg.Done()
    return true
}

func httpStatus(url string, need int) (bool, error) {
    time.Sleep(time.Millisecond * time.Duration(rand.Intn(config.HTTP.MaxWait)))

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

func cerebro(url string, need int) (bool, error) {
    xPsyAgatha, errAgatha := httpStatus(url, need)
    xPsyArthur, errArthur := httpStatus(url, need)
    xPsyDash  , errDash   := httpStatus(url, need)

    if (xPsyAgatha && xPsyArthur) ||
       (xPsyAgatha && xPsyDash)   ||
       (xPsyArthur && xPsyDash)   {
            return true, nil
    }

    err := errors.New("Inconsistent error returned")

    if errAgatha.Error() == errArthur.Error() {
        err = errAgatha
    } else if errArthur.Error() == errDash.Error() {
        err = errArthur
    } else if errDash.Error() == errAgatha.Error() {
        err = errDash
    }

    return false, err
}
