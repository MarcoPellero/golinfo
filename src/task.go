package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
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

type Submission struct {
	SubmissionId       int
	TaskId             int
	SourceCodeId       string
	Timestamp          time.Time
	CompilationSuccess bool
	EvaluationSuccess  bool
	Score              int
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

func GetTaskSubmissions(taskName string, token AuthToken) ([]Submission, error) {
	type payloadShape struct {
		TaskName string `json:"task_name"`
		Action   string `json:"action"`
	}

	type apiFileShape struct {
		Name       string `json:"name"`
		HashDigest string `json:"digest"`
	}

	type apiSubmissionShape struct {
		Id                 int            `json:"id"`
		TaskId             int            `json:"task_id"`
		Timestamp          float64        `json:"timestamp"`
		CompilationOutcome string         `json:"compilation_outcome"`
		EvaluationOutcome  string         `json:"evaluation_outcome"`
		Score              float32        `json:"score"`
		Files              []apiFileShape `json:"files"`
	}

	type responseShape struct {
		Submissions []apiSubmissionShape `json:"submissions"`
		Success     int                  `json:"success"`
		// technically it has an "error" field but 1. it's optional (doesn't appear on success) and 2. is very generic ("error: not found")
	}

	payloadData := payloadShape{taskName, "list"}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		// this should be impossible
		return []Submission{}, errors.New("GetTaskSubmissions(): Error while trying to encode the payload into JSON")
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl+"/submission", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return []Submission{}, errors.New("GetTaskSubmissions(): Error while crafting the request")
	}

	request.Header.Add("content-type", "application/json")
	request.AddCookie(&http.Cookie{Name: "training_token", Value: token})
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []Submission{}, errors.New("GetTaskSubmissions(): The request caused an error")
	}

	if response.StatusCode != 200 {
		return []Submission{}, fmt.Errorf("GetTaskSubmissions(): The response had status code %d, instead of 200", response.StatusCode)
	}

	var responseData responseShape
	err = json.NewDecoder(response.Body).Decode(&responseData)

	if err != nil {
		return []Submission{}, errors.New("GetTaskSubmissions(): Error while decoding the response; maybe it has an unexpected shape")
	}

	if responseData.Success != 1 {
		return []Submission{}, fmt.Errorf("GetTaskSubmissions(): Invalid task name %s", taskName)
	}

	parsedData := make([]Submission, len(responseData.Submissions))
	for i, apiSub := range responseData.Submissions {
		unixSeconds := int64(math.Round(apiSub.Timestamp))
		tmp := apiSub.Timestamp - float64(unixSeconds)
		for tmp-math.Round(tmp) > 0 {
			tmp *= 10
		}
		unixNanoseconds := int64(tmp)

		parsedData[i] = Submission{
			apiSub.Id,
			apiSub.TaskId,
			apiSub.Files[0].HashDigest,
			time.Unix(unixSeconds, unixNanoseconds),
			apiSub.CompilationOutcome == "ok",
			apiSub.EvaluationOutcome == "ok",
			int(apiSub.Score),
		}
	}

	return parsedData, nil
}
