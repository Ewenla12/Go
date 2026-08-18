package main

import "fmt"

func main() {
	// declare without assigning — Go sets everything to its zero value

	var b bool
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var u uint
	var u8 uint8
	var u16 uint16
	var u32 uint32
	var u64 uint64
	var f32 float32
	var f64 float64
	var c64 complex64
	var c128 complex128
	var s string
	var r rune

	fmt.Println("bool      :", b)    // false
	fmt.Println("int       :", i)    // 0
	fmt.Println("int8      :", i8)   // 0
	fmt.Println("int16     :", i16)  // 0
	fmt.Println("int32     :", i32)  // 0
	fmt.Println("int64     :", i64)  // 0
	fmt.Println("uint      :", u)    // 0
	fmt.Println("uint8     :", u8)   // 0
	fmt.Println("uint16    :", u16)  // 0
	fmt.Println("uint32    :", u32)  // 0
	fmt.Println("uint64    :", u64)  // 0
	fmt.Println("float32   :", f32)  // 0
	fmt.Println("float64   :", f64)  // 0
	fmt.Println("complex64 :", c64)  // (0+0i)
	fmt.Println("complex128:", c128) // (0+0i)
	fmt.Println("string    :", s)    // (empty, nothing prints)
	fmt.Println("rune      :", r)    // 0  (rune is int32, zero is code point 0)

	// proving string is empty and not nil
	fmt.Println("string is empty:", s == "") // true
	fmt.Println("string length  :", len(s))  // 0

	// proving rune zero value is not the character '0' but code point 0
	fmt.Println("rune as char   :", string(r)) // (prints nothing — null character)
	fmt.Println("rune as number :", r)         // 0
}
