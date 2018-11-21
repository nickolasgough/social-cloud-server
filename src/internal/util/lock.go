package util

import (
	"sync"
	"fmt"
)


type model struct {
	id string
	lock sync.Mutex
}

var models = []model{
	{id: "profile", lock: sync.Mutex{}},
	{id: "post", lock: sync.Mutex{}},
	{id: "notification", lock: sync.Mutex{}},
	{id: "connection", lock: sync.Mutex{}},
	{id: "feed", lock: sync.Mutex{}},
	{id: "reaction", lock: sync.Mutex{}},
}

func AcquireLocks(lockIds []string) {
	for _, i := range lockIds {
		acquireLock(i)
	}
	fmt.Printf("Acquiring locks for %+v\n", lockIds)
}

func acquireLock(lockId string) {
	for _, m := range models {
		if lockId == m.id {
			m.lock.Lock()
			return
		}
	}
	panic("Invalid lock id")
}

func ReleaseLocks(lockIds []string) {
	fmt.Printf("Releasing locks for %+v\n", lockIds)
	for _, i := range lockIds {
		releaseLock(i)
	}
}

func releaseLock(lockId string) {
	for _, m := range models {
		if lockId == m.id {
			m.lock.Lock()
			return
		}
	}
	panic("Invalid lock id")
}