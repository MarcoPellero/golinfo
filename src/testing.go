package main

import (
	"fmt"
	"os"
)

func avgSubs(username, password string) float64 {
	// gets avg number of subs per task
	token, err := Login(username, password)
	if err != nil {
		panic(err)
	}

	profile, err := GetProfile(username)
	if err != nil {
		panic(err)
	}

	var sum uint64 = 0
	results := make(chan int, len(profile.Scores))

	for _, task := range profile.Scores {
		go func(task ProfileScore) {
			info, err := GetTaskSubmissions(task.Name, token)
			if err != nil {
				panic(err)
			}
			results <- len(info)
		}(task)
	}

	for range profile.Scores {
		sum += uint64(<-results)
	}

	return float64(sum) / float64(len(profile.Scores))
}

func avgScorePerSub(username, password string) float64 {
	// gets avg score of subs per sub

	token, err := Login(username, password)
	if err != nil {
		panic(err)
	}

	profile, err := GetProfile(username)
	if err != nil {
		panic(err)
	}

	var sum uint64 = 0
	results := make(chan float64, len(profile.Scores))

	for _, task := range profile.Scores {
		go func(task ProfileScore) {
			info, err := GetTaskSubmissions(task.Name, token)
			if err != nil {
				panic(err)
			}

			scoreSum := 0
			for _, x := range info {
				scoreSum += x.Score
			}

			results <- float64(scoreSum) / float64(len(info))
		}(task)
	}

	for range profile.Scores {
		sum += uint64(<-results)
	}

	return float64(sum) / float64(len(profile.Scores))
}

func avgScorePerTask(username string) float64 {
	// gets avg score per task

	profile, err := GetProfile(username)
	if err != nil {
		panic(err)
	}

	var sum uint64 = 0
	for _, task := range profile.Scores {
		sum += uint64(task.Score)
	}

	return float64(sum) / float64(len(profile.Scores))
}

func main() {
	x := avgScorePerTask(os.Args[1])
	fmt.Println(x)
}
