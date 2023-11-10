package utils

import "sync"

var mutexMap sync.Map

func GetMutex(name string) *sync.Mutex {
	mutex, ok := mutexMap.Load(name)
	if !ok {
		newMutex := &sync.Mutex{}
		mutex, _ = mutexMap.LoadOrStore(name, newMutex)
	}
	return mutex.(*sync.Mutex)
}
