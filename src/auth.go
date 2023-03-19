package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type AuthToken = string

const apiUrl = "https://training.olinfo.it/api"

func Login(username, password string) (AuthToken, error) {
	type payloadShape struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Action     string `json:"action"`
		KeepSigned bool   `json:"keep_signed"`
	}

	type responseShape struct {
		Success int `json:"success"`
		// technically it has an "error" field but 1. it's optional (doesn't appear on success) and 2. is very generic ("error: login.error")
	}

	payloadData := payloadShape{username, password, "login", true}
	payloadBytes, err := json.Marshal(payloadData)
	if err != nil {
		// this should be impossible
		return "", errors.New("Login(): Error while trying to encode the payload into JSON")
	}

	response, err := http.Post(apiUrl+"/user", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", errors.New("Login(): The request caused an error; this is NOT because of the credentials")
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Login(): The response had status code %d, instead of 200", response.StatusCode)
	}

	var responseData responseShape
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return "", errors.New("Login(): Error while decoding the response; maybe it has an unexpected shape")
	}

	if responseData.Success != 1 {
		return "", errors.New("Login(): Invalid credentials")
	}

	for _, cookie := range response.Cookies() {
		if cookie.Name == "training_token" {
			return cookie.Value, nil
		}
	}

	return "", errors.New("Login(): The request was successful but the AuthToken wasn't found in the cookies; maybe its name has changed")
}
