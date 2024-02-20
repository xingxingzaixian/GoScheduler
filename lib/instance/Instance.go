package instance

import "sync"

// 任务ID作为Key
type Instance struct {
	m sync.Map
}

// 是否有任务处于运行中
func (i *Instance) Has(key uint) bool {
	_, ok := i.m.Load(key)

	return ok
}

func (i *Instance) Add(key uint) {
	i.m.Store(key, struct{}{})
}

func (i *Instance) Done(key uint) {
	i.m.Delete(key)
}
