package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// curl https://api.openai.com/v1/responses \
//   -H "Content-Type: application/json" \
//   -H "Authorization: Bearer $OPENAI_API_KEY" \
//   -d '{
//     "model": "gpt-4o",
//     "input": "Tell me a three sentence bedtime story about a unicorn."
//   }'

type OpenAIInput struct {
	Model string  `json:"model"`
	Input []Input `json:"input"`
}

type Input struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type OpenAIResponse struct {
	ID        string   `json:"id"`
	Object    string   `json:"object"`
	CreatedAt int64    `json:"created_at"`
	Status    string   `json:"status"`
	Error     string   `json:"error"`
	Output    []Output `json:"output"`
}

type Output struct {
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	Status  string    `json:"status"`
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func MakeOpenAIRequest(ctx context.Context, input OpenAIInput) (*OpenAIResponse, error) {

	// Set up the request
	url := "https://api.openai.com/v1/responses"
	method := "POST"

	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	openAIKey := os.Getenv("OPENAIKEY")

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got a non 200 response: %s", resp.Status)
	}

	// Parse the response
	var response OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	// Check for errors in the response
	if response.Error != "" {
		return nil, fmt.Errorf("Error from OpenAI: %s", response.Error)
	}

	// Return the response
	return &response, nil

}
