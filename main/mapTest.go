package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var s sync.RWMutex

func main() {

	tMap := make(map[int]bool)

	for i := 0; i < 1000; i++ {
		if i%2 == 0 {
			tMap[i] = false
		} else {
			tMap[i] = true
		}
	}
	go read(tMap)
	go read(tMap)
	go read(tMap)
	go write(tMap)
	select {}
}

func write(tMap map[int]bool) {
	for {
		s.Lock()
		randIndex := rand.Int() % 1000
		if rand.Int()%2 == 0 {
			tMap[randIndex] = true
		} else {
			tMap[randIndex] = false
		}
		s.Unlock()
	}
}

func read(tMap map[int]bool) {
	var sum = 0
	for {
		s.RLock()
		randIndex := rand.Int() % 1000
		if tMap[randIndex] {
			sum ++
		}
		if sum%10000 == 0 {
			fmt.Println(sum)
		}
		s.RUnlock()
	}
}
