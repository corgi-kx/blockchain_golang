package log

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var logMap = map[logType]*log.Logger{
	Leveltrace: log.New(os.Stdout, "", 0),
	Levelinfo:  log.New(os.Stdout, "", 0),
	Leveldebug: log.New(os.Stdout, "", 0),
	Levelwarn:  log.New(os.Stderr, "", 0),
	Levelerror: log.New(os.Stderr, "", 0),
	Levelpanic: log.New(os.Stderr, "", 0),
	Levelfatal: log.New(os.Stderr, "", 0),
}
var isColor = false

/*
	隐藏单个日志级别的输出
	日志级别不能超过Levelpanic
*/
func SetLogDiscard(t logType) error {
	if t > Levelpanic {
		return errors.New("SetLevel err: can't set log level Discard > Levelerror")
	}
	logMap[t].SetOutput(ioutil.Discard)
	return nil
}

/*
	设置隐藏日志等级
	最高设置到Levelpanic 即是低于panic等级的日志都不显示
*/
func SetLogDiscardLevel(t logType) error {
	if t > Levelpanic {
		return errors.New("SetLevel err: can't set log level Discard more than Levelerror")
	}
	for i := int(t); i >= 0; i-- {
		if i <= 3 {
			logMap[logType(i)].SetOutput(os.Stdout)
		} else {
			logMap[logType(i)].SetOutput(os.Stderr)
		}
	}
	for i := 0; i < int(t); i++ {
		logMap[logType(i)].SetOutput(ioutil.Discard)
	}
	return nil
}

func IsColor(iscolor bool) {
	isColor = iscolor
}

//设置单个日志级别输出到目标位置
func SetOutput(t logType, w io.Writer) {
	logMap[t].SetOutput(w)
}

//设置全部日志级别输出到目的地
func SetOutputAll(w io.Writer) {
	for _, v := range logMap {
		v.SetOutput(w)
	}
}

//设置指定日志级别及以上的输出到目标位置
func SetOutputAbove(t logType, w io.Writer) {
	for i := int(t); i < 7; i++ {
		logMap[logType(i)].SetOutput(w)
	}
}

//设置指定日志级别及以下的输出到目标位置
func SetOutputBelow(t logType, w io.Writer) {
	for i := 0; i <= int(t); i++ {
		logMap[logType(i)].SetOutput(w)
	}
}
