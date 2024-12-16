package npm

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/RichestHumanAlive/ppm_cli/pkg/manager"
)

type NPMManager struct{}

func New() *NPMManager {
	return &NPMManager{}
}

func (n *NPMManager) GetName() string {
	return "npm"
}

func (n *NPMManager) Install(pkg string) error {
	cmd := exec.Command("npm", "install", "-g", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed: %v\n%s", err, string(output))
	}
	return nil
}

type NPMSearchResult struct {
	Objects []struct {
		Package struct {
			Name        string  `json:"name"`
			Version     string  `json:"version"`
			Description string  `json:"description"`
			Author      struct {
				Name string `json:"name"`
			} `json:"author"`
			Links struct {
				Homepage   string `json:"homepage"`
				Repository string `json:"repository"`
			} `json:"links"`
			Score struct {
				Final float64 `json:"final"`
			} `json:"score"`
		} `json:"package"`
	} `json:"objects"`
}

func (n *NPMManager) Search(query string) ([]manager.Package, error) {
	cmd := exec.Command("npm", "search", "--json", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("npm search failed: %v", err)
	}

	var result NPMSearchResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse npm search results: %v", err)
	}

	packages := make([]manager.Package, 0, len(result.Objects))
	for _, obj := range result.Objects {
		pkg := manager.Package{
			Name:        obj.Package.Name,
			Version:     obj.Package.Version,
			Description: obj.Package.Description,
			Author:      obj.Package.Author.Name,
			Provider:    "npm",
			Score:      obj.Package.Score.Final,
			Homepage:   obj.Package.Links.Homepage,
			Repository: obj.Package.Links.Repository,
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

func (n *NPMManager) Update(pkg string) error {
	args := []string{"update", "-g"}
	if pkg != "" {
		args = append(args, pkg)
	}
	
	cmd := exec.Command("npm", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm update failed: %v\n%s", err, string(output))
	}
	return nil
}

func (n *NPMManager) Remove(pkg string) error {
	cmd := exec.Command("npm", "uninstall", "-g", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm uninstall failed: %v\n%s", err, string(output))
	}
	return nil
}

func (n *NPMManager) IsAvailable() bool {
	cmd := exec.Command("npm", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
