package configuration

type Config struct {
	WatchFolder    string              `json:"watch_folder"`
	SortRules      map[string]SortRule `json:"sort_rules"`
	EnableLogging  bool                `json:"enable_logging"`
	LogFile        string              `json:"log_file"`
	DryRun         bool                `json:"dry_run"`
	IgnorePatterns []string            `json:"ignore_patterns"`
	CustomMimeMap  map[string]string   `json:"custom_mime_map"`
}

type SortRule struct {
	Folder     string   `json:"folder"`
	Extensions []string `json:"extension"`
	MimeTypes  []string `json:"mime_types"`
	Keywords   []string `json:"keywords"`
}
