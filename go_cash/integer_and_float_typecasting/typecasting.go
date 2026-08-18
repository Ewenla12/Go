package main

import "fmt"

func main() {
	// ── float64 → int (truncates, not rounds)
	a := int(3.9)
	fmt.Println("float64 → int    :", a) // 3

	// ── int → float64 (safe)
	b := float64(10)
	fmt.Println("int → float64    :", b) // 10

	// ── int64 → int32
	var x int64 = 1000
	c := int32(x)
	fmt.Println("int64 → int32    :", c) // 1000

	// ── int → uint (dangerous if negative)
	var n int = 50
	d := uint(n)
	fmt.Println("int → uint       :", d) // 50

	// ── overflow example
	var big int = 300
	e := int8(big)
	fmt.Println("300 → int8       :", e) // 44 (overflows)

	// ── negative int → uint (wraps to huge number)
	var neg int = -1
	f := uint(neg)
	fmt.Println("-1 → uint        :", f) // 18446744073709551615

	// ── integer division vs float division
	fmt.Println("5 / 2 (int)      :", 5/2)                   // 2
	fmt.Println("5 / 2 (float)    :", float64(5)/float64(2)) // 2.5

	// ── float to int always truncates toward zero
	fmt.Println("int(3.2)         :", int(3.2))  //  3
	fmt.Println("int(3.9)         :", int(3.9))  //  3
	fmt.Println("int(-3.9)        :", int(-3.9)) // -3
}
