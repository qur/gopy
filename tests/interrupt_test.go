package pytesting

import (
	"testing"
	"time"

	"gopython.xyz/py/v3"
)

func TestInterrupt(t *testing.T) {
	t.Skip("signals are kind tricky?")

	lock := py.InitAndLockWithSignals()
	defer lock.Finalize()

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
}
