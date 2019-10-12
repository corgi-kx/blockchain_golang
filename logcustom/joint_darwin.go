package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//传入要打印的信息，拼接好后返回
func joint(prefix, message string, color colorType) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	filename, funcname, line := getpProcInfo()
	s := fmt.Sprint(prefix, ": ", now, " ", filename, ":", line, ":", funcname, ": ", message)
	if isColor {
		//s = fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
		s = fmt.Sprintf("\033[%dm%s\033[0m", color, s)
	}
	return s
}

func (l mylog) joint(prefix, message string, color colorType) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	filename, funcname, line := getpProcInfo()
	s := fmt.Sprint(prefix, ": ", now, " ", filename, ":", line, ":", funcname, ": ", message)
	if l.isColor {
		//s = fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
		s = fmt.Sprintf("\033[%dm%s\033[0m", color, s)
	}
	return s
}

//获取打印日志的进程信息
func getpProcInfo() (filename, funcname string, line int) {
	pc, filename, line, ok := runtime.Caller(3)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()      // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)            // .foo
		funcname = strings.TrimPrefix(funcname, ".") // foo
		filename = filepath.Base(filename)           // /full/path/basename.go => basename.go
	}
	return
}
