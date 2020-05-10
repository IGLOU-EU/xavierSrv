package config

import (
	"fmt"
	"os"

	"git.iglou.eu/xavierSrv/tools"
	"github.com/BurntSushi/toml"
)

// Global : External global var
var Global Config

// Check : External global var
var Check CheckList

// Error : External global var
var Error ErrorList

// Init : for init config and log
func (co *Config) Init(ch *CheckList, er *ErrorList) {
	if _, err := toml.DecodeFile("./examples/etc/xavier-srv/config.toml", &co); err != nil {
		panic(err)
	}
	if _, err := toml.DecodeFile(co.Overall.CheckListFile, &ch); err != nil {
		panic(err)
	}
	if _, err := toml.DecodeFile(co.Overall.ErrorListFile, &er); err != nil {
		panic(err)
	}
	if f, err := os.Create(co.Overall.LogFile); err != nil {
		panic(err)
	} else {
		f.Close()
	}

	tools.PushToLog(-1, fmt.Errorf(
		"=>\n--- [Xavier] ---\nLog file: %v\nApp list file: %v\nErr list file: %v\n---",
		co.Overall.LogFile,
		co.Overall.CheckListFile,
		co.Overall.ErrorListFile,
	))
}
