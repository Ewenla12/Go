package main

import "fmt"

func main() {
	var b bool = true

	// signed integers
	var i int = 42                      // platform-dependent size (32 or 64 bit)
	var i8 int8 = 127                   // -128 to 127
	var i16 int16 = 32000               // -32,760 to 32767
	var i32 int32 = 2000000000          // -2.1B to 2.1B
	var i64 int64 = 9000000000000000000 // -9.2 quintillion to 9.2 quintillion

	// unsigned
	var u uint = 42                       // 32 or 64-bit(Depends on system)
	var u8 uint8 = 255                    // 0 to 255
	var u16 uint16 = 65535                // 0 to 65,535
	var u32 uint32 = 4294967295           // 0 to 4.2B
	var u64 uint64 = 18446744073709551615 // 0 to 18.4 quintillion

	// float
	var f32 float32 = 3.14             //3.14 is a float literal, but Go infers it as float64 by default, so we need to explicitly declare it as float32
	var f64 float64 = 3.14159265358979 // 15 decimal places of precision

	// complex
	var c1 complex128 = 1 + 2i // 1 = real, 2 = imaginary
	c := 3 + 4i
	c2 := 3 - 4i        // Go infers complex128
	c3 := complex(5, 6) // complex(real, imag) -> 5 + 6i

	// string
	var s string = "hello"

	// rune
	var r rune = 'A' // single quotes for rune literals

	fmt.Println(b, i, i8, i16, i32, i64, u, u8, u16, u32, u64, f32, f64, c, c1, c2, c3, s)
	fmt.Println(r)        // 65  ← the Unicode code point for 'A'
	fmt.Printf("%c\n", r) // A   ← %c prints the actual character
}
