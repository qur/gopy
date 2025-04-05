//go:build go1.23

package py

import "iter"

// Range returns a Go iterator that iterates through the values of the given
// Iterator. If an error is returned by the Iterator's Next method, then the
// iterator will panic.
func Range(i Iterator) iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for {
			obj, err := i.Next()
			if err != nil {
				panic(err)
			}

			if obj == nil || !yield(obj) {
				return
			}
		}
	}
}
