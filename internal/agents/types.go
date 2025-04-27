package agents

// openAI response struct
// The JSON is deeply nested (choices[0].message.content)
// lightweight struct that only digs into the parts you care about.
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// the openAI request
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// types concerning the deepseek

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// deepseek message
type DeepSeekMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type DeepSeekPayload struct {
	Messages []DeepSeekMessage `json:"messages"`
	Model    string            `json:"model"`
}
