package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/ln64/voxcptr/internal/config"
	"github.com/ln64/voxcptr/internal/vosk"
)

func main() {
	startFlag := flag.Bool("start", false, "Start speech recognition")
	// modelPath := flag.String("model", "model", "Path to Vosk model")
	flag.Parse()

	if !*startFlag {
		fmt.Println("Usage: go run main.go -start -model=path/to/model")
		os.Exit(1)
	}

	// Retrieve configuration
	configName := "voxcptr.json"
	configData, err := config.GetConfig(configName)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	modelPath := config.GetStringOrDefault(configData, "VoskModelPath", "")

	log.Info(modelPath)

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
			result = strings.TrimSpace(result)
			if result != "" {
				err := vosk.CopyToClipboard(result)
				if err != nil {
					log.Printf("Failed to copy to clipboard: %v", err)
				} else {
					log.Print("Copied to clipboard:", result)
				}
			}
		}
	}()

	err = speechRecognizer.Start(resultChan)
	if err != nil {
		log.Fatalf("Error starting speech recognizer: %v", err)
	}

	fmt.Println("Speech recognition started. Press Enter to stop...")
	fmt.Scanln()

	speechRecognizer.Stop()
}
