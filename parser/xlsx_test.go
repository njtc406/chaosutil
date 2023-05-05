/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package parser
// 模块名: 模块名
// 功能描述: 描述
// 作者:  yr  2023/5/4 0004 23:53
// 最后更新:  yr  2023/5/4 0004 23:53
package parser

import (
	"fmt"
	"testing"
)

type APEObject struct {
	Sn              string  `json:"-"`                         // 设备sn号
	ApeID           string  `json:"ApeID,omitempty"`           // 设备ID(deviceID)
	Name            string  `json:"Name,omitempty"`            // 名称
	Model           string  `json:"Model,omitempty"`           // 型号
	IPAddr          string  `json:"IPAddr,omitempty"`          // IP地址
	IPV6Addr        string  `json:"IPV6Addr,omitempty"`        // IPv6地址
	Port            int     `json:"Port,omitempty"`            // 端口号
	Longitude       float64 `json:"Longitude,omitempty"`       // 经度
	Latitude        float64 `json:"Latitude,omitempty"`        // 纬度
	PlaceCode       string  `json:"PlaceCode,omitempty"`       // 安装地点行政区划代码
	Place           string  `json:"Place,omitempty"`           // 位置名
	OrgCode         string  `json:"OrgCode,omitempty"`         // 管辖单位代码
	CapDirection    int     `json:"CapDirection,omitempty"`    // 车辆抓拍方向
	MonitorDirect   string  `json:"MonitorDirect,omitempty"`   // 监视方向
	MonitorAreaDesc string  `json:"MonitorAreaDesc,omitempty"` // 监视区域说明
	IsOnline        string  `json:"IsOnline,omitempty"`        // 是否在线 (1在线 2离线 9其他)
	OwnerApsID      string  `json:"OwnerApsID,omitempty"`      // 所属采集系统
	UserId          string  `json:"UserId,omitempty"`          // 用户帐号
	Password        string  `json:"Password,omitempty"`        // 口令
	FunctionType    string  `json:"FunctionType,omitempty"`    // 功能类型
	PositionType    string  `json:"PositionType,omitempty"`    // 摄像机位置类型
}

func TestNewXlsxParser(t *testing.T) {
	parser := NewXlsxParser()
	parser.FilePath = "./devices.xlsx"
	parser.StructObjMap["devices"] = APEObject{}
	if err := parser.ParseXlsxFile(); err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range parser.RetObjMap {
		for _, obj := range v {
			fmt.Printf("%+v\n", obj)
		}
	}
}
