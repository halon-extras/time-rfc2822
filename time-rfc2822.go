package main

// #cgo CFLAGS: -I/opt/halon/include
// #cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-all
// #include <HalonMTA.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"math"
	"sync"
	"time"
	"unsafe"
)

var (
	cached_locations []*time.Location
	lock             = sync.Mutex{}
)

func main() {}

func GetArgumentAsString(args *C.HalonHSLArguments, pos uint64, required bool) (string, error) {
	var x = C.HalonMTA_hsl_argument_get(args, C.ulong(pos))
	if x == nil {
		if required {
			return "", fmt.Errorf("missing argument at position %d", pos)
		} else {
			return "", nil
		}
	}
	var y *C.char
	if C.HalonMTA_hsl_value_get(x, C.HALONMTA_HSL_TYPE_STRING, unsafe.Pointer(&y), nil) {
		return C.GoString(y), nil
	} else {
		return "", fmt.Errorf("invalid argument at position %d", pos)
	}
}

func GetArgumentAsFloat(args *C.HalonHSLArguments, pos uint64, required bool) (float64, error) {
	var x = C.HalonMTA_hsl_argument_get(args, C.ulong(pos))
	if x == nil {
		if required {
			return 0, fmt.Errorf("missing argument at position %d", pos)
		} else {
			return 0, nil
		}
	}
	var y C.double
	if C.HalonMTA_hsl_value_get(x, C.HALONMTA_HSL_TYPE_NUMBER, unsafe.Pointer(&y), nil) {
		return float64(y), nil
	} else {
		return 0, fmt.Errorf("invalid argument at position %d", pos)
	}
}

func SetException(hhc *C.HalonHSLContext, msg string) {
	x := C.CString(msg)
	y := unsafe.Pointer(x)
	defer C.free(y)
	exception := C.HalonMTA_hsl_throw(hhc)
	C.HalonMTA_hsl_value_set(exception, C.HALONMTA_HSL_TYPE_EXCEPTION, y, 0)
}

func SetReturnValueToString(ret *C.HalonHSLValue, val string) {
	x := C.CString(val)
	y := unsafe.Pointer(x)
	defer C.free(y)
	C.HalonMTA_hsl_value_set(ret, C.HALONMTA_HSL_TYPE_STRING, y, 0)
}

func GetCachedLocation(location_string string) (*time.Location, error) {
	lock.Lock()
	defer lock.Unlock()
	for _, cached_location := range cached_locations {
		if cached_location.String() == location_string {
			return cached_location, nil
		}
	}

	location, err := time.LoadLocation(location_string)
	if err == nil {
		cached_locations = append(cached_locations, location)
	}

	return location, err
}

//export Halon_version
func Halon_version() C.int {
	return C.HALONMTA_PLUGIN_VERSION
}

//export time_rfc2822
func time_rfc2822(hhc *C.HalonHSLContext, args *C.HalonHSLArguments, ret *C.HalonHSLValue) {
	unixtime_float64, err := GetArgumentAsFloat(args, 0, false)
	if err != nil {
		SetException(hhc, err.Error())
		return
	}

	location_string, err := GetArgumentAsString(args, 1, false)
	if err != nil {
		SetException(hhc, err.Error())
		return
	}

	var old_time time.Time
	if unixtime_float64 != 0 {
		seconds, decimals := math.Modf(unixtime_float64)
		old_time = time.Unix(int64(seconds), int64(decimals))
	} else {
		old_time = time.Now()
	}

	var value string
	if len(location_string) > 0 {
		location, err := GetCachedLocation(location_string)
		if err != nil {
			SetException(hhc, err.Error())
			return
		}
		new_time := old_time.In(location)
		value = new_time.Format(time.RFC1123Z)
	} else {
		value = old_time.Format(time.RFC1123Z)
	}

	SetReturnValueToString(ret, value)
}

//export Halon_hsl_register
func Halon_hsl_register(hhrc *C.HalonHSLRegisterContext) C.bool {
	time_rfc2822_cs := C.CString("time_rfc2822")
	C.HalonMTA_hsl_register_function(hhrc, time_rfc2822_cs, nil)
	C.HalonMTA_hsl_module_register_function(hhrc, time_rfc2822_cs, nil)
	return true
}
