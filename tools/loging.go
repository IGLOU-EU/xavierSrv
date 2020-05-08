package tools

import (
	"fmt"
	"time"
)

var logFile string
var logVerbose bool

// InitLogConf : For init local logFile var
func InitLogConf(f string, v bool) {
	logFile = f
	logVerbose = v
}

// PushToLog : Log manager
func PushToLog(lvlID int, message error) {
	if !logVerbose && lvlID == 6 {
		return
	}

	msgout := ""
	lvlName := "??"

	switch lvlID {
	case -1:
		lvlName = "Start"
	case 6:
		lvlName = "Info"
	case 3:
		fallthrough
	default:
		lvlName = "Error"
	}

	msgout = fmt.Sprintf("[%s] (%s) %s", time.Now().Format(time.RFC3339), lvlName, message)

	//logToFile(msgout)
	fmt.Println(msgout)
}

/*func logToFile(m string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0644)

	//err != nil {
	//	panic(err)
	//}

	if _, err := fmt.Fprintln(f, m); err != nil {
		f.Close()
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}*/
