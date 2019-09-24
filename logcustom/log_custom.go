package log

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type logType int

type mylog struct {
	loggerMap map[logType]*log.Logger
	isColor   bool
}

/*
 	默认设置：
	0 - 2层信息输出到os.stdout
    3 - 6层信息输出到os.stderr
*/
func New() mylog {
	loggerMap := map[logType]*log.Logger{
		Leveltrace: log.New(os.Stdout, "", 0),
		Levelinfo:  log.New(os.Stdout, "", 0),
		Leveldebug: log.New(os.Stdout, "", 0),
		Levelwarn:  log.New(os.Stderr, "", 0),
		Levelerror: log.New(os.Stderr, "", 0),
		Levelpanic: log.New(os.Stderr, "", 0),
		Levelfatal: log.New(os.Stderr, "", 0),
	}
	return mylog{loggerMap, false}
}

/*
	隐藏单个日志级别的输出
	日志级别不能超过Levelpanic
*/
func (l mylog) SetLogDiscard(t logType) error {
	if t > Levelpanic {
		return errors.New("SetLevel err: can't set log level Discard > Levelerror")
	}
	l.loggerMap[t].SetOutput(ioutil.Discard)
	return nil
}

/*
	设置隐藏日志等级
	最高设置到Levelpanic 即是低于panic等级的日志都不显示
*/
func (l mylog) SetLogDiscardLevel(t logType) error {
	if t > Levelpanic {
		return errors.New("SetLevel err: can't set log level Discard more than Levelerror")
	}
	for i := int(t); i >= 0; i-- {
		if i <= 3 {
			l.loggerMap[logType(i)].SetOutput(os.Stdout)
		} else {
			l.loggerMap[logType(i)].SetOutput(os.Stderr)
		}
	}
	for i := 0; i < int(t); i++ {
		l.loggerMap[logType(i)].SetOutput(ioutil.Discard)
	}
	return nil
}

func (l *mylog) IsColor(iscolor bool) {
	l.isColor = iscolor
}

//设置单个日志级别输出到目标位置
func (l mylog) SetOutput(t logType, w io.Writer) {
	l.loggerMap[t].SetOutput(w)
}

//设置全部日志级别输出到目的地
func (l mylog) SetOutputAll(w io.Writer) {
	for _, v := range l.loggerMap {
		v.SetOutput(w)
	}
}

//设置指定日志级别及以上的输出到目标位置
func (l mylog) SetOutputAbove(t logType, w io.Writer) {
	for i := int(t); i < 7; i++ {
		l.loggerMap[logType(i)].SetOutput(w)
	}
}

//设置指定日志级别及以下的输出到目标位置
func (l mylog) SetOutputBelow(t logType, w io.Writer) {
	for i := 0; i <= int(t); i++ {
		l.loggerMap[logType(i)].SetOutput(w)
	}
}
