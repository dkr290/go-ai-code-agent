package utils

import "fmt"

type PromptTemplate struct {
	Language string
	Template string
}

var defaultPrompts = []PromptTemplate{
	{
		Language: "go",
		Template: `You are a code generation assistant. Provide complete, accurate, and well-structured Go code based on the user's requirements.
Format your response like this for each file:

---FILE_PATH: path/to/filename.ext
[code content goes here]
---END_FILE

Make sure to include ALL necessary files to make the application work, including configuration files, main files, package files, etc.
Always include a README.md with instructions on how to run the application.

IMPORTANT:
1. DO NOT include markdown code block markers (like "` + "```" + `go" or "` + "```" + `") in your code content.
2. Use proper Go package structure based on this base package: %s
3. Include a Makefile with common commands (build, run, test, etc.)
4. Ensure all Go files have the correct package declarations based on their directory structure.
%s`,
	},
	{
		Language: "python",
		Template: `You are a code generation assistant. Provide complete, accurate, and well-structured python code based on the user's requirements.
Format your response like this for each file:

---FILE_PATH: path/to/filename.ext
[code content goes here]
---END_FILE

Make sure to include ALL necessary files to make the application work, including configuration files, main files, package files, etc.
Always include a README.md with instructions on how to run the application.

IMPORTANT:
1. DO NOT include markdown code block markers (like "` + "```" + `python" or "` + "```" + `") in your code content.
2. Use proper python package structure based on this base package: %s 
3. Include a Makefile with common commands (build, run, test, etc.)
4. Ensure all python files have the correct package declarations based on their directory structure.
%s`,
	},
	{
		Language: "java",
		Template: `You are a code generation assistant. Provide complete, accurate, and well-structured java code based on the user's requirements.
Format your response like this for each file:

---FILE_PATH: path/to/filename.ext
[code content goes here]
---END_FILE

Make sure to include ALL necessary files to make the application work, including configuration files, main files, package files, etc.
Always include a README.md with instructions on how to run the application.

IMPORTANT:
1. DO NOT include markdown code block markers (like "` + "```" + `java" or "` + "```" + `") in your code content.
2. Use proper java package structure based on this base package: %s
3. Include a Makefile with common commands (build, run, test, etc.)
4. Ensure all java files have the correct package declarations based on their directory structure.
%s`,
	},
}

func GetSystemPrompt(lang string, basepackage string, extraPrompt string) string {
	for _, val := range defaultPrompts {
		if val.Language == lang {
			switch lang {
			case "go":
				return fmt.Sprintf(val.Template, basepackage, extraPrompt)
			case "python":
				return fmt.Sprintf(val.Template, basepackage, extraPrompt)

			case "java":
				return fmt.Sprintf(val.Template, basepackage, extraPrompt)
			default:
				return "unknown language"
			}
		}
	}
	return "unknown language"
}
