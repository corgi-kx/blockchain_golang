package log

import "fmt"

//Trace级别的打印信息，输入要打印的内容，会自动换行
func Trace(v ...interface{}) {
	trace := loggers[Leveltrace]
	s := fmt.Sprint(v...)
	message := joint(lTrace, s, colorGreen)
	trace.Println(message)
}

//Trace级别的打印信息，第一个参数输入格式,第二个参数输入要打印的内容，类似fmt.Printf
func Tracef(format string, v ...interface{}) {
	trace := loggers[Leveltrace]
	s := fmt.Sprintf(format, v...)
	message := joint(lTrace, s, colorGreen)
	trace.Println(message)
}

//Info级别的打印信息，输入要打印的内容，会自动换行
func Info(v ...interface{}) {
	info := loggers[Levelinfo]
	s := fmt.Sprint(v...)
	message := joint(lInfo, s, colorBlue)
	info.Println(message)
}

//Infof级别的打印信息，第一个参数输入格式,第二个参数输入要打印的内容，类似fmt.Printf
func Infof(format string, v ...interface{}) {
	info := loggers[Levelinfo]
	s := fmt.Sprintf(format, v...)
	message := joint(lInfo, s, colorBlue)
	info.Println(message)
}

//Debug级别的打印信息，输入要打印的内容，会自动换行
func Debug(v ...interface{}) {
	debug := loggers[Leveldebug]
	s := fmt.Sprint(v...)
	message := joint(lDebug, s, colorDarkblue)
	debug.Println(message)
}

//Debug级别的打印信息，第一个参数输入格式,第二个参数输入要打印的内容，类似fmt.Printf
func Debugf(format string, v ...interface{}) {
	debug := loggers[Leveldebug]
	s := fmt.Sprintf(format, v...)
	message := joint(lDebug, s, colorDarkblue)
	debug.Println(message)
}

//Warn级别的打印信息，输入要打印的内容，会自动换行
func Warn(v ...interface{}) {
	warn := loggers[Levelwarn]
	s := fmt.Sprint(v...)
	message := joint(lWarn, s, colorYellow)
	warn.Println(message)
}

//Warn级别的打印信息，第一个参数输入格式,第二个参数输入要打印的内容，类似fmt.Printf
func Warnf(format string, v ...interface{}) {
	warn := loggers[Levelwarn]
	s := fmt.Sprintf(format, v...)
	message := joint(lWarn, s, colorYellow)
	warn.Println(message)
}

//Error级别的打印信息，输入要打印的内容，会自动换行
func Error(v ...interface{}) {
	e := loggers[Levelerror]
	s := fmt.Sprint(v...)
	message := joint(lError, s, colorRed)
	e.Println(message)
}

//Errorf级别的打印信息，第一个参数输入格式,第二个参数输入要打印的内容，类似fmt.Printf
func Errorf(format string, v ...interface{}) {
	e := loggers[Levelerror]
	s := fmt.Sprintf(format, v...)
	message := joint(lError, s, colorRed)
	e.Println(message)
}

//Panic级别的打印信息，输入要打印的内容，会自动换行
// 执行Panic，递归执行每层的defer后中断程序
func Panic(v ...interface{}) {
	p := loggers[Levelpanic]
	s := fmt.Sprint(v...)
	message := joint(lPanic, s, colorRed)
	p.Panicln(message)
}

//Panicf级别的打印信息，第一个参数输入格式,第二个参数输入要打印的内容
// 输出错误信息后，执行Panic()，递归执行每层的defer后中断程序
func Panicf(format string, v ...interface{}) {
	p := loggers[Levelpanic]
	s := fmt.Sprintf(format, v...)
	message := joint(lPanic, s, colorRed)
	p.Panicln(message)
}

//Fatal级别的打印信息，输入要打印的内容，会自动换行
// 输出错误信息后，直接执行os.exit(1)中断程序
func Fatal(v ...interface{}) {
	falat := loggers[Levelfatal]
	s := fmt.Sprint(v...)
	message := joint(lFatal, s, colorPurple)
	falat.Fatalln(message)
}

//Fatal级别的打印信息，第一个参数输入格式,第二个参数输入要打印的内容
// 输出错误信息后，直接执行os.exit(1)中断程序
func Fatalf(format string, v ...interface{}) {
	falat := loggers[Levelfatal]
	s := fmt.Sprintf(format, v...)
	message := joint(lFatal, s, colorPurple)
	falat.Fatalln(message)
}

func (l mylog) Trace(v ...interface{}) {
	trace := l.loggers[Leveltrace]
	s := fmt.Sprint(v...)
	message := l.joint(lTrace, s, colorGreen)
	trace.Println(message)
}

func (l mylog) Tracef(format string, v ...interface{}) {
	trace := l.loggers[Leveltrace]
	s := fmt.Sprintf(format, v...)
	message := l.joint(lTrace, s, colorGreen)
	trace.Println(message)
}

func (l mylog) Info(v ...interface{}) {
	info := l.loggers[Levelinfo]
	s := fmt.Sprint(v...)
	message := l.joint(lInfo, s, colorBlue)
	info.Println(message)
}

func (l mylog) Infof(format string, v ...interface{}) {
	info := l.loggers[Levelinfo]
	s := fmt.Sprintf(format, v...)
	message := l.joint(lInfo, s, colorBlue)
	info.Println(message)
}

func (l mylog) Debug(v ...interface{}) {
	debug := l.loggers[Leveldebug]
	s := fmt.Sprint(v...)
	message := l.joint(lDebug, s, colorDarkblue)
	debug.Println(message)
}

func (l mylog) Debugf(format string, v ...interface{}) {
	debug := l.loggers[Leveldebug]
	s := fmt.Sprintf(format, v...)
	message := l.joint(lDebug, s, colorDarkblue)
	debug.Println(message)
}

func (l mylog) Warn(v ...interface{}) {
	warn := l.loggers[Levelwarn]
	s := fmt.Sprint(v...)
	message := l.joint(lWarn, s, colorYellow)
	warn.Println(message)
}

func (l mylog) Warnf(format string, v ...interface{}) {
	warn := l.loggers[Levelwarn]
	s := fmt.Sprintf(format, v...)
	message := l.joint(lWarn, s, colorYellow)
	warn.Println(message)
}

func (l mylog) Error(v ...interface{}) {
	e := l.loggers[Levelerror]
	s := fmt.Sprint(v...)
	message := l.joint(lError, s, colorRed)
	e.Println(message)
}

func (l mylog) Errorf(format string, v ...interface{}) {
	e := l.loggers[Levelerror]
	s := fmt.Sprintf(format, v...)
	message := l.joint(lError, s, colorRed)
	e.Println(message)
}

func (l mylog) Panic(v ...interface{}) {
	p := l.loggers[Levelpanic]
	s := fmt.Sprint(v...)
	message := l.joint(lPanic, s, colorRed)
	p.Panicln(message)
}

func (l mylog) Panicf(format string, v ...interface{}) {
	p := l.loggers[Levelpanic]
	s := fmt.Sprintf(format, v...)
	message := l.joint(lPanic, s, colorRed)
	p.Panicln(message)
}

func (l mylog) Fatal(v ...interface{}) {
	falat := l.loggers[Levelfatal]
	s := fmt.Sprint(v...)
	message := l.joint(lFatal, s, colorPurple)
	falat.Fatalln(message)
}

func (l mylog) Fatalf(format string, v ...interface{}) {
	falat := l.loggers[Levelfatal]
	s := fmt.Sprintf(format, v...)
	message := l.joint(lFatal, s, colorPurple)
	falat.Fatalln(message)
}
