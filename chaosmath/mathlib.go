// Package chaosmath
// Mode Name: 一些公用的数学方法
// Mode Desc: 模块功能描述
package chaosmath

import (
	"github.com/njtc406/chaosutil/chaoserrors"
	"math/rand"
	"reflect"
	"time"
)

// RandomOneFromSlice 在切片中随机一个元素
//
// slice: 切片
func RandomOneFromSlice(slice interface{}) (interface{}, error) {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return nil, chaoserrors.NewErrCode(-1, "must be a slice", nil)
	}
	sVal := reflect.ValueOf(slice)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := rand.Intn(sVal.Len())
	return sVal.Index(idx).Interface(), nil
}
