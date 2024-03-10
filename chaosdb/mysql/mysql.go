// Package mysql
// Mode Name: 模块名
// Mode Desc: 模块功能描述
package mysql

import (
	"fmt"
	"github.com/njtc406/chaosutil/chaoserrors"
	"github.com/njtc406/chaosutil/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	syslog "log"
	"runtime"
	"time"
)

// Conf mysql配置
type Conf struct {
	UserName      string          // 数据库用户名
	Passwd        string          // 数据库密码
	Addr          string          // 数据库地址
	DBName        string          // 使用的数据库
	TimeZone      string          // 时区
	SlowThreshold time.Duration   // 慢查询阈值
	SlowLogLevel  logger.LogLevel // 慢查询日志级别
	MaxIdleTime   time.Duration   // 最大空闲时间
	MaxLifetime   time.Duration   // 最大生命周期
	MaxIdleConn   int             // 最大空闲连接数
	MaxOpenConn   int             // 最大打开连接数
	logger        log.ILogger
}

// checkDataBase 检查数据库是否存在，不存在则创建
func checkDataBase(db *gorm.DB, dbName string) chaoserrors.CError {
	// 检查数据库是否存在
	var count int
	db.Raw("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", dbName).Scan(&count)
	if count == 0 {
		// 数据库不存在，创建它
		sql := fmt.Sprintf(`create database if not exists %s default charset utf8mb4 collate utf8mb4_unicode_ci`,
			dbName)
		if err := db.Exec(sql).Error; err != nil {
			return chaoserrors.NewErrCode(-1, "failed to create database", err)
		}
	}

	return nil
}

// initDB 初始化数据库连接
func initDB(conf *Conf) (chaoserrors.CError, *gorm.DB) {
	slowLogger := logger.New(
		syslog.New(conf.logger.GetOutput(), "\n", syslog.LstdFlags),
		logger.Config{
			// 设定慢查询时间阈值为 默认值：200 * time.Millisecond
			SlowThreshold: conf.SlowThreshold,
			// 设置日志级别
			LogLevel: conf.SlowLogLevel,
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8&parseTime=True&loc=%s",
		conf.UserName,
		conf.Passwd,
		conf.Addr,
		conf.TimeZone,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: slowLogger,
	})
	if err != nil {
		return chaoserrors.NewErrCode(-1, "failed to connect to mysql", err), nil
	}

	return nil, db
}

// initConf 初始化配置,没有的配置给默认值
func initConf(conf *Conf) error {
	if conf == nil {
		return chaoserrors.NewErrCode(-1, "mysql conf is nil", nil)
	}
	if conf.Addr == "" {
		return chaoserrors.NewErrCode(-1, "mysql addr is empty", nil)
	}
	if conf.SlowThreshold == 0 {
		conf.SlowThreshold = 200 * time.Millisecond
	}

	if conf.SlowLogLevel == 0 {
		conf.SlowLogLevel = logger.Warn
	}

	if conf.MaxIdleTime == 0 {
		conf.MaxIdleTime = 10 * time.Second
	}

	if conf.MaxLifetime == 0 {
		conf.MaxLifetime = 10 * time.Second
	}

	if conf.MaxIdleConn == 0 {
		conf.MaxIdleConn = 10 * runtime.NumCPU()
	}

	if conf.TimeZone == "" {
		conf.TimeZone = "Local"
	}

	if conf.MaxOpenConn == 0 {
		conf.MaxOpenConn = 2 * runtime.NumCPU()
	}

	return nil
}

// NewClient 创建mysql客户端
func NewClient(conf *Conf) (error, *gorm.DB) {
	if err := initConf(conf); err != nil {
		return err, nil
	}
	err, db := initDB(conf)
	if err != nil {
		return err, nil
	}
	db.Use(dbresolver.Register(dbresolver.Config{ /* xxx */ }).
		SetConnMaxIdleTime(conf.MaxIdleTime).
		SetConnMaxLifetime(conf.MaxLifetime).
		SetMaxIdleConns(conf.MaxIdleConn).
		SetMaxOpenConns(conf.MaxOpenConn))

	if err = checkDataBase(db, conf.DBName); err != nil {
		return err, nil
	}

	return nil, db
}
