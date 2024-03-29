package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	ApiKey string
	ApiURL string
}

func NewConfig(apiKey, apiURL string) Config {
	return Config{
		ApiKey: apiKey,
		ApiURL: apiURL,
	}
}

type ConversationResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

var languageExtensions = map[string]string{
	"python":     ".py",
	"go":         ".go",
	"ruby":       ".rb",
	"perl":       ".pl",
	"bash":       ".sh",
	"powershell": ".ps1",
	"javascript": ".js",
	"typescript": ".ts",
	"php":        ".php",
	"lua":        ".lua",
}

//go:embed apikey.txt
var key embed.FS

func chat(config Config, prompt string) {
	scriptLanguage := detectLanguage(prompt)
	if scriptLanguage == "" {
		fmt.Println("Unsupported language")
		return
	}
	fmt.Println("Please specify the output directory:")
	reader := bufio.NewReader(os.Stdin)
	outputDir, _ := reader.ReadString('\n')
	outputDir = strings.TrimSpace(outputDir)

	scriptFilename := fmt.Sprintf("gollum_script%s", languageExtensions[scriptLanguage])
	scriptPath := filepath.Join(outputDir, scriptFilename)
	responseText, err := generateScript(config, prompt)
	if err != nil {
		fmt.Println("Error generating script:", err)
		return
	}
	if err := saveScript(scriptPath, responseText); err != nil {
		fmt.Println("Error saving script:", err)
		return
	}

	fmt.Printf("%s script generated in %s\n", cases.Title(language.Und).String(scriptLanguage), outputDir)
}

func detectLanguage(input string) string {
	input = strings.ToLower(input)

	if strings.Contains(input, "python") {
		return "python"
	} else if strings.Contains(input, "ruby") {
		return "ruby"
	} else if strings.Contains(input, "perl") {
		return "perl"
	} else if strings.Contains(input, "bash") {
		return "bash"
	} else if strings.Contains(input, "powershell") {
		return "powershell"
	} else if strings.Contains(input, "javascript") {
		return "javascript"
	} else if strings.Contains(input, "typescript") {
		return "typescript"
	} else if strings.Contains(input, "php") {
		return "php"
	} else if strings.Contains(input, "lua") {
		return "lua"
	} else if strings.Contains(input, "go") || strings.Contains(input, "golang") {
		return "go"
	}

	return ""
}

func generateScript(config Config, prompt string) (string, error) {
	response, err := getResponse(config, prompt)
	if err != nil {
		return "", err
	}

	var code strings.Builder
	for _, choice := range response.Choices {
		code.WriteString(strings.TrimSpace(choice.Text))
		code.WriteString("\n")
	}

	return code.String(), nil
}

func saveScript(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	total := len(content)
	barWidth := 40
	barChar := "="

	writer := bufio.NewWriter(f)

	for i, char := range content {
		_, err := writer.WriteRune(char)
		if err != nil {
			return err
		}

		progress := float64(i+1) / float64(total)
		numBarChars := int(progress * float64(barWidth))
		bar := strings.Repeat(barChar, numBarChars) + strings.Repeat(" ", barWidth-numBarChars)

		print("\r[", bar, "] ", int(progress*100), "%")
	}

	writer.Flush()
	println()

	return nil
}

func getResponse(config Config, prompt string) (ConversationResponse, error) {

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":             "gpt-3.5-turbo-instruct",
		"prompt":            "```" + prompt + "\n```",
		"top_p":             1,
		"stop":              "```",
		"temperature":       0,
		"suffix":            "\n```",
		"max_tokens":        1000,
		"presence_penalty":  0,
		"frequency_penalty": 0,
	})
	if err != nil {
		return ConversationResponse{}, err
	}

	req, err := http.NewRequest("POST", config.ApiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return ConversationResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.ApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ConversationResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ConversationResponse{}, err
	}

	var conversationResponse ConversationResponse
	err = json.Unmarshal(body, &conversationResponse)
	if err != nil {
		return ConversationResponse{}, err
	}

	return conversationResponse, nil
}

func readAPIKey(filename string) (string, error) {
	content, err := key.ReadFile(filename)
	if err != nil {
		return "", err
	}
	apiKey := strings.TrimSpace(string(content))
	return apiKey, nil
}

func main() {
	apiKey, err := readAPIKey("apikey.txt")
	if err != nil {
		fmt.Println("Error reading API key:", err)
		return
	}
	apiURL := "https://api.openai.com/v1/completions"
	config := NewConfig(apiKey, apiURL)

	fmt.Println("Welcome to goLLum. Generate a script or type 'quit' to exit.")

	for {
		fmt.Print("You: ")
		reader := bufio.NewReader(os.Stdin)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if userInput == "quit" {
			break
		}

		chat(config, userInput)
	}
}
