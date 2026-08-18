package main

import "fmt"

// didn't work cause of float typecasting and i had to use only the asteric button for multiplication and not the x button
// func main() {
// 	var num1, num2 float64
// 	var operator string

// 	fmt.Println(" Calculator ")

// 	// for  first number
// 	fmt.Print("Enter first number: ")
// 	fmt.Scan(&num1)

// 	//  operator
// 	fmt.Print("Enter operator (+, -, *, /, %): ")
// 	fmt.Scan(&operator)

// 	// for second number
// 	fmt.Print("Enter second number: ")
// 	fmt.Scan(&num2)

// 	//  calculate
// 	switch operator {

// 	case "+":
// 		fmt.Printf("%.2f + %.2f = %.2f\n", num1, num2, num1+num2)

// 	case "-":
// 		fmt.Printf("%.2f - %.2f = %.2f\n", num1, num2, num1-num2)

// 	case "*":
// 		fmt.Printf("%.2f * %.2f = %.2f\n", num1, num2, num1*num2)

// 	case "/":
// 		if num2 == 0 {
// 			fmt.Println("Error: Cannot divide by zero!")
// 		} else {
// 			fmt.Printf("%.2f / %.2f = %.2f\n", num1, num2, num1/num2)
// 		}

// 	case "%":
// 		// Modulus only works with integers
// 		fmt.Println("Modulus (%) only works with integers.")

// 	default:
// 		fmt.Println("Invalid operator!")
// 	}
// }

// tottally forgot to and the percenatge cayse i was debating if it was really needed
// func main() {
// 	var num1, num2 float64
// 	var operator string

// 	fmt.Println(" Go Calculator ")

// 	// Get first number
// 	fmt.Print("Enter first number: ")
// 	fmt.Scan(&num1)

// 	// Get operator
// 	fmt.Print("Enter operator (+, -, *, /): ")
// 	fmt.Scan(&operator)

// 	// Get second number
// 	fmt.Print("Enter second number: ")
// 	fmt.Scan(&num2)

// 	// Perform calculation
// 	switch operator {
// 	case "+":
// 		fmt.Printf("%.2f + %.2f = %.2f\n", num1, num2, num1+num2)

// 	case "-":
// 		fmt.Printf("%.2f - %.2f = %.2f\n", num1, num2, num1-num2)

// 	case "*":
// 		fmt.Printf("%.2f * %.2f = %.2f\n", num1, num2, num1*num2)

// 	case "/":
// 		if num2 == 0 {
// 			fmt.Println("Error: Cannot divide by zero.")
// 		} else {
// 			fmt.Printf("%.2f / %.2f = %.2f\n", num1, num2, num1/num2)
// 		}

// 	default:
// 		fmt.Println("Invalid operator.")
// 	}
// }

func main() {
	var choice string

	fmt.Println("go calculator")

	for {
		var num1, num2 float64
		var operator string
		var result float64

		fmt.Print("\nEnter first number: ")
		fmt.Scan(&num1)

		fmt.Print("Enter operator (+, -, *, x, X, /, %): ")
		fmt.Scan(&operator)

		fmt.Print("Enter second number: ")
		fmt.Scan(&num2)

		switch operator {
		case "+":
			result = num1 + num2
			fmt.Printf("Result: %.2f\n", result)

		case "-":
			result = num1 - num2
			fmt.Printf("Result: %.2f\n", result)

		case "*", "x", "X":
			result = num1 * num2
			fmt.Printf("Result: %.2f\n", result)

		case "/":
			if num2 == 0 {
				fmt.Println("Error: Cannot divide by zero.")
			} else {
				result = num1 / num2
				fmt.Printf("Result: %.2f\n", result)
			}

		case "%":
			result = (num1 / 100) * num2
			fmt.Printf("%.2f%% of %.2f = %.2f\n", num1, num2, result)

		default:
			fmt.Println("Invalid operator!")
		}

		fmt.Print("\nDo another calculation? (y/n): ")
		fmt.Scan(&choice)

		if choice == "n" || choice == "N" {
			fmt.Println("bye!!!!!!!!!!!!!!!!!!!!!!!")
			break
		}
	}
}

// got the percentage to work and added a loop till user decides to.............
// and i fogort to add it so that a normal x can do multiplication instead of starting all over i just added it

// tried to use strust here cause it accpect all data types
//  struct method
// type Calculator struct {
// 	Num1 float64
// 	Num2 float64
// }

// // Methods
// func (c Calculator) Add() float64 {
// 	return c.Num1 + c.Num2
// }

// func (c Calculator) Subtract() float64 {
// 	return c.Num1 - c.Num2
// }

// func (c Calculator) Multiply() float64 {
// 	return c.Num1 * c.Num2
// }

// func (c Calculator) Divide() (float64, error) {
// 	if c.Num2 == 0 {
// 		return 0, fmt.Errorf("cannot divide by zero")
// 	}
// 	return c.Num1 / c.Num2, nil
// }

// func (c Calculator) Percentage() float64 {
// 	return (c.Num1 / 100) * c.Num2
// }

// func main() {
// 	var calc Calculator
// 	var operator string
// 	var choice string

// 	fmt.Println("Struct calculator`")

// 	for {
// 		fmt.Print("\nEnter first number:")
// 		fmt.Scan(&calc.Num1)

// 		fmt.Print("Enter operator (+, -, *, x, X, /, %): ")
// 		fmt.Scan(&operator)

// 		fmt.Print("Enter second number:")
// 		fmt.Scan(&calc.Num2)

// 		switch operator {

// 		case "+":
// 			fmt.Printf("Result = %.2f\n", calc.Add())

// 		case "-":
// 			fmt.Printf("Result = %.2f\n", calc.Subtract())

// 		case "*", "x", "X":
// 			fmt.Printf("Result = %.2f\n", calc.Multiply())

// 		case "/":
// 			result, err := calc.Divide()
// 			if err != nil {
// 				fmt.Println(err)
// 			} else {
// 				fmt.Printf("Result = %.2f\n", result)
// 			}

// 		case "%":
// 			fmt.Printf("%.2f%% of %.2f = %.2f\n",
// 				calc.Num1, calc.Num2, calc.Percentage())

// 		default:
// 			fmt.Println("Invalid operator")
// 		}

// 		fmt.Print("\nContinue? (y/n): ")
// 		fmt.Scan(&choice)

// 		if choice == "n" || choice == "N" {
// 			fmt.Println("bye!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
// 			break
// 		}
// 	}
// }
// patched up untill it worked
