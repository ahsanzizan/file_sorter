package main

import (
	"file-sorter/internal/configuration"
	"file-sorter/internal/filesorter"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func NewFileSorter(configPath string) (*filesorter.FileSorter, error) {
	config, err := configuration.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create new watcher: %w", err)
	}

	var logger *log.Logger
	if config.EnableLogging {
		logFile, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		logger = log.New(logFile, "[AUTO-SORT] ", log.LstdFlags)
	} else {
		logger = log.New(os.Stdout, "[AUTO-SORT] ", log.LstdFlags)
	}

	return &filesorter.FileSorter{
		Config:  config,
		Watcher: watcher,
		Logger:  logger,
	}, nil
}

func main() {
	var (
		configPath = flag.String("config", "./config.json", "Path to configuration file")
		dryRun     = flag.Bool("dry-run", false, "Run in dry-run mode (don't actually move files)")
		generate   = flag.Bool("generate-config", false, "Generate a default configuration file")
	)
	flag.Parse()

	// Generate config and exit if requested
	if *generate {
		_, err := NewFileSorter(*configPath)
		if err != nil {
			log.Fatalf("Failed to generate config: %v", err)
		}
		fmt.Printf("Configuration file generated at: %s\n", *configPath)
		fmt.Println("Edit the configuration file to customize sorting rules, then run the program again.")
		return
	}

	// Create file sorter
	fs, err := NewFileSorter(*configPath)
	if err != nil {
		log.Fatalf("Failed to create file sorter: %v", err)
	}
	defer fs.Stop()

	// Override dry-run setting if specified via command line
	if *dryRun {
		fs.Config.DryRun = true
	}

	// Start the file sorter
	if err := fs.Start(); err != nil {
		log.Fatalf("Failed to start file sorter: %v", err)
	}

	// Keep the program running
	fmt.Println("Auto-sort is running. Press Ctrl+C to stop.")
	select {}
}
