package main

import (
	"fmt"
)

// Define a struct
type Person struct {
	Name    string
	Age     int
	Country string
}

func main() {
	// Create a map
	codeMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// String variable
	message := "Hello, World!"

	// For loop with a range
	for key, value := range codeMap {
		fmt.Println("Key:", key, "Value:", value)
	}

	// Create an instance of the struct
	person := Person{Name: "John", Age: 30, Country: "USA"}

	// Switch statement
	switch person.Country {
	case "USA":
		fmt.Println("American")
	case "UK":
		fmt.Println("British")
	default:
		fmt.Println("Unknown")
	}

	// If statement
	if person.Age >= 18 {
		fmt.Println(person.Name, "is an adult.")
	} else {
		fmt.Println(person.Name, "is a minor.")
	}

	// Using continue and break in a loop
	for i := 0; i < 5; i++ {
		if i == 2 {
                        // Skip iteration
			continue 
		}
		if i == 4 {
                       // Exit loop
			break 
		}
		fmt.Println(i)
	}

	// Printing a string variable
	fmt.Println(message)
}
