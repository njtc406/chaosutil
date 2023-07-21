/*
 * Copyright (c) 2023. YR. All rights reserved
 */

// Package async
// 模块名: 异步回调模块
// 功能描述: 描述
// 作者:  yr  2023/7/13 0013 20:25
// 最后更新:  yr  2023/7/13 0013 20:25
package async

import (
	"sync"
)

var (
	asyncStopped      int32
	wgAsyncJobWorkers sync.WaitGroup
)

// Callback is a function which will be called after async job is finished with result and error
type Callback func(res interface{}, err error)

//func (ac Callback) callback(res interface{}, err error) {
//	if ac != nil {
//		post.Post(func() {
//			ac(res, err)
//		})
//	}
//}
//
//// Routine is a function that will be executed in the async goroutine and its result and error will be passed to AsyncCallback
//type Routine func() (res interface{}, err error)
//
//// asyncJobWorker 异步任务队列
//type asyncJobWorker struct {
//	jobQueue chan *asyncJobItem
//}
//
//type asyncJobItem struct {
//	routine  Routine
//	callback Callback
//}
//
//func newAsyncJobWorker() *asyncJobWorker {
//	return &asyncJobWorker{
//		jobQueue: make(chan *asyncJobItem, consts.AsyncJobQueueMaxLen),
//	}
//}
//
//func (ajw *asyncJobWorker) appendJob(routine Routine, callback Callback) {
//	ajw.jobQueue <- &asyncJobItem{routine: routine, callback: callback}
//}
//
//func (ajw *asyncJobWorker) run() {
//	defer wgAsyncJobWorkers.Done()
//
//	util.RunPanicedRepeated(func() {
//		for {
//		selectTag:
//			select {
//			case item := <-ajw.jobQueue:
//				if item == nil || item.routine == nil {
//					if item == nil {
//						return
//					}
//
//					break selectTag
//				}
//
//				res, err := item.routine()
//				item.callback.callback(res, err)
//			}
//		}
//	})
//}
//
//var (
//	asyncJobWorkersLock sync.RWMutex
//	asyncJobWorkersMap  = map[string]*asyncJobWorker{}
//)
//
//func getAsyncJobWorker(group string) (ajw *asyncJobWorker) {
//	asyncJobWorkersLock.RLock()
//	ajw, exists := asyncJobWorkersMap[group]
//	asyncJobWorkersLock.RUnlock()
//
//	if ajw != nil || exists {
//		return
//	}
//
//	asyncJobWorkersLock.Lock()
//	defer asyncJobWorkersLock.Unlock()
//
//	ajw, exists = asyncJobWorkersMap[group]
//	if ajw != nil || exists {
//		return
//	}
//
//	ajw = newAsyncJobWorker()
//	asyncJobWorkersMap[group] = ajw
//	wgAsyncJobWorkers.Add(1)
//	go ajw.run()
//
//	return
//}
//
//// AppendAsyncJob append an async job to be executed async (not in the gameserver routine)
//func AppendAsyncJob(group string, routine Routine, callback Callback) {
//	stopped := atomic.LoadInt32(&asyncStopped)
//	if stopped == 1 {
//		return
//	}
//
//	ajw := getAsyncJobWorker(group)
//	if ajw != nil {
//		ajw.appendJob(routine, callback)
//	}
//}
//
//// WaitClear wait for all async job workers to finish (should only be called in the gameserver routine)
//func WaitClear() bool {
//	swapped := atomic.CompareAndSwapInt32(&asyncStopped, 0, 1)
//
//	if !swapped {
//		return false
//	}
//
//	// Close all job queue workers
//	log4j.GetLogger().Info("Waiting for all async job workers to be cleared ...")
//	asyncJobWorkersLock.Lock()
//	for group, alw := range asyncJobWorkersMap {
//		log4j.GetLogger().Info("\tclear %s", group)
//		close(alw.jobQueue)
//		asyncJobWorkersMap[group] = nil
//	}
//	asyncJobWorkersLock.Unlock()
//
//	// wait for all job workers to quit
//	wgAsyncJobWorkers.Wait()
//	return true
//}
