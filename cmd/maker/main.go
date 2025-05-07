package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dkr290/go-ai-code-agent/internal/agents"
	"github.com/dkr290/go-ai-code-agent/internal/utils"
)

var (
	openaiKey   *string
	deepSeekKey *string
	geminiKey   *string
	useLLM      *string
	outputDir   *string
	basePackage *string
	language    *string
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
	language = flag.String("use-language", "go", " use supported language like go,python or java")
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
		a := agents.NewAgent(ctx, nil, deepSeekClient, nil, *outputDir, *basePackage)
		runAgent(a)

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
		a := agents.NewAgent(ctx, openaiClient, nil, nil, *outputDir, *basePackage)
		runAgent(a)

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
		prompt, err := utils.GetSystemPrompt(*language, *basePackage, "")
		if err != nil {
			log.Fatal(err)
		}

		resp, err := geminiClient.QueryGemini(prompt, "create todo app")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("\n\n%+v\n", resp)

		// a := agents.NewAgent(ctx, nil, nil, geminiClient, *outputDir, *basePackage)
		// runAgent(a)

	default:
		return errors.New("wrong option only deepseek, openai and gemini accepted")
	}

	return nil
}

func runAgent(a *agents.Agent) {
	a.WriteFile(
		"main.go",
		"package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println(\"Hello World.\")\n}\n",
	)
	go func() {
		for err := range a.GetErrorChan() {
			fmt.Printf("Error: %v\n", err)
		}
	}()
	a.Wait()
	log.Println("Finished writing to the output directory", *outputDir)
}
