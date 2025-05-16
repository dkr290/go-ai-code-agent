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
	DeepSeekEndpoint = "https://api.deepseek.com/chat/completions"
)

type DeepSeek struct {
	httpClient *http.Client
	ctx        context.Context
	apiKey     string
	modelName  string
}

func NewDeepSeek(
	ctx context.Context,
	apiKey string,
	httpClient *http.Client,
	modelName string,
) *DeepSeek {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 240,
		}
	}

	return &DeepSeek{
		ctx:        ctx,
		apiKey:     apiKey,
		httpClient: httpClient,
		modelName:  modelName,
	}
}

func (o *DeepSeek) Query(systemPrompt, userPrompt string) (DeepSeekResponse, error) {
	if systemPrompt == "" {
		systemPrompt = "You are helpful assistant."
	}

	payload := DeepSeekPayload{
		Messages: []DeepSeekMessage{
			{
				Content: systemPrompt,
				Role:    "system",
			},
			{
				Content: userPrompt,
				Role:    "user",
			},
		},
		Model: o.modelName,
	}
	// Marshal the body
	body, err := json.Marshal(payload)
	if err != nil {
		return DeepSeekResponse{}, fmt.Errorf("error marshalling the payload %v", err)
	}

	req, err := http.NewRequestWithContext(o.ctx, "POST", DeepSeekEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return DeepSeekResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+o.apiKey)

	// Do the request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return DeepSeekResponse{}, err
	}
	defer resp.Body.Close()

	// Decode the response
	var DeepSeekResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&DeepSeekResp); err != nil {
		return DeepSeekResponse{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if DeepSeekResp.Error != nil {
		return DeepSeekResponse{}, fmt.Errorf("API error %v", err)
	}
	if len(DeepSeekResp.Choices) == 0 {
		return DeepSeekResponse{}, errors.New("no choices returned from API")
	}
	return DeepSeekResp, nil
}
