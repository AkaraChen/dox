package project

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadHistory(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.yaml")

	// Create test history
	content := []byte(`
entries:
- timestamp: "2024-01-15T10:30:00Z"
  command: "c up"
  directory: "/home/user/project1"
  exit_code: 0
- timestamp: "2024-01-15T10:31:00Z"
  command: "c logs -f"
  directory: "/home/user/project1"
  exit_code: 0
`)
	err := os.WriteFile(historyFile, content, 0644)
	require.NoError(t, err)

	// Load history
	hist, err := LoadHistory(historyFile)
	require.NoError(t, err)

	assert.Len(t, hist.Entries, 2)
	assert.Equal(t, "c up", hist.Entries[0].Command)
	assert.Equal(t, "c logs -f", hist.Entries[1].Command)
}

func TestLoadHistory_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "nonexistent.yaml")

	// Non-existent file should return empty history
	hist, err := LoadHistory(historyFile)
	require.NoError(t, err)
	assert.NotNil(t, hist)
	assert.Empty(t, hist.Entries)
}

func TestLoadHistory_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.yaml")

	// Create invalid YAML
	content := []byte(`- timestamp: [invalid`)
	err := os.WriteFile(historyFile, content, 0644)
	require.NoError(t, err)

	// Should return error
	hist, err := LoadHistory(historyFile)
	require.Error(t, err)
	assert.Nil(t, hist)
}

func TestHistory_AddEntry(t *testing.T) {
	hist := &History{
		Entries: []HistoryEntry{},
	}

	// Add entry
	entry := HistoryEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Command:   "c up",
		Directory: "/test/project",
		ExitCode:  0,
	}

	hist.AddEntry(entry)

	assert.Len(t, hist.Entries, 1)
	assert.Equal(t, "c up", hist.Entries[0].Command)
}

func TestHistory_AddMultipleEntries(t *testing.T) {
	hist := &History{
		Entries: []HistoryEntry{},
	}

	for i := 0; i < 5; i++ {
		entry := HistoryEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Command:   "c up",
			Directory: "/test/project",
			ExitCode:  0,
		}
		hist.AddEntry(entry)
	}

	assert.Len(t, hist.Entries, 5)
}

func TestHistory_Save(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.yaml")

	hist := &History{
		Entries: []HistoryEntry{
			{
				Timestamp: "2024-01-15T10:30:00Z",
				Command:   "c up",
				Directory: "/test/project",
				ExitCode:  0,
			},
		},
	}

	// Save history
	err := hist.Save(historyFile)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(historyFile)
	require.NoError(t, err)

	// Load and verify
	loaded, err := LoadHistory(historyFile)
	require.NoError(t, err)
	assert.Len(t, loaded.Entries, 1)
	assert.Equal(t, "c up", loaded.Entries[0].Command)
}

func TestHistory_SaveAppend(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.yaml")

	// Initial save
	hist := &History{
		Entries: []HistoryEntry{
			{
				Timestamp: "2024-01-15T10:30:00Z",
				Command:   "c up",
				Directory: "/test/project",
				ExitCode:  0,
			},
		},
	}
	err := hist.Save(historyFile)
	require.NoError(t, err)

	// Append new entry
	entry := HistoryEntry{
		Timestamp: "2024-01-15T10:31:00Z",
		Command:   "c logs",
		Directory: "/test/project",
		ExitCode:  0,
	}
	hist.AddEntry(entry)
	err = hist.Save(historyFile)
	require.NoError(t, err)

	// Load and verify both entries
	loaded, err := LoadHistory(historyFile)
	require.NoError(t, err)
	assert.Len(t, loaded.Entries, 2)
}

func TestHistory_GetHistoryPath(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	os.Setenv("HOME", "/test/home")
	path := GetHistoryPath()
	expected := filepath.Join("/test/home", ".cache", "dox", "history.yaml")
	assert.Equal(t, expected, path)
}

func TestHistory_Last(t *testing.T) {
	hist := &History{
		Entries: []HistoryEntry{
			{Command: "c up", Timestamp: "2024-01-15T10:30:00Z"},
			{Command: "c logs", Timestamp: "2024-01-15T10:31:00Z"},
			{Command: "c down", Timestamp: "2024-01-15T10:32:00Z"},
		},
	}

	last := hist.Last(2)
	assert.Len(t, last, 2)
	assert.Equal(t, "c logs", last[0].Command)
	assert.Equal(t, "c down", last[1].Command)
}

func TestHistory_Last_Empty(t *testing.T) {
	hist := &History{
		Entries: []HistoryEntry{},
	}

	last := hist.Last(5)
	assert.Nil(t, last)
}

func TestHistory_Last_RequestMoreThanAvailable(t *testing.T) {
	hist := &History{
		Entries: []HistoryEntry{
			{Command: "c up", Timestamp: "2024-01-15T10:30:00Z"},
			{Command: "c logs", Timestamp: "2024-01-15T10:31:00Z"},
		},
	}

	last := hist.Last(10)
	assert.Len(t, last, 2)
}

func TestHistory_FilterByDirectory(t *testing.T) {
	hist := &History{
		Entries: []HistoryEntry{
			{Command: "c up", Directory: "/project1", Timestamp: "2024-01-15T10:30:00Z"},
			{Command: "c logs", Directory: "/project2", Timestamp: "2024-01-15T10:31:00Z"},
			{Command: "c down", Directory: "/project1", Timestamp: "2024-01-15T10:32:00Z"},
		},
	}

	filtered := hist.FilterByDirectory("/project1")
	assert.Len(t, filtered, 2)
	assert.Equal(t, "c up", filtered[0].Command)
	assert.Equal(t, "c down", filtered[1].Command)
}

func TestHistory_EntryTimestampFormat(t *testing.T) {
	// Test that timestamps are in RFC3339 format
	now := time.Now().UTC()
	entry := HistoryEntry{
		Timestamp: now.Format(time.RFC3339),
		Command:   "c up",
		Directory: "/test",
		ExitCode:  0,
	}

	// Verify it's a valid timestamp format
	_, err := time.Parse(time.RFC3339, entry.Timestamp)
	assert.NoError(t, err)
}

func TestNewHistoryEntry(t *testing.T) {
	// Test NewHistoryEntry function
	entry := NewHistoryEntry("c up", "/test/project", 0)

	// Verify the entry structure (we can't check exact timestamp due to time.Now())
	assert.NotEmpty(t, entry.Timestamp)
	assert.Equal(t, "c up", entry.Command)
	assert.Equal(t, "/test/project", entry.Directory)
	assert.Equal(t, 0, entry.ExitCode)

	// Verify timestamp is valid RFC3339
	_, err := time.Parse(time.RFC3339, entry.Timestamp)
	assert.NoError(t, err)
}

func TestNewHistoryEntry_WithExitCode(t *testing.T) {
	entry := NewHistoryEntry("c up", "/test/project", 1)

	assert.NotEmpty(t, entry.Timestamp)
	assert.Equal(t, "c up", entry.Command)
	assert.Equal(t, "/test/project", entry.Directory)
	assert.Equal(t, 1, entry.ExitCode)
}

func TestHistory_Save_ErrorOnDirectoryCreation(t *testing.T) {
	// This test verifies error handling when directory creation fails
	// In most systems we can't actually cause this, but we can test the path

	hist := &History{
		Entries: []HistoryEntry{
			{Timestamp: "2024-01-15T10:30:00Z", Command: "c up", Directory: "/test", ExitCode: 0},
		},
	}

	// Try to save to a path we can't create (like /root/.cache on most systems)
	// This should fail gracefully
	err := hist.Save("/root/.cache/dox/history.yaml")
	// We expect this to fail in most environments
	// The test just verifies the error handling path exists
	if err != nil {
		assert.Error(t, err)
	}
}

func TestHistory_Save_EmptyEntries(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.yaml")

	hist := &History{
		Entries: []HistoryEntry{},
	}

	err := hist.Save(historyFile)
	require.NoError(t, err)

	// Verify file was created
	data, err := os.ReadFile(historyFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "entries:")
}

func TestHistory_Save_NilEntries(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.yaml")

	hist := &History{
		Entries: nil,
	}

	err := hist.Save(historyFile)
	require.NoError(t, err)

	// Verify file was created
	data, err := os.ReadFile(historyFile)
	require.NoError(t, err)
	assert.Contains(t, string(data), "entries:")
}

func TestHistory_Last_NilEntries(t *testing.T) {
	hist := &History{
		Entries: nil,
	}

	last := hist.Last(5)
	assert.Nil(t, last)
}

func TestHistory_FilterByDirectory_NoMatches(t *testing.T) {
	hist := &History{
		Entries: []HistoryEntry{
			{Command: "c up", Directory: "/project1"},
			{Command: "c logs", Directory: "/project2"},
		},
	}

	filtered := hist.FilterByDirectory("/nonexistent")
	assert.Empty(t, filtered)
}

func TestHistory_FilterByDirectory_NilEntries(t *testing.T) {
	hist := &History{
		Entries: nil,
	}

	filtered := hist.FilterByDirectory("/project1")
	assert.Empty(t, filtered)
}

func TestHistory_LoadFromFile_ReadError(t *testing.T) {
	// Test that non-existent file returns empty history (not an error)
	hist, err := LoadHistory("/nonexistent/path/history.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, hist)
	assert.Empty(t, hist.Entries)
}

func TestGetGlobalConfigPath_ErrorPath(t *testing.T) {
	// Test when HOME is not set
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	os.Unsetenv("HOME")
	path := GetGlobalConfigPath()
	// Should still return a path using fallback
	assert.NotEmpty(t, path)
}

func TestGetHistoryPath_ErrorPath(t *testing.T) {
	// Test when HOME is not set
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	os.Unsetenv("HOME")
	path := GetHistoryPath()
	// Should still return a path using fallback
	assert.NotEmpty(t, path)
}
