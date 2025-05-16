package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dkr290/go-ai-code-agent/internal/agents"
	"github.com/dkr290/go-ai-code-agent/internal/utils"
)

func run(ctx context.Context, params *Params) error {
	switch params.useLLM {
	case "deepseek":
		if params.deepSeekKey == "" {
			params.deepSeekKey = os.Getenv("DEEPSEEK_API_KEY")
			if params.deepSeekKey == "" {
				return errors.New(
					"NEED DEEPSEEK_API_KEY env var or to be passed as command line arg -deepseek-key",
				)
			}
		}

		deepSeekClient := agents.NewDeepSeek(ctx, params.deepSeekKey, nil, params.modelName)
		a := agents.NewAgent(ctx, nil, deepSeekClient, nil, params.outputDir, params.basePackage)
		prompt, err := getPrompt(a, params)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := deepSeekClient.Query(prompt, params.userPrompt)
		if err != nil {
			log.Fatal(err)
		}
		fileParse, err := utils.ParseCode(resp.Choices[0].Message.Content)
		if err != nil {
			log.Fatal("error parsing code", err)
		}
		runAgent(a, fileParse, params)

	case "openai":
		if params.openaiKey == "" {
			params.openaiKey = os.Getenv("OPENAI_API_KEY")
			if params.openaiKey == "" {
				return errors.New(
					"NEED OPENAI_API_KEY env var or to be passed as command line arg -openai-key",
				)
			}
		}

		openaiClient := agents.NewOpenAI(ctx, params.openaiKey, nil, params.modelName)
		a := agents.NewAgent(ctx, openaiClient, nil, nil, params.outputDir, params.basePackage)
		prompt, err := getPrompt(a, params)
		if err != nil {
			log.Fatal(err)
		}
		resp, err := openaiClient.Query(prompt, params.userPrompt)
		if err != nil {
			log.Fatal(err)
		}
		fileParse, err := utils.ParseCode(resp.Choices[0].Message.Content)
		if err != nil {
			log.Fatal("error parsing code", err)
		}
		runAgent(a, fileParse, params)

	case "gemini":
		if params.geminiKey == "" {
			params.geminiKey = os.Getenv("GEMINI_API_KEY")
			if params.geminiKey == "" {
				return errors.New(
					"NEED GEMINI_API_KEY env var or to be passed as command line arg -gemini-key",
				)
			}
		}
		geminiClient := agents.NewGemini(ctx, params.geminiKey, params.modelName)
		a := agents.NewAgent(ctx, nil, nil, geminiClient, params.outputDir, params.basePackage)

		prompt, err := getPrompt(a, params)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := geminiClient.QueryGemini(prompt, params.userPrompt)
		if err != nil {
			log.Fatal(err)
		}
		fileparse, err := utils.ParseCode(resp)
		if err != nil {
			log.Fatal("error parsing code", err)
		}

		runAgent(a, fileparse, params)

	default:
		return errors.New("wrong option only deepseek, openai and gemini accepted")
	}

	return nil
}

func runAgent(a *agents.Agent, fileparse []utils.FileParser, params *Params) {
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
	log.Println("Finished writing to the output directory", params.outputDir)
}

func getPrompt(a *agents.Agent, params *Params) (string, error) {
	templ, err := a.LoadTemplatesFromFolder()
	if err != nil {
		return "", err
	}
	var p string
	for _, t := range templ {
		if params.language == t.Language && params.addTemplate == t.Name {
			p = t.Prompt
		}
	}

	prompt, err := utils.GetSystemPrompt(
		params.language,
		params.basePackage,
		p,
	)
	if err != nil {
		return "", err
	}

	return prompt, nil
}
