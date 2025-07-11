package filesorter

import (
	"file-sorter/internal/configuration"
	"log"

	"github.com/fsnotify/fsnotify"
)

type FileSorter struct {
	Config  *configuration.Config
	Watcher *fsnotify.Watcher
	Logger  *log.Logger
}
