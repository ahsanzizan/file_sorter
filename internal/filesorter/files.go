package filesorter

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (fs *FileSorter) Start() error {
	// Create sort folders if they don't exist
	if err := fs.createSortFolders(); err != nil {
		return fmt.Errorf("failed to create sort folders: %w", err)
	}

	// Add the watch folder to the watcher
	err := fs.Watcher.Add(fs.Config.WatchFolder)
	if err != nil {
		return fmt.Errorf("failed to add watch folder: %w", err)
	}

	fs.Logger.Printf("Started monitoring: %s", fs.Config.WatchFolder)
	if fs.Config.DryRun {
		fs.Logger.Println("Running in DRY RUN mode - no files will be moved")
	}

	// Process existing files on startup
	fs.processExistingFiles()

	// Start the event loop
	go fs.watchEvents()

	return nil
}

func (fs *FileSorter) Stop() error {
	return fs.Watcher.Close()
}

func (fs *FileSorter) watchEvents() {
	for event := range fs.Watcher.Events {
		if event.Op&fsnotify.Create != 0 {
			go func(filename string) {
				// Wait briefly to ensure file is not being written
				time.Sleep(100 * time.Millisecond)
				if err := fs.processFile(filename); err != nil {
					log.Printf("Failed to process %s: %v", filename, err)
				}
			}(event.Name)
		}
	}
}

func (fs *FileSorter) processExistingFiles() {
	entries, err := os.ReadDir(fs.Config.WatchFolder)
	if err != nil {
		fs.Logger.Printf("Error reading watch folder: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fullPath := filepath.Join(fs.Config.WatchFolder, entry.Name())
			fs.processFile(fullPath)
		}
	}
}

func (fs *FileSorter) processFile(filePath string) error {
	if fs.shouldIgnoreFile(filePath) {
		return fmt.Errorf("%s should be ignored based on the config", filePath)
	}

	// Check if file is still being written (size check)
	if !fs.isFileComplete(filePath) {
		fs.Logger.Printf("File still being written, skipping: %s", filePath)
		return fmt.Errorf("%s is still being written", filePath)
	}

	targetFolder := fs.determineTargetFolder(filePath)
	if targetFolder == "" {
		fs.Logger.Printf("No matching rule for file: %s", filePath)
		return fmt.Errorf("no matching rule for file: %s", filePath)
	}

	// Move the file
	if err := fs.moveFile(filePath, targetFolder); err != nil {
		fs.Logger.Printf("Error moving file %s: %v", filePath, err)
	}

	return nil
}

func (fs *FileSorter) shouldIgnoreFile(filePath string) bool {
	fileName := filepath.Base(filePath)

	for _, pattern := range fs.Config.IgnorePatterns {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}

	// Ignore temporary files
	if strings.HasPrefix(fileName, ".") || strings.HasSuffix(fileName, ".tmp") {
		return true
	}

	return false
}

// isFileComplete checks if a file is complete by checking its size stability
func (fs *FileSorter) isFileComplete(filePath string) bool {
	info1, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	time.Sleep(50 * time.Millisecond)

	info2, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return info1.Size() == info2.Size()
}

// determineTargetFolder determines which folder a file should be moved to
func (fs *FileSorter) determineTargetFolder(filePath string) string {
	fileName := filepath.Base(filePath)
	fileExt := strings.ToLower(filepath.Ext(fileName))

	// Get MIME type
	mimeType := fs.getMimeType(filePath)

	// Check each sort rule
	for _, rule := range fs.Config.SortRules {
		// Check extension
		for _, ext := range rule.Extensions {
			if fileExt == strings.ToLower(ext) {
				return rule.Folder
			}
		}

		// Check MIME type
		for _, mt := range rule.MimeTypes {
			if strings.HasPrefix(mimeType, mt) {
				return rule.Folder
			}
		}

		// Check keywords in filename
		for _, keyword := range rule.Keywords {
			if strings.Contains(strings.ToLower(fileName), strings.ToLower(keyword)) {
				return rule.Folder
			}
		}
	}

	return ""
}

// getMimeType gets the MIME type of a file
func (fs *FileSorter) getMimeType(filePath string) string {
	// Check custom MIME map first
	ext := strings.ToLower(filepath.Ext(filePath))
	if customMime, exists := fs.Config.CustomMimeMap[ext]; exists {
		return customMime
	}

	// Use standard MIME detection
	return mime.TypeByExtension(ext)
}

// moveFile moves a file to the target folder
func (fs *FileSorter) moveFile(filePath, targetFolder string) error {
	fileName := filepath.Base(filePath)
	targetPath := filepath.Join(targetFolder, fileName)

	// Handle duplicate filenames
	targetPath = fs.handleDuplicateFilename(targetPath)

	if fs.Config.DryRun {
		fs.Logger.Printf("DRY RUN: Would move %s -> %s", filePath, targetPath)
		return nil
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Move the file
	if err := os.Rename(filePath, targetPath); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	fs.Logger.Printf("Moved: %s -> %s", filePath, targetPath)
	return nil
}

// handleDuplicateFilename handles duplicate filenames by appending a number
func (fs *FileSorter) handleDuplicateFilename(targetPath string) string {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return targetPath
	}

	dir := filepath.Dir(targetPath)
	name := filepath.Base(targetPath)
	ext := filepath.Ext(name)
	nameWithoutExt := strings.TrimSuffix(name, ext)

	counter := 1
	for {
		newName := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
		newPath := filepath.Join(dir, newName)

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
		counter++
	}
}

// createSortFolders creates the sort folders if they don't exist
func (fs *FileSorter) createSortFolders() error {
	for _, rule := range fs.Config.SortRules {
		if err := os.MkdirAll(rule.Folder, 0755); err != nil {
			return fmt.Errorf("failed to create folder %s: %w", rule.Folder, err)
		}
	}
	return nil
}
