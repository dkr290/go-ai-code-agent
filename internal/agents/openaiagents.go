package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	OpenAIEndpoint = "https://api.openai.com/v1/chat/completions"
	OpenAIModel    = "gpt-4o-mini"
)

// Response that holds the result

type OpenAI struct {
	httpClient *http.Client
	ctx        context.Context
	apiKey     string
}

func NewOpenAI(ctx context.Context, apiKey string, httpClient *http.Client) *OpenAI {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 120,
		}
	}

	return &OpenAI{
		ctx:        ctx,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (o *OpenAI) Query(systemPrompt, userPrompt string) (LLMQueryResponse, error) {
	if systemPrompt == "" {
		systemPrompt = "You are helpful assistant."
	}
	payload := OpenAIRequest{
		Model: OpenAIModel,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}
	// Marshal the body
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling the payload %v", err)
	}

	req, err := http.NewRequestWithContext(o.ctx, "POST", OpenAIEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+o.apiKey)

	// Do the request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return OpenAIResponse{}, err
	}
	defer resp.Body.Close()

	// Decode the response
	var openaiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return OpenAIResponse{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if openaiResp.Error != nil {
		return OpenAIResponse{}, fmt.Errorf("API error %v", err)
	}
	if len(openaiResp.Choices) == 0 {
		return OpenAIResponse{}, errors.New("no choices returned from API")
	}

	return openaiResp, nil
}
