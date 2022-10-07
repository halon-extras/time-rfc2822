package main

// #cgo CFLAGS: -I/opt/halon/include
// #cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-all
// #include <HalonMTA.h>
// #include <stdlib.h>
import "C"
import (
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

//export Halon_version
func Halon_version() C.int {
	return C.HALONMTA_PLUGIN_VERSION
}

func set_ret_value(ret *C.HalonHSLValue, value string) {
	value_cs := C.CString(value)
	value_cs_up := unsafe.Pointer(value_cs)
	defer C.free(value_cs_up)
	C.HalonMTA_hsl_value_set(ret, C.HALONMTA_HSL_TYPE_STRING, value_cs_up, 0)
}

func throw_exception(hhc *C.HalonHSLContext, value string) {
	exception := C.HalonMTA_hsl_throw(hhc)
	value_cs := C.CString(value)
	value_cs_up := unsafe.Pointer(value_cs)
	defer C.free(value_cs_up)
	C.HalonMTA_hsl_value_set(exception, C.HALONMTA_HSL_TYPE_EXCEPTION, value_cs_up, 0)
}

func get_cached_location(location_string string) (*time.Location, error) {
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

//export time_rfc2822
func time_rfc2822(hhc *C.HalonHSLContext, args *C.HalonHSLArguments, ret *C.HalonHSLValue) {
	var unixtime_cd C.double
	var unixtime_float64 float64 = 0
	var location_string string
	var location_cs *C.char

	var args_0 = C.HalonMTA_hsl_argument_get(args, 0)
	if args_0 != nil {
		if !C.HalonMTA_hsl_value_get(args_0, C.HALONMTA_HSL_TYPE_NUMBER, unsafe.Pointer(&unixtime_cd), nil) {
			throw_exception(hhc, "Invalid type of \"unixtime\" argument")
			return
		}
		unixtime_float64 = float64(unixtime_cd)
	}

	var args_1 = C.HalonMTA_hsl_argument_get(args, 1)
	if args_1 != nil {
		if !C.HalonMTA_hsl_value_get(args_1, C.HALONMTA_HSL_TYPE_STRING, unsafe.Pointer(&location_cs), nil) {
			throw_exception(hhc, "Invalid type of \"location\" argument")
			return
		}
		location_string = C.GoString(location_cs)
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
		location, err := get_cached_location(location_string)
		if err != nil {
			throw_exception(hhc, err.Error())
			return
		}
		new_time := old_time.In(location)
		value = new_time.Format(time.RFC1123Z)
	} else {
		value = old_time.Format(time.RFC1123Z)
	}

	set_ret_value(ret, value)
}

//export Halon_hsl_register
func Halon_hsl_register(hhrc *C.HalonHSLRegisterContext) C.bool {
	time_rfc2822_cs := C.CString("time_rfc2822")
	C.HalonMTA_hsl_register_function(hhrc, time_rfc2822_cs, nil)
	C.HalonMTA_hsl_module_register_function(hhrc, time_rfc2822_cs, nil)
	return true
}
