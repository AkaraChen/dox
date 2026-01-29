package project

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// History represents command execution history
type History struct {
	Entries []HistoryEntry `yaml:"entries"`
}

// HistoryEntry represents a single command execution record
type HistoryEntry struct {
	Timestamp string `yaml:"timestamp"`
	Command   string `yaml:"command"`
	Directory string `yaml:"directory"`
	ExitCode  int    `yaml:"exit_code"`
}

// GetHistoryPath returns the default path for the history file
func GetHistoryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}
	return filepath.Join(homeDir, ".cache", "dox", "history.yaml")
}

// LoadHistory loads the command history from the specified path.
// Returns an empty history if the file doesn't exist (not an error).
func LoadHistory(path string) (*History, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{Entries: []HistoryEntry{}}, nil
		}
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var hist History
	if err := yaml.Unmarshal(data, &hist); err != nil {
		return nil, fmt.Errorf("failed to parse history file: %w", err)
	}

	// Initialize entries if nil
	if hist.Entries == nil {
		hist.Entries = []HistoryEntry{}
	}

	return &hist, nil
}

// AddEntry adds a new entry to the history
func (h *History) AddEntry(entry HistoryEntry) {
	h.Entries = append(h.Entries, entry)
}

// Save saves the history to the specified path
func (h *History) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	data, err := yaml.Marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

// Last returns the last n entries from history.
// Returns nil if history is empty.
func (h *History) Last(n int) []HistoryEntry {
	if len(h.Entries) == 0 {
		return nil
	}

	start := len(h.Entries) - n
	if start < 0 {
		start = 0
	}

	result := make([]HistoryEntry, 0, len(h.Entries)-start)
	result = append(result, h.Entries[start:]...)
	return result
}

// FilterByDirectory returns entries that match the specified directory
func (h *History) FilterByDirectory(dir string) []HistoryEntry {
	filtered := make([]HistoryEntry, 0)
	for _, entry := range h.Entries {
		if entry.Directory == dir {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// NewHistoryEntry creates a new history entry with the current timestamp
func NewHistoryEntry(command, directory string, exitCode int) HistoryEntry {
	return HistoryEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Command:   command,
		Directory: directory,
		ExitCode:  exitCode,
	}
}
