package chaoslog

import (
	"fmt"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	Logs := NewDefaultLogger("./run.log", 0, 0, 4, true, false, true)
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	start := time.Now()
	Logs.Logger.Debug("-----------debug test")
	end := time.Now()
	fmt.Println(end.Sub(start))
	//Logs.Fatal("fatal test")
	//Logs.Panic("panic test")
	Logs.Logger.Info("-----------info test")
	Logs.Logger.Error("-----------error test")

	Logs.Close()
}
