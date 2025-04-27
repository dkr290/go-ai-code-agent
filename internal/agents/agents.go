package agents

type FileTask struct {
	Path    string
	Content string
}

type Agent struct {
	openAI      *OpenAI
	deepSeek    *DeepSeek
	gemini      *GeminiApi
	outputDir   string
	basePackage string
	taskQueue   chan FileTask
}
