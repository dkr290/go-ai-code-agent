package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dkr290/go-ai-code-agent/internal/agents"
)

func main() {
	okey := os.Getenv("OPENAI_API_KEY")
	if len(okey) == 0 {
		fmt.Println("NEED OPENAI_API_KEY env var")
		return
	}

	dkey := os.Getenv("DEEPSEEK_API_KEY")
	if len(dkey) == 0 {
		fmt.Println("NEED DEEPSEEK_API_KEY env var")
		return
	}
	gkey := os.Getenv("GEMINI_API_KEY")
	if len(gkey) == 0 {
		fmt.Println("NEED GEMINI_API_KEY env var")
		return
	}

	ctx := context.Background()

	prompt := "Write simnple todo program in go"

	openaiClient := agents.NewOpenAI(ctx, okey, nil)
	resp, err := openaiClient.Query("", prompt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Openai Reponce")
	fmt.Printf("\n\n%+v\n\n", resp)

	deepSeekPrompt := "Write simnple todo program in nodejs"

	deepSeekClient := agents.NewDeepSeek(ctx, dkey, nil)
	d_resp, err := deepSeekClient.Query("", deepSeekPrompt)
	if err != nil {
		panic(err)
	}
	fmt.Println("DeepSeek Response below here is there ")
	fmt.Printf("\n\n%+v\n\n", d_resp)

	geminiPrompt := "Write simple todo program in python"
	geminiClient := agents.NewGemini(ctx, gkey)
	gemini_resp, err := geminiClient.QueryGemini("", geminiPrompt)
	if err != nil {
		panic(err)
	}
	fmt.Println("Gemini Response below here is there ")
	fmt.Printf("\n\n%+v\n\n", gemini_resp)
}
