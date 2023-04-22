package pytesting

import (
	"log"
	"testing"

	"gopython.xyz/py/v3"
)

func TestDynamic(t *testing.T) {
	lock := py.InitAndLock()
	defer lock.Finalize()

	unittest, err := py.Import("unittest")
	if err != nil {
		t.Fatal(err)
	}

	testResult, err := unittest.GetAttrString("TestResult")
	if err != nil {
		t.Fatal(err)
	}

	goTestResultType := py.Class{
		Name:     "TestResult",
		BaseType: testResult,
		Flags:    py.ClassHeapType,
	}

	if err := goTestResultType.Create(); err != nil {
		t.Fatal(err)
	}

	log.Printf("TRT: %d", py.RefCount(&goTestResultType))

	tr, err := goTestResultType.CallGo(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("TRT: %d, TR: %d", py.RefCount(&goTestResultType), py.RefCount(tr))

	tr.Decref()

	log.Printf("TRT: %d, TR: %d", py.RefCount(&goTestResultType), py.RefCount(tr))
}
