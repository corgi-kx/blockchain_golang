package main

import (
	log "github.com/corgi-kx/logcustom"
	"os"
)

func main() {
	log.Info("Write something you want to print !")
	log.Warn("Write something you want to print !")
	log.Trace("Write something you want to print !")
	log.Debug("Write something you want to print !")
	log.Error("Write something you want to print !")

	//设置输出信息隐藏等级
	err := log.SetLogDiscardLevel(log.Leveldebug)
	if err != nil {
		log.Error(err)
	}

	log.Info("SetLogDiscardLevel test  !") //INFO不会被打印
	log.Debug("SetLogDiscardLevel test  !")
	log.Warn("SetLogDiscardLevel test  !")

	//创建新的日志对象
	mylog := log.New()
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
	}
	//将日志信息输出到指定文件
	mylog.SetOutputAbove(file, log.Levelwarn) //WARN及WARN以上级别的日志会输出到指定文件
	mylog.Trace("SetOutputAll test !")
	mylog.Info("SetOutputAll test  !")
	mylog.Debug("SetOutputAll test  !")
	mylog.Warn("SetOutputAll test  !")
	mylog.Error("SetOutputAll test  !")
}

//func DisplayEffect () {
//	log.Info("Write something you want to print !")
//	log.Warn("Write something you want to print !")
//	log.Trace("Write something you want to print !")
//	log.Debug("Write something you want to print !")
//	log.Error("Write something you want to print !")
//	log.Fatal("Write something you want to print !")
//}