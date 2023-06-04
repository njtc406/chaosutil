/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package mysql
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/6/5 0005 0:46
// 最后更新:  yr  2023/6/5 0005 0:46
package mysql

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

func TestNewMysqlClient(t *testing.T) {
	cli, err := NewMysqlClient(&mysql.Config{
		User:         "root",
		Passwd:       "chaospwd",
		Net:          "tcp",
		Addr:         "192.168.2.101:3306",
		DBName:       "game",
		Timeout:      time.Second * 15,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	if err != nil {
		fmt.Println("db", err)
		return
	}

	_, err = cli.conn.QueryContext(context.Background(), "create database game if not exsists charset utf-8")
	if err != nil {
		fmt.Println("conn", err)
		return
	}
}
