package pip

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/RichestHumanAlive/ppm_cli/pkg/manager"
)

type PIPManager struct{}

func New() *PIPManager {
	return &PIPManager{}
}

func (p *PIPManager) GetName() string {
	return "pip"
}

func (p *PIPManager) Install(pkg string) error {
	cmd := exec.Command("pip", "install", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pip install failed: %v\n%s", err, string(output))
	}
	return nil
}

type PyPIResponse struct {
	Info struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Summary     string `json:"summary"`
		Author      string `json:"author"`
		HomePage    string `json:"home_page"`
		ProjectURL  string `json:"project_url"`
		DownloadURL string `json:"download_url"`
	} `json:"info"`
}

func (p *PIPManager) Search(query string) ([]manager.Package, error) {
	// Use pip search command (Note: pip search is deprecated, using pip index instead)
	cmd := exec.Command("pip", "index", "versions", query)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback to simple package info
		cmd = exec.Command("pip", "show", query)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("pip search failed: %v", err)
		}
	}

	// Parse the output into a Package struct
	pkg := manager.Package{
		Name:     query,
		Provider: "pip",
		Score:    0.8, // Default score for exact matches
	}

	// Parse the output line by line
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch strings.ToLower(key) {
		case "version":
			pkg.Version = value
		case "summary":
			pkg.Description = value
		case "author":
			pkg.Author = value
		case "home-page":
			pkg.Homepage = value
		}
	}

	return []manager.Package{pkg}, nil
}

func (p *PIPManager) Update(pkg string) error {
	args := []string{"install", "--upgrade"}
	if pkg != "" {
		args = append(args, pkg)
	} else {
		args = append(args, "-r", "requirements.txt")
	}

	cmd := exec.Command("pip", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pip update failed: %v\n%s", err, string(output))
	}
	return nil
}

func (p *PIPManager) Remove(pkg string) error {
	cmd := exec.Command("pip", "uninstall", "-y", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pip uninstall failed: %v\n%s", err, string(output))
	}
	return nil
}

func (p *PIPManager) IsAvailable() bool {
	cmd := exec.Command("pip", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
