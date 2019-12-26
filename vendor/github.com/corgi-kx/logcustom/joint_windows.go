package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

//传入要打印的信息，拼接好后返回
func joint(prefix, message string, color colorType) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	filename, funcname, line := getpProcInfo()
	s := fmt.Sprint(prefix, ": ", now, " ", filename, ":", line, ":", funcname, ": ", message)
	//如果是windows系统彩色打印，需要调试kernel32.dll文件
	if isColor {
		winKernelOpen(color)
	}
	return s
}

func (l mylog) joint(prefix, message string, color colorType) string {
	now := time.Now().Format("2006/01/02 15:04:05")
	filename, funcname, line := getpProcInfo()
	s := fmt.Sprint(prefix, ": ", now, " ", filename, ":", line, ":", funcname, ": ", message)
	if l.isColor {
		winKernelOpen(color)
	}
	return s
}

//传入颜色代码，windows系统开启彩色打印
func winKernelOpen(color colorType) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	proc.Call(uintptr(syscall.Stdout), uintptr(color))
}

//windows系统关闭彩色打印
func winKernelColse() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("SetConsoleTextAttribute")
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(7))
	CloseHandle := kernel32.NewProc("CloseHandle")
	CloseHandle.Call(handle)
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
