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
}

// Package represents a package in any package manager
type Package struct {
	Name        string
	Version     string
	Description string
	Source      string // Which package manager this package is from
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
	if pm.IsAvailable() {
		m.managers = append(m.managers, pm)
	}
}
