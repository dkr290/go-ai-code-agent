package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Agent struct {
	openAI       *OpenAI
	deepSeek     *DeepSeek
	gemini       *GeminiApi
	outputDir    string
	basePackage  string
	ctx          context.Context
	writtenFiles map[string]bool
	mu           sync.Mutex
	wg           sync.WaitGroup
	errorChan    chan error
}

func NewAgent(
	ctx context.Context,
	openAI *OpenAI,
	deepSeek *DeepSeek,
	gemini *GeminiApi,
	outputDir string,
	basePackage string,
) *Agent {
	return &Agent{
		ctx:          ctx,
		openAI:       openAI,
		deepSeek:     deepSeek,
		gemini:       gemini,
		outputDir:    outputDir,
		basePackage:  basePackage,
		writtenFiles: make(map[string]bool),
		errorChan:    make(chan error),
	}
}

func (a *Agent) WriteFile(path string, content string) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		a.mu.Lock()

		if a.writtenFiles[path] {
			a.mu.Unlock()
			return
		}
		a.writtenFiles[path] = true
		a.mu.Unlock()

		fullPath := filepath.Join(a.outputDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			a.errorChan <- fmt.Errorf("failed to create directory %s: %w", fullPath, err)
			return
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			a.errorChan <- fmt.Errorf("failed to write file %s: %w", fullPath, err)
			return
		}
	}()
}

func (a *Agent) Wait() {
	a.wg.Wait()
	close(a.errorChan)
}

func (a *Agent) GetErrorChan() chan error {
	return a.errorChan
}
