package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var keyStore map[string]string

func init() {
	println("Initializing the key store...")
	keyStore = make(map[string]string)
}

func addKey(key, value string) {
	keyStore[key] = value
	println("Added key:", key)
}

func getKey(key string) string {
	value, exists := keyStore[key]
	if !exists {
		println("Key not found:", key)
		return ""
	}
	println("Retrieved key:", key, "with value:", value)
	return value
}

func main() {
	println("Hello, World!")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter a command (add/get/exit):")
	for scanner.Scan() {
		command := scanner.Text()
		switch command {
		case "add":
			fmt.Println("Enter key and value (key value):")
			scanner.Scan()
			parts := strings.Split(scanner.Text(), " ")
			if len(parts) == 2 {
				addKey(parts[0], parts[1])
			} else {
				fmt.Println("Invalid input. Please enter key and value separated by a space.")
			}
		case "get":
			fmt.Println("Enter key:")
			scanner.Scan()
			getKey(scanner.Text())
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Unknown command. Please enter add, get, or exit.")
		}
	}
}
