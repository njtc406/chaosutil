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
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Client struct {
	db   *sql.DB
	conn *sql.Conn
}

func NewMysqlClient(conf *mysql.Config) (*Client, error) {
	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:   db,
		conn: conn,
	}, nil
}
