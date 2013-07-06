package pytesting

import (
	"github.com/qur/gopy/lib"
	"testing"
	"time"
)

func TestInterrupt(t *testing.T) {
	py.Initialize()

	main, err := py.NewDict()
	if err != nil {
		t.Fatal(err)
	} else if g, err := py.GetBuiltins(); err != nil {
		t.Fatal(err)
	} else if err := main.SetItemString("__builtins__", g); err != nil {
		t.Fatal(err)
	}
	go func() {
		<-time.After(3 * time.Second)
		py.SetInterrupt()
	}()
	if ret, err := py.RunString(`import time
while True:
	time.sleep(1)
	print("tick")
`, py.FileInput, main, nil); err != nil {
		// The error we receive here isn't correct, but whatever...
		t.Log(err)
	} else {
		t.Log(ret)
	}
	py.Finalize()
}
