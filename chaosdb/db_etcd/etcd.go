// Package db_etcd
// Mode Name: etcd
// Mode Desc: 模块功能描述
package db_etcd

import (
	"context"
	"fmt"
	"github.com/njtc406/chaosutil/chaoserrors"
	"github.com/njtc406/chaosutil/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	defaultRootPath       = "/chaos"
	defaultDataPath       = "/data"
	defaultLockPath       = "/lock"
	defaultTTL      int64 = 30
)

type ETCDConf struct {
	AddrList []string
	RootPath string
	UserName string
	Password string
	Log      log.ILogger
}

type ETCDClient struct {
	c        *clientv3.Client     // etcd客户端
	session  *concurrency.Session // etcd会话
	conf     clientv3.Config      // 配置
	rootPath string               // etcd数据存储根路径
	lockPath string               // etcd锁存储路径
	closed   bool                 // 是否已关闭
	lockMap  sync.Map             // 锁映射表
	log      log.ILogger          // 日志 (这里使用的是chaosutil的日志,可以修改为系统日志,然后兼容chaosutil的日志)
}

func NewETCDClient() *ETCDClient {
	return &ETCDClient{
		lockMap: sync.Map{},
	}
}

func (ec *ETCDClient) Init(conf *ETCDConf) chaoserrors.CError {
	ec.conf.Endpoints = conf.AddrList

	// 支持无账号连接
	if conf.UserName != "" && conf.Password != "" {
		ec.conf.Username = conf.UserName
		ec.conf.Password = conf.Password
	}

	rootPath := conf.RootPath
	if rootPath == "" {
		rootPath = defaultRootPath
	} else {
		ec.rootPath = rootPath
	}

	if conf.Log == nil {
		// TODO 这里可能需要使用一个默认的log,直接写入到标准输出
		return chaoserrors.NewErrCode(-1, "etcd log is nil", nil)
	}

	return nil
}

func (ec *ETCDClient) Connect() error {
	etcdClient, err := clientv3.New(ec.conf)
	if err != nil {
		return chaoserrors.NewErrCode(-1, "etcd failed to connect", err)
	}

	ec.log.Infof("connect to etcd server, addr: %s", ec.conf.Endpoints)

	ec.c = etcdClient

	ec.newSession()

	// TODO 看是否需要增加OnConnect钩子函数
	return nil
}

// newSession 创建一个新的会话
func (ec *ETCDClient) newSession() error {
	ctx := context.Background()
	sessionCtx, cancel := context.WithCancel(ctx)
	s, err := concurrency.NewSession(ec.c, concurrency.WithTTL(int(defaultTTL)), concurrency.WithContext(sessionCtx))
	if err != nil {
		cancel()
		return chaoserrors.NewErrCode(-1, "etcd failed to create session", err)
	}
	ec.session = s

	go ec.watchSession(sessionCtx, cancel)

	return nil
}

// watchSession 监听会话
func (ec *ETCDClient) watchSession(ctx context.Context, cancel context.CancelFunc) {
	defer cancel() // ensure context is canceled when exiting this function

	select {
	case <-ec.session.Done():
		// Session is expired or closed
		if !ec.closed {
			// If the client wasn't intentionally closed, try to reconnect
			ec.newSession()
		}
	case <-ctx.Done():
		// Context was canceled, so just return
		return
	}
}

// Close 关闭etcd客户端
//
// 关闭前应该保证业务已经退出,否则退出时自动关闭所有的锁可能会导致业务异常
func (ec *ETCDClient) Close() {
	ec.closed = true

	ec.lockMap.Range(func(key, value any) bool {
		_ = ec.Unlock(key.(string))
		return true
	})

	if ec.c != nil {
		_ = ec.c.Close()
	}
}

// Get 获取存储记录
func (ec *ETCDClient) Get(key string, withPrefix bool) (*clientv3.GetResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 根据key查询存储记录
	var opts []clientv3.OpOption
	if withPrefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	resp, err := ec.c.Get(ctx, path.Join(ec.rootPath, defaultDataPath, key), opts...)
	if err != nil {
		return nil, chaoserrors.NewErrCode(-1, fmt.Sprintf("etcd failed to search key[%s]", key), err)
	}

	return resp, nil
}

// Put 存储记录
func (ec *ETCDClient) Put(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 存储数据
	_, err := ec.c.Put(ctx, path.Join(ec.rootPath, defaultDataPath, key), value)
	if err != nil {
		return chaoserrors.NewErrCode(-1, fmt.Sprintf("etcd failed to put key[%s]", key), err)
	}

	return nil
}

// Del 删除存储记录
func (ec *ETCDClient) Del(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 删除数据
	_, err := ec.c.Delete(ctx, path.Join(ec.rootPath, defaultDataPath, key))
	if err != nil {
		return chaoserrors.NewErrCode(-1, fmt.Sprintf("etcd failed to delete key[%s]", key), err)
	}

	return nil
}

// Lock 加锁
func (ec *ETCDClient) Lock(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查找key的存储记录
	if cLock, ok := ec.lockMap.Load(key); !ok {
		mutex := concurrency.NewMutex(ec.session, path.Join(ec.rootPath, defaultLockPath, key))
		ec.lockMap.Store(key, mutex)
		return mutex.Lock(ctx)
	} else {
		err := cLock.(*concurrency.Mutex).Lock(ctx)
		if err != nil && isLeaseNotFoundErr(err) { // 如果错误是关于lease的
			// 从映射表中删除该锁
			ec.lockMap.Delete(key)
			// 重新创建session
			ec.newSession()
			// 重新尝试加锁
			return ec.Lock(key)
		}
		return cLock.(*concurrency.Mutex).Lock(ctx)
	}
}

// Unlock 解锁
func (ec *ETCDClient) Unlock(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查找key的存储记录
	if cLock, ok := ec.lockMap.Load(key); ok {
		return cLock.(*concurrency.Mutex).Unlock(ctx)
	} else {
		return chaoserrors.NewErrCode(-1, fmt.Sprintf("etcd failed to unlock key[%s]", key), nil)
	}
}

// 检查是否是lease未找到的错误
func isLeaseNotFoundErr(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "requested lease not found"))
}
