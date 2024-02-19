package models

import (
	"github.com/spf13/viper"
)

// 并发队列
type ConcurrencyQueue struct {
	queue chan struct{}
}

func NewConcurrencyQueue() *ConcurrencyQueue {
	return &ConcurrencyQueue{
		make(chan struct{}, viper.GetInt("queue")),
	}
}
func (cq *ConcurrencyQueue) Add() {
	cq.queue <- struct{}{}
}

func (cq *ConcurrencyQueue) Done() {
	<-cq.queue
}
