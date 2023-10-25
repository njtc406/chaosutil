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
	"fmt"
	"github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

//func TestNewMysqlClient(t *testing.T) {
//	cli, err := NewMysqlClient(&mysql.Config{
//		User:                 "root",
//		Passwd:               "chaospwd",
//		Net:                  "tcp",
//		Addr:                 "192.168.2.101:3306",
//		DBName:               "game",
//		Timeout:              time.Second * 15,
//		ReadTimeout:          time.Second,
//		WriteTimeout:         time.Second,
//		AllowNativePasswords: true,
//	})
//	if err != nil {
//		fmt.Println("db", err)
//		return
//	}
//
//	//_, err = cli.conn.QueryContext(context.Background(), "create database if not exists game default charset utf8 collate utf8_general_ci")
//	//if err != nil {
//	//	fmt.Println("conn", err)
//	//	return
//	//}
//
//	_, err = cli.conn.QueryContext(context.Background(), "CREATE TABLE `game` (`id` int(11) NOT NULL AUTO_INCREMENT, `game_secret` char(32) NOT NULL,  `game_client_id` char(32) DEFAULT NULL,  `game_name` char(32) DEFAULT NULL, `gateWay` char(32) DEFAULT 'tcp',  PRIMARY KEY (`id`)) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8")
//	if err != nil {
//		fmt.Println("conn", err)
//		return
//	}
//
//	row, err := cli.conn.QueryContext(context.Background(), "INSERT INTO `game` VALUES ('1', '21e5c3249f6cf405826027d59773a7a3', '1001436429702627', 'c1', 'tcp')")
//	if err != nil {
//		fmt.Println("conn", err)
//		return
//	}
//
//}

type Game struct {
	ID           uint32 `gorm:"column:id;type:int(11);primaryKey;size:11;autoIncrement"`
	GameSecret   string `gorm:"column:game_secret;type:char(32);not null"`
	GameClientID string `gorm:"column:game_client_id;type:char(32);default:null"`
	GameName     string `gorm:"column:game_name;type:char(32);default:null"`
	GateWay      string `gorm:"column:gateway;type:char(32);default:tcp"`
}

func TestNewClient(t *testing.T) {
	conf := &mysql.Config{
		User:                 "root",
		Passwd:               "chaospwd",
		Net:                  "tcp",
		Addr:                 "192.168.2.101:3306",
		DBName:               "game",
		Timeout:              time.Second * 15,
		ReadTimeout:          time.Second,
		WriteTimeout:         time.Second,
		AllowNativePasswords: true,
	}
	db, err := NewMysqlClient(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	g := &Game{}
	db.Debug().Where("id = ?", 1).First(g)

	fmt.Printf("%#v", g)

	//gameInfo := &Game{
	//	GameSecret:   "12345",
	//	GameClientID: "1",
	//	GameName:     "h1",
	//}
	//
	//db.AutoMigrate(&Game{})
	//db.Create(gameInfo)
	//db.UpdateColumn("gateway", "tcp")
}
