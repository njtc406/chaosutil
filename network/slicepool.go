/*
 * Copyright (c) 2023. YR. All rights reserved
 */

package network

import (
	"sync"
)

type INetMempool interface {
	MakeByteSlice(size int) []byte
	ReleaseByteSlice(byteBuff []byte) bool
}

type memAreaPool struct {
	minAreaValue int //最小范围值
	maxAreaValue int //最大范围值
	growthValue  int //内存增长值
	pool         []sync.Pool
}

// memAreaPoolList 不同大小的缓存池
var memAreaPoolList = [3]*memAreaPool{
	{
		minAreaValue: 1,
		maxAreaValue: 4096,
		growthValue:  512,
	},
	{
		minAreaValue: 4097,
		maxAreaValue: 40960,
		growthValue:  4096,
	},
	{
		minAreaValue: 40961,
		maxAreaValue: 417792,
		growthValue:  16384,
	},
}

func init() {
	for i := 0; i < len(memAreaPoolList); i++ {
		memAreaPoolList[i].makePool()
	}
}

func NewMemAreaPool() *memAreaPool {
	return &memAreaPool{}
}

func (areaPool *memAreaPool) makePool() {
	poolLen := (areaPool.maxAreaValue - areaPool.minAreaValue + 1) / areaPool.growthValue
	areaPool.pool = make([]sync.Pool, poolLen)
	for i := 0; i < poolLen; i++ {
		memSize := (areaPool.minAreaValue - 1) + (i+1)*areaPool.growthValue
		areaPool.pool[i] = sync.Pool{New: func() interface{} {
			return make([]byte, memSize)
		}}
	}
}

func (areaPool *memAreaPool) makeByteSlice(size int) []byte {
	pos := areaPool.getPosByteSize(size)
	if pos > len(areaPool.pool) || pos == -1 {
		return nil
	}

	return areaPool.pool[pos].Get().([]byte)[:size]
}

func (areaPool *memAreaPool) getPosByteSize(size int) int {
	pos := (size - areaPool.minAreaValue) / areaPool.growthValue
	if pos >= len(areaPool.pool) {
		return -1
	}

	return pos
}

func (areaPool *memAreaPool) releaseByteSlice(byteBuff []byte) bool {
	pos := areaPool.getPosByteSize(cap(byteBuff))
	if pos > len(areaPool.pool) || pos == -1 {
		panic("assert!")
		return false
	}

	areaPool.pool[pos].Put(byteBuff)
	return true
}

func (areaPool *memAreaPool) MakeByteSlice(size int) []byte {
	for i := 0; i < len(memAreaPoolList); i++ {
		if size <= memAreaPoolList[i].maxAreaValue {
			return memAreaPoolList[i].makeByteSlice(size)
		}
	}

	return make([]byte, size)
}

func (areaPool *memAreaPool) ReleaseByteSlice(byteBuff []byte) bool {
	for i := 0; i < len(memAreaPoolList); i++ {
		if cap(byteBuff) <= memAreaPoolList[i].maxAreaValue {
			return memAreaPoolList[i].releaseByteSlice(byteBuff)
		}
	}

	return false
}
