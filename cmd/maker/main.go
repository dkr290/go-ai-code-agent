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
	addTemplate *string
	userPrompt  *string
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
	addTemplate = flag.String(
		"use-template",
		"",
		"supported template for additional instructions - java-spring, go-gin",
	)
	userPrompt = flag.String(
		"user-prompt",
		"sample todo app",
		"User prompt for the application to create",
	)

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
		runAgent(a, nil)

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
		runAgent(a, nil)

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
		a := agents.NewAgent(ctx, nil, nil, geminiClient, *outputDir, *basePackage)

		templ, err := a.LoadTemplatesFromFolder()
		if err != nil {
			log.Fatal(err)
		}
		var p string
		for _, t := range templ {
			if *language == t.Language && *addTemplate == t.Name {
				p = t.Prompt
			}
		}

		prompt, err := utils.GetSystemPrompt(
			*language,
			*basePackage,
			p,
		)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := geminiClient.QueryGemini(prompt, *userPrompt)
		if err != nil {
			log.Fatal(err)
		}
		fileparse, err := utils.ParseCode(resp)
		if err != nil {
			log.Fatal("error parsing code", err)
		}

		runAgent(a, fileparse)

	default:
		return errors.New("wrong option only deepseek, openai and gemini accepted")
	}

	return nil
}

func runAgent(a *agents.Agent, fileparse []utils.FileParser) {
	go func() {
		for err := range a.GetErrorChan() {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	for _, file := range fileparse {
		a.WriteFile(
			file.Path,
			file.Content,
		)
	}
	a.Close()
	log.Println("Finished writing to the output directory", *outputDir)
}
