package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/nlatham1999/aiapp/openai"
)

//go:embed index.html
var homePageHtml string

type Server struct {
	// ipMessages map[string]
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Perform health check logic
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func (s *Server) PromptHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the prompt request
	type Prompt struct {
		Prompt string `json:"prompt"`
	}
	var prompt Prompt
	if err := json.NewDecoder(r.Body).Decode(&prompt); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	input := openai.OpenAIInput{
		Model: "gpt-4o-mini",
		Input: []openai.Input{
			{
				Role: "system",
				Content: []openai.Content{
					{
						Type: "input_text",
						Text: "you are a soothsayer. Your answers should a yes or no, not a maybe, but the answers should be warpped in your style so the exact words yes or no whould never be used. You should be witty and playful yet a little mysterious. You do not require all the context to respond since you know all. You should be somewhat positive but not all the time.",
					},
				},
			},
			{
				Role: "user",
				Content: []openai.Content{
					{
						Type: "input_text",
						Text: prompt.Prompt,
					},
				},
			},
		},
	}

	response, err := openai.MakeOpenAIRequest(context.Background(), input)
	if err != nil {
		http.Error(w, "Error processing request "+err.Error(), http.StatusInternalServerError)
		return
	}

	outputs := response.Output

	if len(outputs) == 0 {
		http.Error(w, "No output from OpenAI", http.StatusInternalServerError)
		return
	}

	// Assuming the first output is the one we want to return
	output := outputs[0]

	content := output.Content

	if len(content) == 0 {
		http.Error(w, "No content in output", http.StatusInternalServerError)
		return
	}

	// Assuming the first content is the one we want to return
	contentText := content[0].Text
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Create a JSON response
	responseData := map[string]string{
		"response": contentText,
	}
	jsonResponse, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {

	htmlTmpl, err := template.New("content").Parse(homePageHtml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	htmlTmpl.Execute(w, nil)

}
