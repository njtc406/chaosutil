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
	"time"
)

type Conf struct {
	UserName      string
	Passwd        string
	Addr          string
	DBName        string
	SlowThreshold time.Duration
	SlowLogLevel  logger.LogLevel
	MaxIdleTime   time.Duration
	MaxLifetime   time.Duration
	MaxIdleConn   int
	MaxOpenConn   int
}

func Init(conf *Conf) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8&parseTime=True&loc=Local",
		conf.UserName,
		conf.Passwd,
		conf.Addr,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return chaoserrors.NewErrCode(-1, err.Error(), nil)
	}

	// 检查数据库是否存在
	var count int
	db.Raw("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", conf.DBName).Scan(&count)
	if count == 0 {
		// 数据库不存在，创建它
		sql := fmt.Sprintf(`create database if not exists %s default charset utf8mb4 collate utf8mb4_unicode_ci`,
			conf.DBName)
		if err = db.Exec(sql).Error; err != nil {
			// 创建失败
			return chaoserrors.NewErrCode(-1, err.Error(), nil)
		}
	}

	return nil
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
		conf.SlowThreshold = 5 * time.Second
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
		conf.MaxIdleConn = 10
	}

	if conf.MaxOpenConn == 0 {
		conf.MaxOpenConn = 10
	}

	return nil
}

func NewClient(conf *Conf, cLogger log.ILogger) *gorm.DB {
	if err := initConf(conf); err != nil {
		cLogger.Fatal(err)
	}
	if err := Init(conf); err != nil {
		cLogger.Fatal(err)
	}
	slowLogger := logger.New(
		syslog.New(cLogger.Writer(), "\n", syslog.LstdFlags),
		logger.Config{
			// 设定慢查询时间阈值为 5s（默认值：200 * time.Millisecond）
			SlowThreshold: conf.SlowThreshold,
			// 设置日志级别
			LogLevel: conf.SlowLogLevel,
		},
	)
	// 初始化mysql
	// 数据源格式: username:password@protocol(address)/dbname?param=value
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.UserName, conf.Passwd, conf.Addr, conf.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: slowLogger,
	})
	if err != nil {
		cLogger.Fatal(err)
	}

	cLogger.Infof("connect to mysql server, addr: %s database:%s", conf.Addr, conf.DBName)

	// TODO 需要详细看下文档这里怎么设置
	db.Use(dbresolver.Register(dbresolver.Config{ /* xxx */ }).
		SetConnMaxIdleTime(conf.MaxIdleTime).
		SetConnMaxLifetime(conf.MaxLifetime).
		SetMaxIdleConns(conf.MaxIdleConn).
		SetMaxOpenConns(conf.MaxOpenConn))

	return db
}
