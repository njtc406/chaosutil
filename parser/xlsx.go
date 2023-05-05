/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package parser
// 模块名: xlsx解析器
// 模块功能简介: 根据传入struct自动将excel表解析成对应的对象列表
package parser

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"reflect"
	"strconv"
	"strings"
)

type xlsxParser struct {
	CanParsePrefix string                   // 可以被解析的标签页前缀 默认为"c_"
	FilePath       string                   // 文件路径
	StructObjMap   map[string]interface{}   // [标签页名称(不包含前缀)]标签对应的结构体值或者指针
	RetObjMap      map[string][]interface{} // [标签页名称(不包含前缀)][]标签对应的结构体指针
}

func NewXlsxParser() *xlsxParser {
	return &xlsxParser{
		CanParsePrefix: canParsePrefix,
		FilePath:       "",
		StructObjMap:   make(map[string]interface{}),
		RetObjMap:      make(map[string][]interface{}),
	}
}

func (x *xlsxParser) ParseXlsxFile() error {
	xlFile, err := xlsx.OpenFile(x.FilePath)
	if err != nil {
		return err
	}

	for _, sheet := range xlFile.Sheets {
		if x.CanParsePrefix == "" || strings.HasPrefix(sheet.Name, x.CanParsePrefix) {
			labelName := strings.TrimPrefix(sheet.Name, x.CanParsePrefix)
			var list []interface{}
			err = parseSheet(sheet, x.StructObjMap[labelName], &list)
			if err != nil {
				return err
			}
			x.RetObjMap[labelName] = list
		}
	}

	return nil
}

type myCell struct {
	colDesc    *string
	colName    *string
	isRequired bool
}

const canParsePrefix = `c_`

// ParseXlsxFile 解析xlsx文件
// xlsx文件中的所有标签页,默认带有"c_"前缀的的标签页才会被解析
// xlsx的文件格式为(第一行为列名的字段中文说明,第二行为struct中的导出字段,第三行为字段选填属性,第四行为字段的类型)
// filePath为文件路径(以执行程序所在位置的相对路径或者根目录为基准的全路径)
// structObj为将要解析的结构对象(值或者指针均可)
// 返回值为传入的*structObj类型的一个切片
func ParseXlsxFile(filePath string, structObjMap map[string]interface{}) (map[string][]interface{}, error) {
	ret := make(map[string][]interface{})
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return ret, err
	}

	for _, sheet := range xlFile.Sheets {
		if strings.HasPrefix(sheet.Name, canParsePrefix) {
			structName := strings.TrimPrefix(sheet.Name, canParsePrefix)
			var list []interface{}
			err = parseSheet(sheet, structObjMap[structName], &list)
			if err != nil {
				return ret, err
			}
			ret[structName] = list
		}
	}

	return ret, err
}

// parseSheet 解析标签页
func parseSheet(sheet *xlsx.Sheet, structObj interface{}, ret *[]interface{}) error {
	var realType = reflect.TypeOf(structObj)
	if realType.Kind() == reflect.Ptr {
		realType = realType.Elem()
	}

	var nameMap = make(map[int]map[int]*reflect.Value)
	var descMap = make(map[int]*myCell)
	// 解析列名称
	row0 := sheet.Row(0) // 字段中文说明
	row2 := sheet.Row(2) // 是否必填
	for c := 0; c < len(row2.Cells); c++ {
		cell, ok := descMap[c]
		if !ok {
			cell = &myCell{}
			descMap[c] = cell
		}
		// 填充必填属性
		required := row2.Cells[c].String()
		if required == "required" {
			cell.isRequired = true
		}
		// 填充中文说明
		desc := row0.Cells[c].String()
		cell.colDesc = &desc
	}

	row1 := sheet.Row(1) // 字段名称
	//row3 := sheet.Row(3) // 字段类型(这个之后会做更多的类型)
	for r := 4; r < sheet.MaxRow; r++ {
		// 创建一个新的对象
		obj := reflect.New(realType).Interface()
		// 获取对象的反射类型
		t := reflect.TypeOf(obj)
		// 获取对象的反射值
		v := reflect.ValueOf(obj)

		// 首先遍历字段名称,建立对应关系
		for c := 0; c < len(row1.Cells); c++ {
			content := row1.Cells[c].String()
			for i := 0; i < t.Elem().NumField(); i++ {
				field := t.Elem().Field(i)
				// 如果配置的字段名称和结构的字段名称相同,那么认为是匹配的
				if field.Name == content {
					// 解析内容
					vField := v.Elem().Field(i)
					// 为正文内容的每一行都预先创建一个结构体,然后将结构体的反射保存下来,在后面读取正文时,按照对应位置直接设置值就好了
					realMap, ok := nameMap[r]
					if !ok {
						nameMap[r] = make(map[int]*reflect.Value)
						realMap = nameMap[r]
					}
					realMap[c] = &vField

					// 填充字段名称
					cell, ok := descMap[c]
					if ok && cell.colName == nil {
						cell.colName = &field.Name
					}
					break
				}
			}
		}

		*ret = append(*ret, obj)
	}

	// 解析正文内容
	for r := 4; r < len(sheet.Rows); r++ {
		for c := 0; c < len(sheet.Rows[r].Cells); c++ {
			content := sheet.Rows[r].Cells[c].String()
			info, ok := descMap[c]
			if !ok {
				return errors.New("unknown row")
			}
			if info.isRequired && len(content) == 0 {
				return errors.New(fmt.Sprintf("%s[%s] is required and cannot be empty\n", *info.colDesc, *info.colName))
			}

			realMap, ok := nameMap[r]
			if !ok {
				return errors.New("unknown row")
			}
			if store, ok := realMap[c]; ok {
				// 这里就是在给之前保存下来的反射设置值,目前只有几个基本类型,之后会逐渐扩充
				switch store.Kind() {
				case reflect.String:
					store.SetString(content)
				case reflect.Int:
					v, err := strconv.ParseInt(content, 10, 10)
					if err != nil {
						return err
					}
					store.SetInt(v)
				case reflect.Float64:
					v, err := strconv.ParseFloat(content, 8)
					if err != nil {
						return err
					}
					store.SetFloat(v)
				default:
					return errors.New("unknown type value")
				}
			}
		}
	}

	return nil
}
