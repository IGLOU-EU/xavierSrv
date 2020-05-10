package main

import (
	"errors"
	"math/rand"
	"strings"
	"sync"
	"time"

	"git.iglou.eu/xavierSrv/config"
	"git.iglou.eu/xavierSrv/tools"
)

var wg sync.WaitGroup

func main() {
	//init:
	rand.Seed(42)
	config.Global.Init(&config.Check, &config.Error)
	tools.InitLogConf(
		config.Global.Overall.LogFile,
		config.Global.Overall.LogVerbose,
	)

	//loop:
	for index, team := range config.Check.Team {
		if team.Enable {
			wg.Add(1)
			go xTeam(index)
		}
	}

	wg.Wait()
	time.Sleep(time.Duration(config.Global.Overall.LoopWait))
	//goto loop
}

func xTeam(teamID int) bool {
	team := &config.Check.Team[teamID]
	var failures [][]string

	for _, app := range team.App {
		tools.PushToLog(6, errors.New("Checking '"+app.Name+"' ..."))

		if _, err := cerebro(app.URL, app.Response); err != nil {
			returned := []string{app.Name, err.Error()}
			failures = append(failures, returned)

			tools.PushToLog(6, err)
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

	if config.Global.Reporting.Local {
		//report.LocalSave(team, report.LocalBuild(failures))
	}

	// Build message
	reportMessage := ""
	for i := range failures {
		reportMessage += "[" + failures[i][0] + "]" + "\n" +
			failures[i][1] + "\n"
	}

	// for reporting
	for _, reportProcess := range config.Error.Team[team].Report {
		switch reportProcess.Process {
		case "smtp":
			subject := strings.ReplaceAll(reportProcess.Subject, "[%TEAM]", team)
			message := strings.ReplaceAll(reportProcess.Body, "[%TEAM]", team)
			message = strings.ReplaceAll(message, "[%ERRORS]", reportMessage)

			tools.SmtpReport(
				reportProcess.Host,
				reportProcess.Encrypt,
				reportProcess.From,
				reportProcess.User,
				reportProcess.Passwd,
				reportProcess.Recipients,
				subject, message,
			)
		case "http":
			message := strings.ReplaceAll(reportProcess.Body, "[%TEAM]", team)

			if reportProcess.Methods == "POST" {
				message = strings.ReplaceAll(message, "[%ERRORS]", tools.JsonEscapeString(reportMessage))
			} else {
				message = strings.ReplaceAll(message, "[%ERRORS]", reportMessage)
			}

			tools.HttpReport(reportProcess.Methods, reportProcess.URL, message)
		default:
			tools.PushToLog(3, errors.New("Unknow '"+reportProcess.Process+"' reporting process"))
		}
	}
	// exec type

	/*if config.Global.Reporting.ReportJSON {
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

func cerebro(url string, need int) (bool, error) {
	xPsyAgatha, errAgatha := tools.HttpStatus(url, need, config.Global.HTTP.MaxWait)
	xPsyArthur, errArthur := tools.HttpStatus(url, need, config.Global.HTTP.MaxWait)
	xPsyDash, errDash := tools.HttpStatus(url, need, config.Global.HTTP.MaxWait)

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
