package main

import (
	"fmt"
	"log"
)

func main() {
	data, err := GetTaskStats("ois_text2")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}
