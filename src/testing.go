package main

import (
	"fmt"
	"os"
)

func main() {
	token, err := Login(os.Args[1], os.Args[2])
	fmt.Println(err)
	fmt.Println(token)

	submissions, err := GetTaskSubmissions("ois_intervalxor", token)
	fmt.Println(err)

	for _, sub := range submissions {
		fmt.Println(sub)
		details, err := GetSubmissionDetails(sub.SubmissionId, token)
		fmt.Println(err)
		fmt.Println(details)
		fmt.Println()
	}
}
