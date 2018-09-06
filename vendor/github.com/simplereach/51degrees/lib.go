package fiftyonedegrees

/*
#cgo CFLAGS: -I . -Wimplicit-function-declaration
#cgo darwin LDFLAGS: -lm
#cgo linux LDFLAGS: -lm -lrt
#include "51Degrees.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type FiftyoneDegreesProvider struct {
	provider *C.fiftyoneDegreesProvider
}

type FiftyoneDegreesDataSet struct {
	dataSet *C.fiftyoneDegreesDataSet
}

func NewFiftyoneDegreesProvider(fileName string, properties string, poolSize int, cacheSize int) (*FiftyoneDegreesProvider, error) {
	item := &FiftyoneDegreesProvider{
		provider: new(C.fiftyoneDegreesProvider),
	}

	if poolSize == 0 {
		poolSize = 10
	}
	if cacheSize == 0 {
		cacheSize = 1000
	}

	var cFileName *C.char = C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	var cProperties *C.char = C.CString(properties)
	defer C.free(unsafe.Pointer(cProperties))

	status := C.fiftyoneDegreesInitProviderWithPropertyString(cFileName, item.provider, cProperties, C.int(poolSize), C.int(cacheSize))

	// e_fiftyoneDegrees_DataSetInitStatus.DATA_SET_INIT_STATUS_SUCCESS == 0
	if status != 0 {
		return nil, errors.New(fmt.Sprintln("InitWithPropertyString Error,Status:", status))
	}
	return item, nil
}
func (fdp *FiftyoneDegreesProvider) Close() {
	C.fiftyoneDegreesProviderFree(fdp.provider)
}

func (fdp *FiftyoneDegreesProvider) Parse(userAgent string) string {
	ws := C.fiftyoneDegreesProviderWorksetGet(fdp.provider)
	defer C.fiftyoneDegreesWorksetRelease(ws)

	// needs to be done outside of the inline function call so CGo explicitly frees the memory
	var cUserAgent *C.char = C.CString(userAgent)
	defer C.free(unsafe.Pointer(cUserAgent))

	C.fiftyoneDegreesMatch(ws, cUserAgent)
	resultLength := 50000
	buff := make([]byte, resultLength)
	length := int32(C.fiftyoneDegreesProcessDeviceJSON(ws, (*C.char)(unsafe.Pointer(&buff[0]))))
	result := buff[:length]
	return string(result)
}

func NewFiftyoneDegreesDataSet(fileName, properties string) (*FiftyoneDegreesDataSet, error) {
	var cFileName *C.char = C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	var cProperties *C.char = C.CString(properties)
	defer C.free(unsafe.Pointer(cProperties))

	fdds := &FiftyoneDegreesDataSet{dataSet: new(C.fiftyoneDegreesDataSet)}
	status := C.fiftyoneDegreesInitWithPropertyString(cFileName, fdds.dataSet, cProperties)
	if status != 0 {
		return nil, errors.New(fmt.Sprintln("InitWithPropertyString Error,Status:", status))
	}
	return fdds, nil
}
func (fdds *FiftyoneDegreesDataSet) Close() {
	C.fiftyoneDegreesDestroy(fdds.dataSet)
}

func (fdds *FiftyoneDegreesDataSet) Parse(userAgent string) string {
	ws := C.fiftyoneDegreesCreateWorkset(fdds.dataSet)
	defer C.fiftyoneDegreesFreeWorkset(ws)

	var cUserAgent *C.char = C.CString(userAgent)
	defer C.free(unsafe.Pointer(cUserAgent))

	C.fiftyoneDegreesMatch(ws, cUserAgent)
	resultLength := 50000
	buff := make([]byte, resultLength)
	length := int32(C.fiftyoneDegreesProcessDeviceJSON(ws, (*C.char)(unsafe.Pointer(&buff[0]))))
	result := buff[:length]
	return string(result)
}
