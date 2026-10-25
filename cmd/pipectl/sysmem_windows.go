package main

import (
	"syscall"
	"unsafe"
)

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

var (
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procGlobalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
)

// totalSystemMemory returns total physical RAM in bytes, or false if it
// could not be determined.
func totalSystemMemory() (int64, bool) {
	var m memoryStatusEx
	m.Length = uint32(unsafe.Sizeof(m))
	ret, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&m)))
	if ret == 0 {
		return 0, false
	}
	return int64(m.TotalPhys), true
}
