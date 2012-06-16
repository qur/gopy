#include <ffi.h>
#include "ffi_type.h"

ffi_type *get_ffi_type_pointer(void) { return &ffi_type_pointer; }
ffi_type *get_ffi_type_sint(void) { return &ffi_type_sint; }
ffi_type *get_ffi_type_sint8(void) { return &ffi_type_sint8; }
ffi_type *get_ffi_type_sint16(void) { return &ffi_type_sint16; }
ffi_type *get_ffi_type_sint32(void) { return &ffi_type_sint32; }
ffi_type *get_ffi_type_sint64(void) { return &ffi_type_sint64; }
ffi_type *get_ffi_type_uint(void) { return &ffi_type_uint; }
ffi_type *get_ffi_type_uint8(void) { return &ffi_type_uint8; }
ffi_type *get_ffi_type_uint16(void) { return &ffi_type_uint16; }
ffi_type *get_ffi_type_uint32(void) { return &ffi_type_uint32; }
ffi_type *get_ffi_type_uint64(void) { return &ffi_type_uint64; }
ffi_type *get_ffi_type_float(void) { return &ffi_type_float; }
ffi_type *get_ffi_type_double(void) { return &ffi_type_double; }
