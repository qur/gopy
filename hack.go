// +build darwin linux freebsd

package py

/*
#include <signal.h>
#include <sys/types.h>
#include <sys/mman.h>
#include <unistd.h>

extern void cinit();
*/
import "C"

import (
	"log"
	"reflect"
	"unsafe"
)

//export stub
func stub(ptrxx unsafe.Pointer) {
	ptr := uintptr(ptrxx)
	var data []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	header.Data = ptr
	header.Cap = 10
	header.Len = 10

	replacement := []byte{
		0x31, 0xc0, // xor    %eax,%eax
		0xff, 0xc8, // dec    %eax
		0xc3, // ret
	}

	pagesize := C.sysconf(C._SC_PAGE_SIZE)
	if pagesize == -1 {
		log.Fatalln("sysconf claims a -1 page size..")
	}

	start := ptr &^ uintptr(pagesize-1) // align address to page start
	ustart := unsafe.Pointer(start)

	if start+uintptr(pagesize) < ptr+uintptr(len(replacement)) {
		// Just in case the code we want to change spans two pages
		pagesize *= 2
	}
	if err := C.mprotect(ustart, C.size_t(pagesize), C.PROT_READ|C.PROT_WRITE|C.PROT_EXEC); err != 0 {
		log.Fatalln(err)
	}
	copy(data, replacement)

	if err := C.mprotect(ustart, C.size_t(pagesize), C.PROT_READ|C.PROT_EXEC); err != 0 {
		log.Fatalln(err)
	}
}

func init() {
	C.cinit()
	//	stub(uintptr(C.addr()))
}
