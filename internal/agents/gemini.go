package agents

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

type GeminiApi struct {
	ctx    context.Context
	apiKey string
	model  string
}

func NewGemini(ctx context.Context, apiKey string, model string) *GeminiApi {
	return &GeminiApi{
		ctx:    ctx,
		apiKey: apiKey,
		model:  model,
	}
}

func (g *GeminiApi) QueryGemini(systemPrompt, userPrompt string) (string, error) {
	if systemPrompt == "" {
		systemPrompt = "You are helpful assistant."
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  g.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("error creating new genai agent %v", err)
	}

	result, err := client.Models.GenerateContent(g.ctx, g.model,
		genai.Text(userPrompt),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: systemPrompt}},
			},
			Tools: []*genai.Tool{{
				GoogleSearch: &genai.GoogleSearch{},
			}},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error generate context %v", err)
	}

	// Decode the response
	return result.Text(), nil
}
