package main

import (
	"fmt"
	"os"
)

func main() {
	token, err := Login(os.Args[1], os.Args[2])
	fmt.Println(err)
	fmt.Println(token)

	data, err := GetTaskSubmissions("ois_intervalxor", token)
	fmt.Println(err)

	for _, sub := range data {
		fmt.Println(sub)
	}
}
