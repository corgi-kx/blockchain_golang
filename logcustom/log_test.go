package log

import (
	"sync"
	"testing"
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
