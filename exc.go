// Code generated by gen_exc.py. DO NOT EDIT.

// This file is automatically generated.  To regenerate:
//   go generate ./...

package py

// #include "utils.h"
import "C"

var (
	BaseException             = newException(C.PyExc_BaseException)
	Exception                 = newException(C.PyExc_Exception)
	BaseExceptionGroup        = newException(C.PyExc_BaseExceptionGroup)
	StopAsyncIteration        = newException(C.PyExc_StopAsyncIteration)
	StopIteration             = newException(C.PyExc_StopIteration)
	GeneratorExit             = newException(C.PyExc_GeneratorExit)
	ArithmeticError           = newException(C.PyExc_ArithmeticError)
	LookupError               = newException(C.PyExc_LookupError)
	AssertionError            = newException(C.PyExc_AssertionError)
	AttributeError            = newException(C.PyExc_AttributeError)
	BufferError               = newException(C.PyExc_BufferError)
	EOFError                  = newException(C.PyExc_EOFError)
	FloatingPointError        = newException(C.PyExc_FloatingPointError)
	OSError                   = newException(C.PyExc_OSError)
	ImportError               = newException(C.PyExc_ImportError)
	ModuleNotFoundError       = newException(C.PyExc_ModuleNotFoundError)
	IndexError                = newException(C.PyExc_IndexError)
	KeyError                  = newException(C.PyExc_KeyError)
	KeyboardInterrupt         = newException(C.PyExc_KeyboardInterrupt)
	MemoryError               = newException(C.PyExc_MemoryError)
	NameError                 = newException(C.PyExc_NameError)
	OverflowError             = newException(C.PyExc_OverflowError)
	RuntimeError              = newException(C.PyExc_RuntimeError)
	RecursionError            = newException(C.PyExc_RecursionError)
	NotImplementedError       = newException(C.PyExc_NotImplementedError)
	SyntaxError               = newException(C.PyExc_SyntaxError)
	IndentationError          = newException(C.PyExc_IndentationError)
	TabError                  = newException(C.PyExc_TabError)
	ReferenceError            = newException(C.PyExc_ReferenceError)
	SystemError               = newException(C.PyExc_SystemError)
	SystemExit                = newException(C.PyExc_SystemExit)
	TypeError                 = newException(C.PyExc_TypeError)
	UnboundLocalError         = newException(C.PyExc_UnboundLocalError)
	UnicodeError              = newException(C.PyExc_UnicodeError)
	UnicodeEncodeError        = newException(C.PyExc_UnicodeEncodeError)
	UnicodeDecodeError        = newException(C.PyExc_UnicodeDecodeError)
	UnicodeTranslateError     = newException(C.PyExc_UnicodeTranslateError)
	ValueError                = newException(C.PyExc_ValueError)
	ZeroDivisionError         = newException(C.PyExc_ZeroDivisionError)
	BlockingIOError           = newException(C.PyExc_BlockingIOError)
	BrokenPipeError           = newException(C.PyExc_BrokenPipeError)
	ChildProcessError         = newException(C.PyExc_ChildProcessError)
	ConnectionError           = newException(C.PyExc_ConnectionError)
	ConnectionAbortedError    = newException(C.PyExc_ConnectionAbortedError)
	ConnectionRefusedError    = newException(C.PyExc_ConnectionRefusedError)
	ConnectionResetError      = newException(C.PyExc_ConnectionResetError)
	FileExistsError           = newException(C.PyExc_FileExistsError)
	FileNotFoundError         = newException(C.PyExc_FileNotFoundError)
	InterruptedError          = newException(C.PyExc_InterruptedError)
	IsADirectoryError         = newException(C.PyExc_IsADirectoryError)
	NotADirectoryError        = newException(C.PyExc_NotADirectoryError)
	PermissionError           = newException(C.PyExc_PermissionError)
	ProcessLookupError        = newException(C.PyExc_ProcessLookupError)
	TimeoutError              = newException(C.PyExc_TimeoutError)
	EnvironmentError          = newException(C.PyExc_EnvironmentError)
	IOError                   = newException(C.PyExc_IOError)
	Warning                   = newException(C.PyExc_Warning)
	UserWarning               = newException(C.PyExc_UserWarning)
	DeprecationWarning        = newException(C.PyExc_DeprecationWarning)
	PendingDeprecationWarning = newException(C.PyExc_PendingDeprecationWarning)
	SyntaxWarning             = newException(C.PyExc_SyntaxWarning)
	RuntimeWarning            = newException(C.PyExc_RuntimeWarning)
	FutureWarning             = newException(C.PyExc_FutureWarning)
	ImportWarning             = newException(C.PyExc_ImportWarning)
	UnicodeWarning            = newException(C.PyExc_UnicodeWarning)
	BytesWarning              = newException(C.PyExc_BytesWarning)
	EncodingWarning           = newException(C.PyExc_EncodingWarning)
	ResourceWarning           = newException(C.PyExc_ResourceWarning)
)
