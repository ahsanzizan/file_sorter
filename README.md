# Auto-Sort Downloads

A hopefully smart Go application that automatically organizes your Downloads folder by monitoring file changes and sorting files into appropriate folders based on their type, extension, or filename patterns.

## Features

- **Real-time Monitoring**: Uses `fsnotify` to watch your Downloads folder for new files
- **Intelligent Sorting**: Categorizes files by extension, MIME type, or filename keywords
- **Configurable Rules**: Fully customizable sorting rules via JSON configuration
- **Duplicate Handling**: Automatically handles duplicate filenames by appending numbers
- **Dry Run Mode**: Test your configuration without actually moving files
- **Comprehensive Logging**: Optional logging to file or console
- **Cross-platform**: Works on Windows, macOS, and Linux
- **Startup Processing**: Processes existing files when the program starts
- **File Completion Detection**: Waits for files to finish downloading before moving them

## Quick Start

1. **Install dependencies:**

   ```bash
   go mod tidy
   ```

2. **Generate default configuration:**

   ```bash
   make config
   # or
   go run main.go -generate-config
   ```

3. **Run in dry-run mode first (recommended):**

   ```bash
   make dry-run
   # or
   go run main.go -dry-run
   ```

4. **Run the application:**
   ```bash
   make run
   # or
   go run main.go
   ```

## Configuration

The application uses a `config.json` file to define sorting rules. Here's what each section does:

### Basic Settings

```json
{
  "watch_folder": "./Downloads",
  "enable_logging": true,
  "log_file": "./auto-sort.log",
  "dry_run": false,
  "ignore_patterns": ["*.tmp", "*.part", ".*"]
}
```

### Sort Rules

Each rule can match files by:

- **Extensions**: File extensions (e.g., `.jpg`, `.pdf`)
- **MIME Types**: MIME type prefixes (e.g., `image/`, `video/`)
- **Keywords**: Words in the filename (e.g., `screenshot`, `invoice`)

```json
{
  "sort_rules": {
    "images": {
      "folder": "./Downloads/Images",
      "extensions": [".jpg", ".jpeg", ".png", ".gif"],
      "mime_types": ["image/"],
      "keywords": ["screenshot", "photo"]
    }
  }
}
```

### Custom MIME Types

Add custom MIME type mappings for files that aren't recognized by default:

```json
{
  "custom_mime_map": {
    ".dmg": "application/x-apple-diskimage",
    ".deb": "application/x-debian-package"
  }
}
```

## Command Line Options

```bash
# Generate default configuration
go run main.go -generate-config

# Run in dry-run mode (shows what would be moved without actually moving)
go run main.go -dry-run

# Specify custom config file location
go run main.go -config /path/to/config.json
```

## Building

### Single Platform

```bash
make build
```

### Multiple Platforms

```bash
make build-all
```

This creates binaries for:

- Windows (amd64)
- macOS (amd64, arm64)
- Linux (amd64, arm64)

## Installation

### Install to System PATH

```bash
make install
```

### Uninstall

```bash
make uninstall
```

## Default Sort Categories

The application comes with predefined categories:

- **Images**: JPG, PNG, GIF, SVG, WebP, etc.
- **Documents**: PDF, Word, Excel, PowerPoint, Text files
- **Videos**: MP4, AVI, MKV, MOV, WebM, etc.
- **Audio**: MP3, WAV, FLAC, AAC, OGG, etc.
- **Archives**: ZIP, RAR, 7Z, TAR, GZ, etc.
- **Programs**: EXE, MSI, DEB, RPM, DMG, etc.
- **Code**: Go, Python, JavaScript, HTML, CSS, etc.

## Advanced Features

### Ignore Patterns

Specify glob patterns for files to ignore:

```json
{
  "ignore_patterns": ["*.tmp", "*.part", ".*", "desktop.ini", "Thumbs.db"]
}
```

### Duplicate File Handling

When a file with the same name exists in the destination folder, the application automatically appends a number:

- `document.pdf` → `document_1.pdf`
- `image.jpg` → `image_1.jpg`

### File Completion Detection

The application waits for files to finish downloading by checking file size stability before moving them.

## Logging

Enable logging to track all file movements:

```json
{
  "enable_logging": true,
  "log_file": "./auto-sort.log"
}
```

Log entries include:

- File movement operations
- Errors and warnings
- Startup and shutdown events
- Dry-run simulation results

## Usage Examples

### Basic Usage

```bash
# Start monitoring with default settings
./auto-sort

# Test configuration first
./auto-sort -dry-run
```

### Custom Configuration

```bash
# Use custom config file
./auto-sort -config /home/user/my-config.json

# Generate config in specific location
./auto-sort -generate-config -config /home/user/custom-config.json
```

### Running as a Service

#### Linux (systemd)

Create `/etc/systemd/system/auto-sort.service`:

```ini
[Unit]
Description=Auto-Sort Downloads
After=network.target

[Service]
Type=simple
User=yourusername
ExecStart=/usr/local/bin/auto-sort -config /home/yourusername/.config/auto-sort.json
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable auto-sort.service
sudo systemctl start auto-sort.service
```

#### macOS (launchd)

Create `~/Library/LaunchAgents/com.auto-sort.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.auto-sort</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/auto-sort</string>
        <string>-config</string>
        <string>/Users/yourusername/.config/auto-sort.json</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

Load the service:

```bash
launchctl load ~/Library/LaunchAgents/com.auto-sort.plist
```

## Troubleshooting

### Common Issues

1. **Permission Denied**: Ensure the application has read/write access to the Downloads folder and destination folders.

2. **Files Not Moving**: Check if files match the ignore patterns or if they're still being written.

3. **Configuration Not Found**: Use `-generate-config` to create a default configuration file.

4. **Watcher Errors**: Some filesystems (like network drives) may not support file watching. Check the logs for specific error messages.

### Debug Mode

Run with logging enabled and check the log file for detailed information about what the application is doing.

## Contributing

Feel free to submit issues and enhancement requests!

## License

This project is open source and available under the MIT License.
