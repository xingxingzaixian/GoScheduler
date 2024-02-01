package models

import "GoScheduler/internal/modules/global"

// 并发队列
type ConcurrencyQueue struct {
	queue chan struct{}
}

func NewConcurrencyQueue() *ConcurrencyQueue {
	return &ConcurrencyQueue{
		make(chan struct{}, global.Setting.Queue),
	}
}
func (cq *ConcurrencyQueue) Add() {
	cq.queue <- struct{}{}
}

func (cq *ConcurrencyQueue) Done() {
	<-cq.queue
}
