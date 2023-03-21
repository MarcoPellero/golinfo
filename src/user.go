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

type ProfileScore struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Score int    `json:"score"`
}

type Profile struct {
	AccessLevel       int    `json:"access_level"`
	GlobalAccessLevel int    `json:"global_access_level"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Username          string `json:"username"`
	JoinDate          time.Time
	MailHash          string         `json:"mail_hash"`
	TotalScore        int            `json:"score"`
	Scores            []ProfileScore `json:"scores"`
}

func GetProfile(username string) (Profile, error) {
	type payloadShape struct {
		Username string `json:"username"`
		Action   string `json:"action"`
	}

	type responseShape struct {
		Profile
		Institute     struct{} `json:"institute"` // i have no clue
		Success       int      `json:"success"`
		JoinDateEpoch float64  `json:"join_date"`
	}

	payloadData := payloadShape{username, "get"}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		// this should be impossible
		return Profile{}, errors.New("GetProfile(): Error while trying to encode the payload into JSON")
	}

	response, err := http.Post(apiUrl+"/user", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return Profile{}, errors.New("GetProfile(): The request caused an error; this is NOT because of the username")
	}

	if response.StatusCode != 200 {
		return Profile{}, fmt.Errorf("GetProfile(): The response had status code %d, instead of 200", response.StatusCode)
	}

	var responseData responseShape
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return Profile{}, errors.New("GetProfile(): Error while decoding the response; maybe it has an unexpected shape")
	}

	if responseData.Success != 1 {
		return Profile{}, errors.New("GetProfile(): Invalid username")
	}

	unixSec, unixNanosec := math.Modf(responseData.JoinDateEpoch)
	responseData.Profile.JoinDate = time.Unix(int64(unixSec), int64(unixNanosec*1e9))

	return responseData.Profile, nil
}
