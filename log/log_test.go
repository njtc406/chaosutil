package log

import (
	"fmt"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	logger, err := NewDefaultLogger("./run.log", time.Hour*24*7, time.Hour*24, WarnLevelStr, true, false, true, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	start := time.Now()
	logger.Debug("-----------debug test")
	end := time.Now()
	fmt.Println(end.Sub(start))
	//Logs.Fatal("fatal test")
	//Logs.Panic("panic test")
	logger.Info("-----------info test")
	logger.Error("-----------error test")

	Release(logger)
}

func BenchmarkName(b *testing.B) {
	logger, err := NewDefaultLogger("/dev/null", time.Hour*24*7, time.Hour*24, WarnLevelStr, true, false, true, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < b.N; i++ {
		logger.Error("aaaa")
	}

	Release(logger)
}
