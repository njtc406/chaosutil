/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package parser
// 模块名: xlsx解析器
// 模块功能简介: 根据传入struct自动将excel表解析成对应的对象列表
package parser

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"reflect"
	"strconv"
)

var filePath = "./devices.xlsx"

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

type myField struct {
	ColName    string
	Value      *reflect.Value
	IsRequired bool
}

func readXlsxFile(s interface{}) (list []interface{}) {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 只读取devices标签页
	sheet, ok := xlFile.Sheet["devices"]
	if !ok {
		fmt.Println("not found devices sheet in xlsx file")
		return
	}

	//tp := reflect.TypeOf(s)
	var realType = reflect.TypeOf(s)
	if realType.Kind() == reflect.Ptr {
		fmt.Println("11111111111111111111111")
		realType = realType.Elem()
	}

	var nameMap = make(map[int]map[int]*myField)
	// 解析列名称
	row1 := sheet.Row(1) // 列名
	//row3 := sheet.Row(2) // 是否必填
	for r := 4; r < sheet.MaxRow; r++ {
		obj := reflect.New(realType)
		//t := reflect.TypeOf(obj)
		//fmt.Println()
		for c := 0; c < len(row1.Cells); c++ {
			content := row1.Cells[c].String()
			//for i := 0; i < t.; i++ {
			//	fmt.Println(t.Field(i).Name)
			//}

			rType := reflect.TypeOf(obj)
			//rValue := reflect.ValueOf(obj)
			for i := 0; i < rType.NumField(); i++ {
				field := rType.Field(i)
				//fmt.Println(field.Type)
				if field.Name == content {
					// 解析内容
					vField := obj.Elem().Field(i)

					realMap, ok := nameMap[r]
					if !ok {
						nameMap[r] = make(map[int]*myField)
						realMap = nameMap[r]
					}
					realMap[c] = &myField{
						ColName:    field.Name,
						Value:      &vField,
						IsRequired: false,
					}
					break
				}
			}
		}
		//for c := 0; c < len(row3.Cells); c++ {
		//	content := row3.Cells[c].String()
		//	if content == "required" {
		//		nameMap[r][c].IsRequired = true
		//	}
		//}
		//list = append(list, obj)
	}

	for r := 4; r < len(sheet.Rows); r++ {
		for c := 0; c < len(sheet.Rows[r].Cells); c++ {
			content := sheet.Rows[r].Cells[c].String()
			realMap, ok := nameMap[r]
			if !ok {
				fmt.Println("unknown row")
				return
			}
			if store, ok := realMap[c]; ok {
				if store.IsRequired && len(content) == 0 {
					fmt.Printf("%s is required, must be not empty", store.ColName)
					return
				}
				switch store.Value.Kind() {
				case reflect.String:
					store.Value.SetString(content)
				case reflect.Int:
					v, e := strconv.ParseInt(content, 10, 10)
					if e != nil {
						fmt.Println(e)
						return
					}
					store.Value.SetInt(v)
				case reflect.Float64:
					v, e := strconv.ParseFloat(content, 8)
					if e != nil {
						fmt.Println(e)
						return
					}
					store.Value.SetFloat(v)
				default:
					fmt.Println("unknown type of value")
				}
			}
		}
	}

	return
}

func xlsxTest() {
	list := readXlsxFile(APEObject{})
	for _, v := range list {
		switch v.(type) {
		case *APEObject:
			val := v.(*APEObject)
			fmt.Printf("%+v", val)
		}
	}
}
