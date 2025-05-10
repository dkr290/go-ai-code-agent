package agents

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

//go:embed templates/*
var templateFiles embed.FS

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

func (a *Agent) Close() {
	a.wg.Wait()
	close(a.errorChan)
}

func (a *Agent) GetErrorChan() chan error {
	return a.errorChan
}

func (a *Agent) LoadTemplatesFromFolder() (TemplateConfigs, error) {
	templates := make(TemplateConfigs)
	entries, err := templateFiles.ReadDir(
		"templates",
	) // Read the contents of the embedded "templates" directory
	if err != nil {
		return nil, fmt.Errorf("error reading embedded templates directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			filePath := filepath.Join(
				"templates",
				entry.Name(),
			) // Construct the full path within the embedded FS
			content, err := templateFiles.ReadFile(filePath)
			if err != nil {
				fmt.Printf("error reading embedded file '%s': %v\n", filePath, err)
				continue // Continue to the next file
			}

			var config TemplateConfig
			err = json.Unmarshal(content, &config)
			if err != nil {
				fmt.Printf("error unmarshalling JSON from embedded file '%s': %v\n", filePath, err)
				continue // Continue to the next file
			}

			nameWithoutExt := entry.Name()[:len(entry.Name())-len(filepath.Ext(entry.Name()))]
			templates[nameWithoutExt] = config
		}
	}

	return templates, nil
}
