package agents

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func (a *Agent) Start() {
	fmt.Printf("Starting agent with %d workers\n", a.workerCount)

	for i := 0; i < a.workerCount; i++ {
		a.wg.Add(1)
		go a.worker(i)
	}
}

func (a *Agent) worker(workerID int) {
	defer a.wg.Done()
	fmt.Printf("Worker %d started\n", workerID)
	for {
		select {
		case task, ok := <-a.taskQueue:
			if !ok {
				fmt.Printf("Worker %d exiting\n", workerID)
				return
			}
			a.fileWriterMutex.Lock()
			if a.fileWritten[task.Path] {
				a.fileWriterMutex.Unlock()
				fmt.Printf("Worker %d skipping already written file: %s\n", workerID, task.Path)
				continue
			}
			a.fileWritten[task.Path] = true
			a.fileWriterMutex.Unlock()
			err := a.writeFile(task)
			if err != nil {
				fmt.Printf("Worker %d error writing file %s: %v\n", workerID, task.Path, err)
			} else {
				fmt.Printf("Worker %d successfully wrote file: %s\n", workerID, task.Path)
			}
		case <-a.ctx.Done():
			fmt.Printf("Worker %d exiting due to context cancellation\n", workerID)
			return
		}
	}
}

func (a *Agent) writeFile(task FileTask) error {
	fullPath := filepath.Join(a.outputDir, task.Path)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	err := os.WriteFile(fullPath, []byte(task.Content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}
	log.Printf("Successfully wrote file: %s\n", fullPath)

	return nil
}

func (a *Agent) SendFileTask(path string, content string) {
	task := FileTask{
		Path:    path,
		Content: content,
	}
	go func() {
		a.taskQueue <- task
	}()
}
