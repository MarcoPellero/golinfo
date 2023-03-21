package main

import (
	"fmt"
)

func main() {
	profile, err := GetProfile("marco_pellero")
	fmt.Println(err)
	fmt.Println(profile)
}
