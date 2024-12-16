package scoop

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/RichestHumanAlive/ppm_cli/pkg/manager"
)

type ScoopManager struct{}

func New() *ScoopManager {
	return &ScoopManager{}
}

func (s *ScoopManager) GetName() string {
	return "scoop"
}

func (s *ScoopManager) Install(pkg string) error {
	cmd := exec.Command("scoop", "install", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scoop install failed: %v\n%s", err, string(output))
	}
	return nil
}

type ScoopApp struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	License     string `json:"license"`
}

func (s *ScoopManager) Search(query string) ([]manager.Package, error) {
	// First, update the scoop database
	updateCmd := exec.Command("scoop", "update")
	updateCmd.Run() // Ignore errors, just try to update

	// Search for the package
	cmd := exec.Command("scoop", "search", query)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("scoop search failed: %v", err)
	}

	// Parse the output
	results := make([]manager.Package, 0)
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		// Scoop output format is typically: 'name (bucket): description'
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		
		namePart := strings.TrimSpace(parts[0])
		description := strings.TrimSpace(parts[1])
		
		// Extract name and bucket
		name := namePart
		if strings.Contains(namePart, " (") {
			nameParts := strings.Split(strings.Trim(namePart, ")"), " (")
			name = nameParts[0]
		}

		// Get more details about the package
		infoCmd := exec.Command("scoop", "info", name)
		infoOutput, err := infoCmd.CombinedOutput()
		if err != nil {
			continue
		}

		pkg := manager.Package{
			Name:        name,
			Description: description,
			Provider:    "scoop",
			Score:      0.7, // Default score
		}

		// Parse the info output
		infoLines := strings.Split(string(infoOutput), "\n")
		for _, infoLine := range infoLines {
			if strings.Contains(infoLine, "Version:") {
				pkg.Version = strings.TrimSpace(strings.TrimPrefix(infoLine, "Version:"))
			} else if strings.Contains(infoLine, "Website:") {
				pkg.Homepage = strings.TrimSpace(strings.TrimPrefix(infoLine, "Website:"))
			}
		}

		results = append(results, pkg)
	}

	return results, nil
}

func (s *ScoopManager) Update(pkg string) error {
	args := []string{"update"}
	if pkg != "" {
		args = append(args, pkg)
	}

	cmd := exec.Command("scoop", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scoop update failed: %v\n%s", err, string(output))
	}
	return nil
}

func (s *ScoopManager) Remove(pkg string) error {
	cmd := exec.Command("scoop", "uninstall", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scoop uninstall failed: %v\n%s", err, string(output))
	}
	return nil
}

func (s *ScoopManager) IsAvailable() bool {
	cmd := exec.Command("scoop", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
