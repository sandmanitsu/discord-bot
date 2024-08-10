package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RequestBody struct {
	Inputs string `json:"inputs"`
}

type ResponseBody struct {
	GeneratedText string `json:"generated_text"`
}

type ErrorResponse struct {
	Error         string  `json:"error"`
	EstimatedTime float64 `json:"estimated_time"`
}

func Request(messages string) string {
	requestData := RequestBody{
		Inputs: messages,
	}

	jsonData, err := json.Marshal(requestData)
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Printf("Error encoding JSON: %s", err)
		return ""
	}

	request, err := http.NewRequest("POST", ModelURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %s", err)
		return ""
	}

	request.Header.Set("Authorization", "Bearer "+Token)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	for {
		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Error sending request: %s", err)
			return ""
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading response: %s", err)
			return ""
		}

		var errorResponse ErrorResponse
		err = json.Unmarshal(body, &errorResponse)
		if err == nil && errorResponse.Error != "" {
			fmt.Printf("Model loading, wait for %f second...and try again\n", errorResponse.EstimatedTime)
			time.Sleep(time.Duration(errorResponse.EstimatedTime) * time.Second)
			continue
		}

		var responseBody []ResponseBody
		err = json.Unmarshal(body, &responseBody)
		if err != nil {
			fmt.Printf("Error decoding JSON: %s", err)
			return ""
		}

		return responseBody[0].GeneratedText
	}
}
