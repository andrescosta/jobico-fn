package sdk

// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"
)

func log(message string) {
	ptr, size := stringToPtr(message)
	_log(ptr, size)
	runtime.KeepAlive(message) // keep message alive until ptr is no longer needed.
}

//go:wasmimport env log
func _log(ptr, size uint32)

func ptrToString(ptr uint32, size uint32) string {
	return unsafe.String((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

func stringToPtr(s string) (uint32, uint32) {
	ptr := unsafe.Pointer(unsafe.StringData(s))
	return uint32(uintptr(ptr)), uint32(len(s))
}

func stringToLeakedPtr(s string) (uint32, uint32) {
	size := C.ulong(len(s))
	ptr := unsafe.Pointer(C.malloc(size))
	copy(unsafe.Slice((*byte)(ptr), size), s)
	return uint32(uintptr(ptr)), uint32(size)
}
