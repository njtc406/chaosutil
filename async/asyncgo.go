// Package async
// Mode Name: 异步执行
// Mode Desc: 使用协程池中的协程执行任务,防止出现瞬间创建大量协程,出现性能问题
package async

import (
	"context"
	"fmt"
	"github.com/njtc406/chaosutil/chaoserrors"
	"github.com/panjf2000/ants/v2"
	"github.com/smallnest/rpcx/log"
	//"github.com/smallnest/rpcx/log"
	"runtime"
	"time"
)

// antsPool 协程池
var antsPool *ants.Pool

func init() {
	p, err := ants.NewPool(runtime.NumCPU() * 100)
	if err != nil {
		panic(err)
	}

	antsPool = p
}

func Go(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			errString := fmt.Sprint(r)
			err = chaoserrors.NewErrCode(-1, errString, nil)
			//log.Error()
		}
	}()

	return antsPool.Submit(f)
}

// GoWithTimeout 这是一个同步的方法,如果超时,会返回超时错误(暂时没啥用,因为这个方法是同步的,使用ants毫无意义)
func GoWithTimeout(f func() error, timeout time.Duration) (err error) {
	defer func() {
		if r := recover(); r != nil {
			errString := fmt.Sprint(r)
			err = chaoserrors.NewErrCode(-1, errString, nil)
			log.Error()
		}
	}()
	resultChan := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = antsPool.Submit(func() {
		resultChan <- jobFun(f, ctx)
	})
	if err != nil {
		return err
	}

	select {
	case res := <-resultChan:
		return res
	case <-ctx.Done():
		return ctx.Err() // 这将返回一个超时错误
	}
}

func jobFun(f func() error, ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return f()
		}
	}
}
