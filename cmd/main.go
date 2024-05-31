package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/ln64/voxcptr/internal/clipboard"
	"github.com/ln64/voxcptr/internal/config"
	"github.com/ln64/voxcptr/internal/vosk"
)

type TextResult struct {
	Text string `json:"text"`
}

func main() {
	// Retrieve configuration
	configName := "voxcptr.json"
	configData, err := config.GetConfig(configName)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	modelPath := config.GetStringOrDefault(configData, "VoskModelPath", "")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Fatalf("Model path does not exist: %s", modelPath)
	}
	speechRecognizer, err := vosk.NewSpeechRecognizer(modelPath)
	if err != nil {
		log.Fatalf("Error initializing speech recognizer: %v", err)
	}
	resultChan := make(chan string)
	go func() {
		for result := range resultChan {
			var textResult TextResult
			err := json.Unmarshal([]byte(result), &textResult)
			if err != nil {
				log.Printf("Failed to parse JSON: %v", err)
				continue
			}
			text := strings.TrimSpace(textResult.Text)
			if text != "" {
				err := clipboard.CopyToClipboard(text)
				if err != nil {
					log.Printf("Failed to copy to clipboard: %v", err)
				} else {
					log.Print(text)
				}
			}
		}
	}()
	err = speechRecognizer.Start(resultChan)
	if err != nil {
		log.Fatalf("Error starting speech recognizer: %v", err)
	}
	speechRecognizer.Stop()
}
