/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package define
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/6/5 0005 0:27
// 最后更新:  yr  2023/6/5 0005 0:27
package define

import "fmt"

type MysqlConf struct {
	UserName string
	Password string
	Addr     string
	DataBase string
}

// GetConnPath username:password@protocol(address)/dbname?param=value
// mysql中有个config.FormatDSN可以用来转换,后面看怎么用的
func (m *MysqlConf) GetConnPath() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", m.UserName, m.Password, m.Addr, m.DataBase)
	//return fmt.Sprintf("%s:%s@unix(%s)/%s", m.UserName, m.Password, m.Addr, m.DataBase) // addr是mysql.sock的路径
}
