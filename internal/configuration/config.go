package configuration

import (
	"encoding/json"
	"os"
)

func LoadConfig(configPath string) (*Config, error) {
	defaultConfig := &Config{
		WatchFolder:   "./Downloads",
		EnableLogging: true,
		LogFile:       "./auto-sort.log",
		DryRun:        false,
		IgnorePatterns: []string{
			"*.tmp",
			"*.part",
			".*",
		},
		CustomMimeMap: map[string]string{
			".dmg":  "application/x-apple-diskimage",
			".deb":  "application/x-debian-package",
			".rpm":  "application/x-rpm",
			".appx": "application/appx",
		},
		SortRules: map[string]SortRule{
			"images": {
				Folder:     "./Downloads/Images",
				Extensions: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico"},
				MimeTypes:  []string{"image/"},
			},
			"documents": {
				Folder:     "./Downloads/Documents",
				Extensions: []string{".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt", ".xls", ".xlsx", ".ppt", ".pptx"},
				MimeTypes:  []string{"application/pdf", "application/msword", "text/"},
			},
			"videos": {
				Folder:     "./Downloads/Videos",
				Extensions: []string{".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm", ".m4v"},
				MimeTypes:  []string{"video/"},
			},
			"audio": {
				Folder:     "./Downloads/Audio",
				Extensions: []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a"},
				MimeTypes:  []string{"audio/"},
			},
			"archives": {
				Folder:     "./Downloads/Archives",
				Extensions: []string{".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz"},
				MimeTypes:  []string{"application/zip", "application/x-rar"},
			},
			"executables": {
				Folder:     "./Downloads/Programs",
				Extensions: []string{".exe", ".msi", ".deb", ".rpm", ".dmg", ".pkg", ".appx"},
				MimeTypes:  []string{"application/x-executable", "application/x-msdos-program"},
			},
			"code": {
				Folder:     "./Downloads/Code",
				Extensions: []string{".go", ".py", ".js", ".html", ".css", ".java", ".cpp", ".c", ".h", ".php", ".rb", ".rs"},
				MimeTypes:  []string{"text/x-go", "text/x-python", "text/javascript"},
			},
		},
	}

	// If config file doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return defaultConfig, saveConfig(configPath, defaultConfig)
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(configPath string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
