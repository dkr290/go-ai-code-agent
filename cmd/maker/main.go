package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/dkr290/go-ai-code-agent/internal/agents"
)

func main() {
	openaiKey := flag.String("openai-key", "", "OpenAi API key")
	deepSeekKey := flag.String("deepseek-key", "", "DeepSeek API key")
	geminiKey := flag.String("gemini-key", "", "Gemini API key")
	useLLM := flag.String("use-llm", "openai", "LLM to use (openai, deepseek, gemini)")
	outputDir := flag.String("output-dir", "./output", "Output directory for generated files")
	basePackage := flag.String(
		"base-package",
		"github.com/user/package",
		"Base package for generated files",
	)
	workerCount := flag.Int("worker-count", 4, "Number of workers to use for file generation")
	flag.Parse()

	if useLLM == nil {
		fmt.Println("NEED use-llm flag")
		return
	}

	ctx := context.Background()

	if *useLLM == "openai" {
		if *openaiKey == "" {
			*openaiKey = os.Getenv("OPENAI_API_KEY")
			if *openaiKey == "" {
				fmt.Println(
					"NEED OPENAI_API_KEY env var or to be passed as command line arg -openai-key",
				)
				return
			}
		}

		openaiClient := agents.NewOpenAI(ctx, *openaiKey, nil)
		a := agents.NewAgent(ctx, openaiClient, nil, nil, *outputDir, *basePackage, *workerCount)
		a.Start()

	}
	if *useLLM == "deepseek" {
		if *deepSeekKey == "" {
			*deepSeekKey = os.Getenv("DEEPSEEK_API_KEY")
			if *deepSeekKey == "" {
				fmt.Println(
					"NEED DEEPSEEK_API_KEY env var or to be passed as command line arg -deepseek-key",
				)
				return
			}
		}
		deepSeekClient := agents.NewDeepSeek(ctx, *deepSeekKey, nil)
		a := agents.NewAgent(ctx, nil, deepSeekClient, nil, *outputDir, *basePackage, *workerCount)
		a.Start()
	}
	if *useLLM == "gemini" {
		if *geminiKey == "" {
			*geminiKey = os.Getenv("GEMINI_API_KEY")
			if *geminiKey == "" {
				fmt.Println(
					"NEED GEMINI_API_KEY env var or to be passed as command line arg -gemini-key",
				)
				return
			}
		}
		geminiClient := agents.NewGemini(ctx, *geminiKey)
		agents.NewAgent(ctx, nil, nil, geminiClient, *outputDir, *basePackage, *workerCount)
	}
}
