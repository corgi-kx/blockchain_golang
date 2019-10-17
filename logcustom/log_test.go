package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLogPrint(t *testing.T) {
	t.Log("测试打印效果")
	{
		IsColor(false)
		Warn("测试打印效果Warn")
		Info("测试打印效果Info")
		Debug("测试打印效果Info")
		Trace("测试打印效果Trace")
		Error("测试打印效果Error")
		//Panic("测试打印效果Panic")
		//Fatal("测试打印效果Fatal")
		IsColor(true)
		Trace("测试打印效果Trace")
		Tracef("%s", "测试打印效果Tracef")
		Info("测试打印效果Info")
		Infof("%s", "测试打印效果Infof")
		Debug("测试打印效果Info")
		Debugf("%s", "测试打印效果Debugf")
		Warn("测试打印效果Warn")
		Warnf("%s", "测试打印效果Warnf")
		Error("测试打印效果Error")
		Errorf("%s", "测试打印效果Errorf")
		SetLogDiscard(Levelwarn)
		Warn("测试打印效果Warn")
		Warnf("%s", "测试打印效果Warnf")
		//Fatal("测试打印效果Fatal")
		//Panic("测试打印效果Panic")
		mylog := New()
		mylog.Trace("mylog测试打印效果Trace")
		mylog.Tracef("%s", "mylog测试打印效果Tracef")
		mylog.Info("mylog测试打印效果Info")
		mylog.Infof("%s", "mylog测试打印效果Infof")
		mylog.Debug("mylog测试打印效果Info")
		mylog.Debugf("%s", "mylog测试打印效果Debugf")
		mylog.Warn("mylog测试打印效果Warn")
		mylog.Warnf("%s", "mylog测试打印效果Warnf")
		mylog.Error("mylog测试打印效果Error")
		mylog.Errorf("%s", "mylog测试打印效果Errorf")
		mylog2 := New()
		mylog2.IsColor(true)
		mylog2.SetLogDiscardLevel(Levelerror)
		mylog2.Trace("mylog2测试打印效果Trace")
		mylog2.Tracef("%s", "mylog2测试打印效果Tracef")
		mylog2.Info("测试打印效果Info")
		mylog2.Infof("%s", "mylog2测试打印效果Infof")
		mylog2.Debug("测试打印效果Debug")
		mylog2.Debugf("%s", "mylog2测试打印效果Debugf")
		mylog2.Warn("测试打印效果Warn")
		mylog2.Warnf("%s", "mylog2测试打印效果Warnf")
		mylog2.Error("测试打印效果Error")
		mylog2.Errorf("%s", "mylog2测试打印效果Errorf")
	}
}

func TestLogMultithread(t *testing.T) {
	t.Log("测试日志是否线程安全")
	{
		wait := sync.WaitGroup{}
		wait.Add(200)
		for i := 0; i < 100; i++ {
			go func(i int) {
				for j := 0; j <= 10; j++ {
					Tracef("1试试可以正常打印吗！这是第%d次 第%d次！", i, j)
				}
				wait.Done()
			}(i)
		}
		for i := 0; i < 100; i++ {
			go func(i int) {
				for j := 0; j <= 10; j++ {
					Tracef("2试试可以正常打印吗！这是第%d次 第%d次！", i, j)
				}
				wait.Done()
			}(i)
		}
		wait.Wait()
	}
}


func TestMapVSArray(t *testing.T) {
	t.Log("测试字典与切片，使用哪个性能更好")
	count := 1000000

	loggerMap := map[logType]*log.Logger{
		Leveltrace: log.New(os.Stdout, "", 0),
		Levelinfo:  log.New(os.Stdout, "", 0),
		Leveldebug: log.New(os.Stdout, "", 0),
		Levelwarn:  log.New(os.Stderr, "", 0),
		Levelerror: log.New(os.Stderr, "", 0),
		Levelpanic: log.New(os.Stderr, "", 0),
		Levelfatal: log.New(os.Stderr, "", 0),
	}

	loggerArr := []*log.Logger{
		Leveltrace: log.New(os.Stdout, "", 0),
		Levelinfo:  log.New(os.Stdout, "", 0),
		Leveldebug: log.New(os.Stdout, "", 0),
		Levelwarn:  log.New(os.Stderr, "", 0),
		Levelerror: log.New(os.Stderr, "", 0),
		Levelpanic: log.New(os.Stderr, "", 0),
		Levelfatal: log.New(os.Stderr, "", 0),
	}

	fmt.Println("---------------count:", count)

	st := time.Now()
	for i := 0; i < count; i++ {
		assert.NotNil(t, loggerMap[Leveltrace])
		assert.NotNil(t, loggerMap[Levelinfo])
		assert.NotNil(t, loggerMap[Leveldebug])
		assert.NotNil(t, loggerMap[Levelwarn])
		assert.NotNil(t, loggerMap[Levelerror])
		assert.NotNil(t, loggerMap[Levelpanic])
		assert.NotNil(t, loggerMap[Levelfatal])
	}
	fmt.Println("map:", time.Since(st))

	st = time.Now()
	for i := 0; i < count; i++ {
		assert.NotNil(t, loggerArr[Leveltrace])
		assert.NotNil(t, loggerArr[Levelinfo])
		assert.NotNil(t, loggerArr[Leveldebug])
		assert.NotNil(t, loggerArr[Levelwarn])
		assert.NotNil(t, loggerArr[Levelerror])
		assert.NotNil(t, loggerArr[Levelpanic])
		assert.NotNil(t, loggerArr[Levelfatal])
	}
	fmt.Println("array:", time.Since(st))
}


func TestColorCode(t *testing.T) {
	t.Log("测试颜色代码")
	{
		for i:=30;i<=40;i++  {
			fmt.Printf("\033[%sm%s\033[0m\n",strconv.Itoa(i),"some thins you want to print out.")
		}
	}
}

