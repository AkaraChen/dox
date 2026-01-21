package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildCommand_Status(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildStatus([]string{})
	require.NoError(t, err)

	// Status command is essentially ps with enhanced output
	assert.True(t, sliceContains(cmd, "ps"))
}

func TestBuildCommand_StatusWithWatch(t *testing.T) {
	fixtureDir := setupFixture(t, "simple")

	b := NewBuilder(fixtureDir, nil, "")
	cmd, err := b.BuildStatus([]string{"--watch"})
	require.NoError(t, err)

	assert.True(t, sliceContains(cmd, "ps"))
	// Watch flag might be handled differently, but command should include ps
}

func TestFormatStatusOutput(t *testing.T) {
	// Test status output formatting
	output := FormatStatusOutput([]ServiceStatus{
	 {Name: "web", State: "running", Ports: "0.0.0.0:8080->80/tcp"},
	 {Name: "db", State: "running", Ports: "0.0.0.0:5432->5432/tcp"},
	})

	assert.Contains(t, output, "web")
	assert.Contains(t, output, "running")
	assert.Contains(t, output, "0.0.0.0:8080->80/tcp")
	assert.Contains(t, output, "db")
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name  string
	State string
	Ports string
}

// FormatStatusOutput formats status output for display
func FormatStatusOutput(services []ServiceStatus) string {
	var output string
	for _, s := range services {
	 output += s.Name + "\t" + s.State + "\t" + s.Ports + "\n"
	}
	return output
}
