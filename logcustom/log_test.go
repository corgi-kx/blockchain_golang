package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLogMultithread(t *testing.T) {
	t.Log("测试日志是否线程安全")
	{
		wait := sync.WaitGroup{}
		wait.Add(20000)
		for i := 0; i < 10000; i++ {
			go func(i int) {
				for j := 0; j <= 10; j++ {
					Tracef("1试试可以正常打印吗！这是第%d次 第%d次！", i, j)
				}
				wait.Done()
			}(i)
		}
		for i := 0; i < 10000; i++ {
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
	//count := 10000
	/*
		---------------count: 10000
			map: 1.39726ms
			array: 696.726µs
	*/
	//count := 100000
	/*
		---------------count: 100000
		map: 8.799891ms
		array: 4.922019ms
	*/
	//count := 200000
	/*
		---------------count: 200000
		map: 17.465112ms
		array: 10.901369ms
	*/
	count := 500000
	/*
		---------------count: 500000
		map: 42.735246ms
		array: 25.192769ms
	*/
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
