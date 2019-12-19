package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

// Config is for
type Config struct {
	Overall struct {
		LoopWait      int    `toml:"loopWait"`
		CheckListFile string `toml:"checkListFile"`
		ErrorListFile string `toml:"errorListFile"`
	} `toml:"overall"`
	HTTP struct {
		MaxWait int `toml:"maxWait"`
	} `toml:"http"`
	Reporting struct {
		Local bool   `toml:"local"`
		Dir   string `toml:"dir"`
	} `toml:"reporting"`
}

// CheckList is for
type CheckList struct {
	Team []struct {
		Enable bool   `toml:"enable"`
		Name   string `toml:"name"`
		App    []struct {
			Name     string `toml:"name"`
			URL      string `toml:"url"`
			Response int    `toml:"response"`
		} `toml:"app"`
	} `toml:"team"`
}

// ErrorList  is for
type ErrorList struct {
	Team map[string]struct {
		Report []struct {
			Process     string   `toml:"process"`
			Host        string   `toml:"host,omitempty"`
			From        string   `toml:"from,omitempty"`
			User        string   `toml:"user,omitempty"`
			Passwd      string   `toml:"passwd,omitempty"`
			Recipients  []string `toml:"recipients,omitempty"`
			Subject     string   `toml:"subject,omitempty"`
			Body        string   `toml:"body"`
			Methods     string   `toml:"methods,omitempty"`
			URL         string   `toml:"url,omitempty"`
			ContentType string   `toml:"Content-Type,omitempty"`
		} `toml:"report"`
	} `toml:"team"`
}

var config Config
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

func pushLog(lvlID int, message error) {
	lvlName := "??"

	switch lvlID {
	case 6:
		lvlName = "Info"
	case 3:
		lvlName = "Error"
	case 0:
		lvlName = "Panic"
	}

	fmt.Printf("[%s] (%s) %s\n", time.Now().Format(time.RFC3339), lvlName, message)
}

func xTeam(teamID int) bool {
	team := &checkList.Team[teamID]
	var failures [][]string

	for _, app := range team.App {
		pushLog(6, errors.New("Checking '"+app.Name+"' ..."))

		if _, err := cerebro(app.URL, app.Response); err != nil {
			returned := []string{app.Name, err.Error()}
			failures = append(failures, returned)

			pushLog(3, err)
		}
	}

	if failures != nil {
		xTeamReport(team.Name, failures)
	}

	defer wg.Done()
	return true
}

func xTeamReport(team string, failures [][]string) (bool, error) {
	//var message string {"name":"","error":""}

	if config.Reporting.Local {
		//report.LocalSave(team, report.LocalBuild(failures))
	}

	// Build message
	reportMessage := ""
	for i := range failures {
		reportMessage += "[" + failures[i][0] + "]" + "\n" +
			failures[i][1] + "\n\n"
	}

	// for reporting
	for _, reportProcess := range errorList.Team[team].Report {
		switch reportProcess.Process {
		case "smtp":
			subject := strings.ReplaceAll(reportProcess.Subject, "[%TEAM]", team)
			message := strings.ReplaceAll(reportProcess.Body, "[%TEAM]", team)
			message = strings.ReplaceAll(reportProcess.Body, "[%ERRORS]", team)
			_ = subject + message
			//report.byMail(host, from, user, passwd, subject, message string, recip []string)
		case "http":
			message := strings.ReplaceAll(reportProcess.Body, "[%TEAM]", team)
			message = strings.ReplaceAll(reportProcess.Body, "[%ERRORS]", team)
			_ = message
			//report.byHttp(methods, url, message string)
		default:
			pushLog(3, errors.New("Unknow '"+reportProcess.Process+"' reporting process"))
		}
	}
	// exec type

	/*if config.Reporting.ReportJSON {
		jsonFile := "{\"status\":\"error\",\"date\":\""
		jsonFile += time.Now().Format(time.RFC3339)
		jsonFile += "\",\"listing\":["

		i, l := 0, len(failures)
		for {
			if i%2 == 0 {
				jsonFile += "{\"name\":\"" + failures[i] + "\""
			} else {
				jsonFile += ",\"error\":\"" + failures[i] + "\"},"
			}

			i++
			if i == l {
				jsonFile = jsonFile[:len(jsonFile)-1]
				break
			}
		}

		jsonFile += "]}"
		fmt.Println(jsonFile)
	}*/

	/*
		Si on a une erreur
		- on construit le json de la team
		- on formate le message de sortie

		On boucle sur les reports
		- pour les reports smtp on envoi le message
		- pour les retapi on export les info format json puis send

		Ne pas oublier de replace les balises [%***]
	*/

	return true, nil
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
		}

		return false, errors.New("Head " + url + ": Response " + strconv.Itoa(resp.StatusCode) + ": Expected status code " + strconv.Itoa(need))
	}

	return false, err
}

func cerebro(url string, need int) (bool, error) {
	xPsyAgatha, errAgatha := httpStatus(url, need)
	xPsyArthur, errArthur := httpStatus(url, need)
	xPsyDash, errDash := httpStatus(url, need)

	if (xPsyAgatha && xPsyArthur) ||
		(xPsyAgatha && xPsyDash) ||
		(xPsyArthur && xPsyDash) {
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
