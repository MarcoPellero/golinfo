package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	token, err := Login(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)
}
