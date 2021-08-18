/*
#cgo CFLAGS: -I/home/bus710/repo/tinygo-nrf52-walkthrough/src/heartrate/include
#include hello.h
*/
package app

import "C"

func Hello() {
	println("hello")
	r := C.print()
	println(r)
}
