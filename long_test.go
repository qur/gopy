package pytesting

import (
    "lime/3rdparty/libs/gopy/lib"
    "testing"
)

func TestLong(t *testing.T) {
    py.Initialize()
    defer py.Finalize()
    l := py.NewLong(10)
    d, err := py.NewDict()
    if err != nil {
        t.Fatal(err)
    }
    if err = d.SetItemString("gopy.count", l); err != nil {
        t.Fatal(err)
    }
    if l2, err := d.GetItemString("gopy.count"); err != nil {
        t.Fatal(err)
    } else {
        if l3, ok := l2.(*py.Long); !ok {
            t.Fatalf("Unable to cast to long")
        } else {
            if 10 != l3.Int64() {
                t.Fatalf("Bad value")
            }
        }
    }
}
