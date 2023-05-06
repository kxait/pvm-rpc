package pvm

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS: -lpvm3
/*
	#include "pvm_wrapper.h"
	#include<stdlib.h>
*/
import "C"
import (
	"reflect"
	"strings"
	"unsafe"
)

func Mytid() (int, error) {
	result := int(C.pvm_mytid())

	dbgln("pvm_mytid() = %d", result)

	if result < 0 {
		return 0, PvmErrorFromInt(result)
	}

	return result, nil
}

func Parent() (int, error) {
	result := int(C.pvm_parent())

	dbgln("pvm_parent() = %d", result)

	if result < 0 {
		return 0, PvmErrorFromInt(result)
	}

	return result, nil
}

func CatchoutStdout() error {
	if result := C.pvm_catchout_stdout(); result < 0 {
		return PvmErrorFromCInt(result)
	}

	dbgln("pvm_catchout(stdout)")

	return nil
}

type SpawnResult struct {
	Numt int
	TIds []int
}

// if error code > 0, then some tasks failed to spawn - check TIds for error results
func Spawn(task string, args []string, flag SpawnOptions, where string, ntask int) (*SpawnResult, error) {
	taskCstr := C.CString(task)
	defer C.free(unsafe.Pointer(taskCstr))

	argsCstrArr := stringSliceToCStringArray(args)
	defer func() {
		for _, c := range argsCstrArr {
			C.free(unsafe.Pointer(c))
		}
	}()

	tIdsCintPtr := (*C.int)(C.malloc(C.sizeof_ulong * C.ulong(ntask)))
	defer C.free(unsafe.Pointer(tIdsCintPtr))

	whereCstr := C.CString(where)
	defer C.free(unsafe.Pointer(whereCstr))

	numtCint := C.pvm_spawn(taskCstr, &argsCstrArr[0], C.int(flag), whereCstr, C.int(ntask), tIdsCintPtr)

	if int(numtCint) < 0 {
		return nil, PvmErrorFromCInt(numtCint)
	}

	tIds := cArrayToSlice(tIdsCintPtr, ntask)

	dbgln("pvm_spawn(%s, %s, %d, %s, %d) = %d, %s", task, strings.Join(args, ", "), flag, where, ntask, numtCint, tIds)

	if int(numtCint) < ntask {
		return &SpawnResult{
			Numt: int(numtCint),
			TIds: tIds,
		}, PvmErrorFromCInt(numtCint)
	}

	return &SpawnResult{
		Numt: int(numtCint),
		TIds: tIds,
	}, nil
}

func Perror(msg string) error {
	msgCstr := C.CString(msg)
	defer C.free(unsafe.Pointer(msgCstr))

	if info := C.pvm_perror(msgCstr); info != 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_perror(%s)", msg)

	return nil
}

func Exit() error {
	if info := C.pvm_exit(); info != 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_exit()")

	return nil
}

func Initsend(encoding DataPackingStyle) (int, error) {
	bufid := C.pvm_initsend(C.int(encoding))

	if bufid < 0 {
		return 0, PvmErrorFromCInt(bufid)
	}

	dbgln("pvm_initsend(%d) = %d", encoding, bufid)

	return int(bufid), nil
}

func Kill(tId int) error {
	info := C.pvm_kill(C.int(tId))

	if info < 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_kill(%d) = %d", tId, info)

	return nil
}

/* PVM_PACKF STATIC BINDINGS */

func PackfString(fmt string, arg string) (int, error) {
	fmtCstr := C.CString(fmt)
	defer C.free(unsafe.Pointer(fmtCstr))

	argCstr := C.CString(arg)
	defer C.free(unsafe.Pointer(argCstr))

	info := C.pvm_packf_string(fmtCstr, argCstr)
	if info < 0 {
		return 0, PvmErrorFromCInt(info)
	}

	dbgln("pvm_packf(%s, %s)", fmt, arg)

	return int(info), nil
}

/* PVM_UNPACKF STATIC BINDINGS */

func UnpackfString(fmt string, buflen int) (string, error) {
	fmtCstr := C.CString(fmt)
	defer C.free(unsafe.Pointer(fmtCstr))

	argCstr := (*C.char)(C.malloc(C.sizeof_char * C.ulong(buflen)))
	defer C.free(unsafe.Pointer(argCstr))

	info := C.pvm_unpackf_string(fmtCstr, argCstr)

	if info < 0 {
		return "", PvmErrorFromCInt(info)
	}

	dbgln("pvm_unpackf(%s, %d) = %s", fmt, buflen, C.GoString(argCstr))

	return C.GoString(argCstr), nil
}

func Send(tid int, msgtag int) error {
	if info := C.pvm_send(C.int(tid), C.int(msgtag)); info < 0 {
		return PvmErrorFromCInt(info)
	}

	dbgln("pvm_send(%d, %d)", tid, msgtag)

	return nil
}

func Recv(tid int, msgtag int) (int, error) {
	info := C.pvm_recv(C.int(tid), C.int(msgtag))

	if info < 0 {
		return 0, PvmErrorFromCInt(info)
	}

	dbgln("pvm_send(%d, %d) = %d", tid, msgtag, info)

	return int(info), nil
}

func Nrecv(tid int, msgtag int) (int, error) {
	info := C.pvm_nrecv(C.int(tid), C.int(msgtag))

	if info < 0 {
		return 0, PvmErrorFromCInt(info)
	}

	dbgln("pvm_nrecv(%d, %d) = %d", tid, msgtag, info)

	return int(info), nil
}

type HostInfo struct {
	HiTid   int
	HiName  string
	HiArch  string
	HiSpeed int
}

type ConfigResult struct {
	Info  int
	Nhost int
	Narch int
	Infos []HostInfo
}

func Config() (*ConfigResult, error) {
	nhostCintPtr := (*C.int)(C.malloc(C.sizeof_int))
	defer C.free(unsafe.Pointer(nhostCintPtr))

	narchCintPtr := (*C.int)(C.malloc(C.sizeof_int))
	defer C.free(unsafe.Pointer(narchCintPtr))

	hostsPtr := C.hostinfo_ptr()

	infoCint := C.pvm_config(nhostCintPtr, narchCintPtr, &hostsPtr)

	result := ConfigResult{
		Info:  int(infoCint),
		Nhost: int(*nhostCintPtr),
		Narch: int(*narchCintPtr),
		Infos: cPvmHostinfoArrayToSlice(hostsPtr, int(*nhostCintPtr)),
	}

	if infoCint != 0 {
		return &result, PvmErrorFromCInt(infoCint)
	}

	return &result, nil
}

type BufinfoResult struct {
	Bytes  int
	MsgTag int
	TId    int
}

func Bufinfo(bufid int) (*BufinfoResult, error) {
	bytesCintPtr := (*C.int)(C.malloc(C.sizeof_ulong))
	defer C.free(unsafe.Pointer(bytesCintPtr))

	msgtagCintPtr := (*C.int)(C.malloc(C.sizeof_ulong))
	defer C.free(unsafe.Pointer(msgtagCintPtr))

	tidCintPtr := (*C.int)(C.malloc(C.sizeof_ulong))
	defer C.free(unsafe.Pointer(tidCintPtr))

	if info := C.pvm_bufinfo(C.int(bufid), bytesCintPtr, msgtagCintPtr, tidCintPtr); info < 0 {
		return nil, PvmErrorFromCInt(info)
	}

	return &BufinfoResult{
		Bytes:  int(*bytesCintPtr),
		MsgTag: int(*msgtagCintPtr),
		TId:    int(*tidCintPtr),
	}, nil
}

func stringSliceToCStringArray(strings []string) []*C.char {
	var result []*C.char
	for _, c := range strings {
		result = append(result, C.CString(c))
	}

	return result
}

func cArrayToSlice(array *C.int, len int) []int {
	var result []int
	var list []C.int
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(array))

	for _, c := range list {
		result = append(result, int(c))
	}

	return result
}

func cPvmHostinfoArrayToSlice(array *C.pvmhostinfo, len int) []HostInfo {
	var result []HostInfo
	var list []C.pvmhostinfo
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(array))

	tidCintPtr := (*C.int)(C.malloc(C.sizeof_int))
	defer C.free(unsafe.Pointer(tidCintPtr))

	speedCintPtr := (*C.int)(C.malloc(C.sizeof_int))
	defer C.free(unsafe.Pointer(speedCintPtr))

	bufsize := 1024
	nameCstrPtr := (*C.char)(C.malloc(C.sizeof_char * C.ulong(bufsize)))
	defer C.free(unsafe.Pointer(nameCstrPtr))

	archCstrPtr := (*C.char)(C.malloc(C.sizeof_char * C.ulong(bufsize)))
	defer C.free(unsafe.Pointer(archCstrPtr))

	for _, c := range list {
		C.unwrap_hostinfo(&c, tidCintPtr, nameCstrPtr, C.int(bufsize), archCstrPtr, C.int(bufsize), speedCintPtr)

		hostinfo := HostInfo{
			HiTid:   int(*tidCintPtr),
			HiName:  C.GoString(nameCstrPtr),
			HiArch:  C.GoString(archCstrPtr),
			HiSpeed: int(*speedCintPtr),
		}

		result = append(result, hostinfo)
	}

	return result
}
