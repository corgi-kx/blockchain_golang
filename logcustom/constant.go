package log

//日志等级
const (
	Leveltrace logType = iota //基本输出,
	Levelinfo                 //系统信息
	Leveldebug                //系统调试
	Levelwarn                 //系统警告,提示有可预测的错误
	Levelerror                //程序错误,不影响继续使用
	Levelpanic                //程序异常，递归调用本层及上层defer后，中断程序进程，类似java的异常处理
	Levelfatal                //程序直接结束，打印错误信息后直接调用os.Exit(1)结束程序，不会调用各层defer
)

//日志打印颜色
const (
	colorRed colorType = iota + 91
	colorGreen
	colorYellow
	colorBlue
	colorMagenta  //洋红
	colorDarkblue //碧蓝
)

//日志头部标志
const (
	lTrace = "[TRACE]"
	lInfo  = "[INFO]"
	lDebug = "[DEBUG]"
	lWarn  = "[WARNNING]"
	lError = "[ERROR]"
	lPanic = "[PANIC]"
	lFatal = "[FATAL]"
)
