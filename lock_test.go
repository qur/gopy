package pytesting

import (
	"lime/3rdparty/libs/gopy/lib"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func test2() {
	l := py.NewLock()
	defer l.Unlock()
}

func test(wg *sync.WaitGroup) {
	t := time.Now()
	for time.Since(t) < time.Second*2 {
		func() {
			l := py.NewLock()
			defer l.Unlock()
			test2()
			time.Sleep(time.Duration(float64(time.Millisecond) * (1 + rand.Float64())))
		}()
	}
	wg.Done()
}

func TestLock(t *testing.T) {
	l := py.InitAndLock()
	l.Unlock()
	defer func() {
		l.Lock()
		py.Finalize()
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(4)
	go test(&wg)
	go test(&wg)
	go test(&wg)
	test(&wg)
	wg.Wait()
}
