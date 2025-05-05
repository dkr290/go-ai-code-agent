package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dkr290/go-ai-code-agent/internal/agents"
)

var (
	openaiKey   *string
	deepSeekKey *string
	geminiKey   *string
	useLLM      *string
	outputDir   *string
	basePackage *string
	workerCount *int
)

func main() {
	openaiKey = flag.String("openai-key", "", "OpenAi API key")
	deepSeekKey = flag.String("deepseek-key", "", "DeepSeek API key")
	geminiKey = flag.String("gemini-key", "", "Gemini API key")
	useLLM = flag.String("use-llm", "openai", "LLM to use (openai, deepseek, gemini)")
	outputDir = flag.String("output-dir", "./output", "Output directory for generated files")
	basePackage = flag.String(
		"base-package",
		"github.com/user/package",
		"Base package for generated files",
	)
	workerCount = flag.Int("worker-count", 4, "Number of workers to use for file generation")
	flag.Parse()

	if useLLM == nil {
		fmt.Println("NEED use-llm flag")
		return
	}

	ctx := context.Background()

	err := run(ctx, *useLLM)
	if err != nil {
		log.Println(err)
		return

	}
}

func run(ctx context.Context, isType string) error {
	switch isType {
	case "deepseek":
		if *deepSeekKey == "" {
			*deepSeekKey = os.Getenv("DEEPSEEK_API_KEY")
			if *deepSeekKey == "" {
				return errors.New(
					"NEED DEEPSEEK_API_KEY env var or to be passed as command line arg -deepseek-key",
				)
			}
		}
		deepSeekClient := agents.NewDeepSeek(ctx, *deepSeekKey, nil)
		a := agents.NewAgent(ctx, nil, deepSeekClient, nil, *outputDir, *basePackage, *workerCount)
		a.Start()

	case "openai":
		if *openaiKey == "" {
			*openaiKey = os.Getenv("OPENAI_API_KEY")
			if *openaiKey == "" {
				return errors.New(
					"NEED OPENAI_API_KEY env var or to be passed as command line arg -openai-key",
				)
			}
		}

		openaiClient := agents.NewOpenAI(ctx, *openaiKey, nil)
		a := agents.NewAgent(ctx, openaiClient, nil, nil, *outputDir, *basePackage, *workerCount)
		a.Start()

	case "gemini":
		if *geminiKey == "" {
			*geminiKey = os.Getenv("GEMINI_API_KEY")
			if *geminiKey == "" {
				return errors.New(
					"NEED GEMINI_API_KEY env var or to be passed as command line arg -gemini-key",
				)
			}
		}
		geminiClient := agents.NewGemini(ctx, *geminiKey)
		a := agents.NewAgent(ctx, nil, nil, geminiClient, *outputDir, *basePackage, *workerCount)
		a.Start()

	default:
		return errors.New("wrong option only deepseek, openai and gemini accepted")
	}

	return nil
}
