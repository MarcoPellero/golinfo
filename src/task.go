package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type leaderboardElement struct {
	// the JSON definitions are so that i don't have to redefine the same struct in GetTaskStats()
	Username    string  `json:"username"`
	TimeSeconds float32 `json:"time"`
}

type TaskStats struct {
	TotalSubmissions int
	GoodSubmissions  int
	BadSubmissions   int
	TotalUsers       int
	GoodUsers        int
	BadUsers         int
	Leaderboard      []leaderboardElement
}

func GetTaskStats(taskName string) (TaskStats, error) {
	type payloadShape struct {
		TaskName string `json:"name"`
		Action   string `json:"action"`
	}

	type responseShape struct {
		TotalSubmissions int                  `json:"nsubs"`
		GoodSubmissions  int                  `json:"nsubscorrect"`
		Success          int                  `json:"success"`
		TotalUsers       int                  `json:"nusers"`
		GoodUsers        int                  `json:"nuserscorrect"`
		Leaderboard      []leaderboardElement `json:"best"`
		// technically it has an "error" field but 1. it's optional (doesn't appear on success) and 2. is very generic ("error: not found")
	}

	payloadData := payloadShape{taskName, "stats"}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		// this should be impossible
		return TaskStats{}, errors.New("GetTaskStats(): Error while trying to encode the payload into JSON")
	}

	response, err := http.Post(apiUrl+"/task", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return TaskStats{}, errors.New("GetTaskStats(): The request caused an error")
	}

	if response.StatusCode != 200 {
		return TaskStats{}, fmt.Errorf("GetTaskStats(): The response had status code %d, instead of 200", response.StatusCode)
	}

	var responseData responseShape
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return TaskStats{}, errors.New("GetTaskStats(): Error while decoding the response; maybe it has an unexpected shape")
	}

	if responseData.Success != 1 {
		return TaskStats{}, fmt.Errorf("GetTaskStats(): Invalid task name %s", taskName)
	}

	parsedData := TaskStats{
		responseData.TotalSubmissions,
		responseData.GoodSubmissions,
		responseData.TotalSubmissions - responseData.GoodSubmissions,
		responseData.TotalUsers,
		responseData.GoodUsers,
		responseData.TotalUsers - responseData.GoodUsers,
		responseData.Leaderboard,
	}

	return parsedData, nil
}
