package rpc

import (
	"sync"
	"time"
)

var locM sync.Mutex
var locMap = make(map[string]*time.Location)

var UTC8, _ = time.LoadLocation("Asia/Shanghai")

func LoadLocation(location string) (*time.Location, error) {
	if location == "" {
		return UTC8, nil
	}
	locM.Lock()
	defer locM.Unlock()
	loc, ok := locMap[location]
	if ok {
		return loc, nil
	}
	loc, err := time.LoadLocation(location)
	if err != nil {
		return nil, err
	}
	locMap[location] = loc
	return loc, nil
}
