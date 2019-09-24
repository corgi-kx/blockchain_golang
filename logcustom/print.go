package log

import "fmt"

func Trace(v ...interface{}) {
	trace := logMap[Leveltrace]
	s := fmt.Sprint(v...)
	message := joint(lTrace, s, colorGreen)
	trace.Println(message)
}

func Tracef(format string, v ...interface{}) {
	trace := logMap[Leveltrace]
	s := fmt.Sprintf(format, v...)
	message := joint(lTrace, s, colorGreen)
	trace.Println(message)
}

func Info(v ...interface{}) {
	info := logMap[Levelinfo]
	s := fmt.Sprint(v...)
	message := joint(lInfo, s, colorBlue)
	info.Println(message)
}

func Infof(format string, v ...interface{}) {
	info := logMap[Levelinfo]
	s := fmt.Sprintf(format, v...)
	message := joint(lInfo, s, colorBlue)
	info.Println(message)
}

func Debug(v ...interface{}) {
	debug := logMap[Leveldebug]
	s := fmt.Sprint(v...)
	message := joint(lDebug, s, colorDarkblue)
	debug.Println(message)
}

func Debugf(format string, v ...interface{}) {
	debug := logMap[Leveldebug]
	s := fmt.Sprintf(format, v...)
	message := joint(lDebug, s, colorDarkblue)
	debug.Println(message)
}

func Warn(v ...interface{}) {
	warn := logMap[Levelwarn]
	s := fmt.Sprint(v...)
	message := joint(lWarn, s, colorYellow)
	warn.Println(message)
}

func Warnf(format string, v ...interface{}) {
	warn := logMap[Levelwarn]
	s := fmt.Sprintf(format, v...)
	message := joint(lWarn, s, colorYellow)
	warn.Println(message)
}

func Error(v ...interface{}) {
	e := logMap[Levelerror]
	s := fmt.Sprint(v...)
	message := joint(lError, s, colorRed)
	e.Println(message)
}

func Errorf(format string, v ...interface{}) {
	e := logMap[Levelerror]
	s := fmt.Sprintf(format, v...)
	message := joint(lError, s, colorRed)
	e.Println(message)
}

func Panic(v ...interface{}) {
	p := logMap[Levelpanic]
	s := fmt.Sprint(v...)
	message := joint(lPanic, s, colorMagenta)
	p.Panicln(message)
}

func Panicf(format string, v ...interface{}) {
	p := logMap[Levelpanic]
	s := fmt.Sprintf(format, v...)
	message := joint(lPanic, s, colorMagenta)
	p.Panicln(message)
}

func Fatal(v ...interface{}) {
	falat := logMap[Levelfatal]
	s := fmt.Sprint(v...)
	message := joint(lFatal, s, colorMagenta)
	falat.Fatalln(message)
}

func Fatalf(format string, v ...interface{}) {
	falat := logMap[Levelfatal]
	s := fmt.Sprintf(format, v...)
	message := joint(lFatal, s, colorMagenta)
	falat.Fatalln(message)
}
