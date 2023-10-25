/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package mysql
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/6/5 0005 0:07
// 最后更新:  yr  2023/6/5 0005 0:07
package mysql

import (
	"fmt"
	gsql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"runtime"
	"time"
)

func NewMysqlClient(conf *gsql.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.FormatDSN()), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sql, err := db.DB()
	if err != nil {
		return nil, err
	}

	sql.SetMaxOpenConns(runtime.NumCPU() * 10)
	sql.SetMaxIdleConns(runtime.NumCPU() * 2)
	sql.SetConnMaxIdleTime(time.Minute * 5)

	return db, nil
}
