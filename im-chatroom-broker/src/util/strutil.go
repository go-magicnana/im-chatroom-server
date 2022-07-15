package util

import (
	"reflect"
	"strings"
	"unsafe"
)

func B2s(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

func S2b(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func IsEmpty(s string) bool {
	if s == "" {
		return true
	}

	return false
}

func IsNotEmpty(s string) bool{
	return !IsEmpty(s)
}

func StartWith(s string, sub string) bool {
	return strings.Index(s, sub) == 0
}


