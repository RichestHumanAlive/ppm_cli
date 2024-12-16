package manager

// PackageManager defines the interface that all package managers must implement
type PackageManager interface {
	// Install installs a package
	Install(pkg string) error
	
	// Search searches for a package
	Search(query string) ([]Package, error)
	
	// Update updates a package or all packages if pkg is empty
	Update(pkg string) error
	
	// Remove removes a package
	Remove(pkg string) error
	
	// IsAvailable checks if this package manager is available on the system
	IsAvailable() bool
	
	// GetName returns the name of the package manager
	GetName() string
}

// Package represents a package in any package manager
type Package struct {
	Name        string  // Package name
	Version     string  // Latest version
	Description string  // Package description
	Author      string  // Package author/maintainer
	Provider    string  // Package manager (npm, pip, scoop)
	Score      float64 // Relevance score (0-1)
	Downloads  int64   // Number of downloads (if available)
	Homepage   string  // Package homepage URL
	Repository string  // Source code repository URL
}

// Manager handles operations across multiple package managers
type Manager struct {
	managers []PackageManager
}

// New creates a new package manager handler
func New() *Manager {
	return &Manager{
		managers: make([]PackageManager, 0),
	}
}

// RegisterManager adds a package manager to the handler
func (m *Manager) RegisterManager(pm PackageManager) {
	m.managers = append(m.managers, pm)
}

// GetManagers returns the list of registered package managers
func (m *Manager) GetManagers() []PackageManager {
	return m.managers
}

// SearchAcrossAll searches for packages across all available package managers
func (m *Manager) SearchAcrossAll(query string) ([]Package, error) {
	results := make([]Package, 0)
	
	for _, pm := range m.managers {
		if !pm.IsAvailable() {
			continue
		}
		
		pkgs, err := pm.Search(query)
		if err != nil {
			// Log error but continue with other package managers
			continue
		}
		results = append(results, pkgs...)
	}
	
	return results, nil
}
