package utils

import (
	"log"
	"regexp"
	"strings"
)

type FileParser struct {
	Path    string
	Content string
}

func ParseCode(content string) ([]FileParser, error) {
	codeBlockRegEx := regexp.MustCompile(`(?s)---FILE_PATH: (.+?)\n(.*?)---END_FILE`)
	matches := codeBlockRegEx.FindAllStringSubmatch(content, -1)

	if len(matches) == 0 {
		log.Println("No code  blocks found, could not find FILE_PATH in the content")
		return nil, nil
	}

	var fp []FileParser
	for _, match := range matches {
		if len(match) < 3 {
			log.Println("Invalid match  found, skipping")
			continue
		}
		filePath := strings.TrimSpace(match[1])
		code := strings.TrimSpace(match[2])

		code = regexp.MustCompile("^```[a-zA-Z0-9]*\n").ReplaceAllString(code, "")
		code = regexp.MustCompile("\n```$").ReplaceAllString(code, "")

		fp = append(fp, FileParser{Path: filePath, Content: code})

	}

	return fp, nil
}
