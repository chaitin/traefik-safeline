package misc

import (
	"sync"
	"testing"
)

var (
	uuids     = make(map[string]uint64)
	mut_uuids sync.Mutex
)

func testOnce(t *testing.T) {
	curr := GenUUID()
	mut_uuids.Lock()
	defer mut_uuids.Unlock()
	if _, ok := uuids[curr]; ok {
		t.Fatalf("duplicate uuid : %s", curr)
	} else {
		uuids[curr] = 1
	}
}

func TestGenUUID(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		testOnce(t)
	}
	rng = NewMT19937()
	for i := 0; i < 1000000; i++ {
		testOnce(t)
	}
}

func TestConcurrentGenUUID(t *testing.T) {
	var wg sync.WaitGroup
	testRoutine := func() {
		for i := 0; i < 500000; i++ {
			testOnce(t)
		}
		wg.Done()
	}
	wg.Add(4)
	go testRoutine()
	go testRoutine()
	go testRoutine()
	go testRoutine()
	wg.Wait()
}
