package log

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

/*
默认设置：
0 - 2层信息输出到os.stdout
3 - 6层信息输出到os.stderr
*/
var loggers = []*log.Logger{
	Leveltrace: log.New(os.Stdout, "", 0),
	Levelinfo:  log.New(os.Stdout, "", 0),
	Leveldebug: log.New(os.Stdout, "", 0),
	Levelwarn:  log.New(os.Stderr, "", 0),
	Levelerror: log.New(os.Stderr, "", 0),
	Levelpanic: log.New(os.Stderr, "", 0),
	Levelfatal: log.New(os.Stderr, "", 0),
}

//是否彩色打印
func IsColor(iscolor bool) {
	isColor = iscolor
}

//SetLogDiscard隐藏单个日志级别的输出信息
//传入需要隐藏的日志级别，日志级别不能超过Levelpanic
func SetLogDiscard(t logType) error {
	if t > Levelpanic {
		return errors.New("SetLevel err: can't set log level Discard > Levelerror")
	}
	loggers[t].SetOutput(ioutil.Discard)
	return nil
}

//SetLogDiscardLevel隐藏多个日志级别输出信息
//传入需要隐藏的日志级别，最高设置到Levelpanic 即是低于panic等级的日志都不显示
func SetLogDiscardLevel(t logType) error {
	if t > Levelpanic {
		return errors.New("SetLevel err: can't set log level Discard more than Levelerror")
	}
	for i := int(t); i >= 0; i-- {
		if i <= 3 {
			loggers[logType(i)].SetOutput(os.Stdout)
		} else {
			loggers[logType(i)].SetOutput(os.Stderr)
		}
	}
	for i := 0; i < int(t); i++ {
		loggers[logType(i)].SetOutput(ioutil.Discard)
	}
	return nil
}

//SetOutput设置单个日志级别输出到目标位置
//传入文件的句柄（或者实现了io.Writer接口的对象）与日志级别，则该日志级别的日志将会输出到指定的文件或位置
func SetOutput(w io.Writer, t logType) {
	loggers[t].SetOutput(w)
}

//SetOutputAll设置全部日志级别输出到目标位置
//传入文件的句柄（或者实现了io.Writer接口的对象）与日志级别，则全部日志级别的日志将会输出到指定的文件或位置
func SetOutputAll(w io.Writer) {
	for i := range loggers {
		loggers[i].SetOutput(w)
	}
}

//SetOutputAbove设置指定日志级别及以上的输出到目标位置
//传入文件的句柄（或者实现了io.Writer接口的对象）与日志级别，则该日志级别以上的日志（包括此日志级别）将会输出到指定的文件或位置
func SetOutputAbove(w io.Writer, t logType) {
	for i := int(t); i < 7; i++ {
		loggers[logType(i)].SetOutput(w)
	}
}

//SetOutputBelow设置指定日志级别及以下的输出到目标位置
//传入文件的句柄（或者实现了io.Writer接口的对象）与日志级别，则该日志级别以下的日志（包括此日志级别）将会输出到指定的文件或位置
func SetOutputBelow(w io.Writer, t logType) {
	for i := 0; i <= int(t); i++ {
		loggers[logType(i)].SetOutput(w)
	}
}
