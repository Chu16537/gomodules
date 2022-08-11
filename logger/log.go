package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func Log(errorCode int32, v interface{}) {

}

func Debug(s string, v ...interface{}) {
	color.Cyan("{'level':'Debug','time': '%v','msg':'%+v'}\n", time.Now().UTC().Format("2006/01/02 15:04:05"), fmt.Sprintf(s, v...))
}

func Info(s string, v ...interface{}) {
	color.Green("{'level':'Info','time': '%v','msg':'%+v'}\n", time.Now().UTC().Format("2006/01/02 15:04:05"), fmt.Sprintf(s, v...))
}

func Warn(s string, v ...interface{}) {
	color.Yellow("{'level':'Warn','time': '%v','msg':'%+v'}\n", time.Now().UTC().Format("2006/01/02 15:04:05"), fmt.Sprintf(s, v...))
}

func Error(s string, v ...interface{}) {
	color.Red("{'level':'Error','time': '%v','msg':'%+v'}\n", time.Now().UTC().Format("2006/01/02 15:04:05"), fmt.Sprintf(s, v...))
}
