package pytesting

import (
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/qur/gopy/lib"
)

func test2() {
	l := py.NewLock()
	defer l.Unlock()
}

func test() {
	t := time.Now()
	for time.Since(t) < time.Second*2 {
		func() {
			l := py.NewLock()
			defer l.Unlock()
			test2()
			time.Sleep(time.Duration(float64(time.Millisecond) * (1 + rand.Float64())))
		}()
	}
}

func TestLock(t *testing.T) {
	l := py.InitAndLock()
	l.Unlock()
	defer l.Lock()
	runtime.GOMAXPROCS(runtime.NumCPU())
	go test()
	go test()
	go test()
	test()
}
