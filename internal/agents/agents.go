package agents

import (
	"context"
	"sync"
)

type FileTask struct {
	Path    string
	Content string
}

type Agent struct {
	openAI          *OpenAI
	deepSeek        *DeepSeek
	gemini          *GeminiApi
	outputDir       string
	basePackage     string
	taskQueue       chan FileTask
	wg              sync.WaitGroup
	workerCount     int
	ctx             context.Context
	cancel          context.CancelFunc
	fileWriterMutex sync.Mutex
	fileWritten     map[string]bool
}

func NewAgent(
	ctx context.Context,
	openAI *OpenAI,
	deepSeek *DeepSeek,
	gemini *GeminiApi,
	outputDir string,
	basePackage string,
	workerCount int,
) *Agent {
	return &Agent{
		ctx:         ctx,
		openAI:      openAI,
		deepSeek:    deepSeek,
		gemini:      gemini,
		outputDir:   outputDir,
		basePackage: basePackage,
		workerCount: workerCount,
		taskQueue:   make(chan FileTask, 100),
		fileWritten: make(map[string]bool),
	}
}
