package cmd

import (
	"fmt"

	"github.com/RichestHumanAlive/ppm_cli/pkg/manager"
	"github.com/RichestHumanAlive/ppm_cli/pkg/manager/npm"
	"github.com/RichestHumanAlive/ppm_cli/pkg/manager/pip"
	"github.com/RichestHumanAlive/ppm_cli/pkg/manager/scoop"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type model struct {
	spinner  spinner.Model
	quitting bool
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n %s Installing package...\n", m.spinner.View())
}

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [package]",
		Short: "Install a package using the appropriate package manager",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg := args[0]

			// Initialize manager
			mgr := manager.New()
			mgr.RegisterManager(npm.New())
			mgr.RegisterManager(pip.New())
			mgr.RegisterManager(scoop.New())

			// Create a channel to receive installation result
			errChan := make(chan error)

			// Create and start spinner
			s := spinner.New()
			s.Spinner = spinner.Dot
			s.Style = s.Style.Foreground(s.Style.GetForeground())

			m := model{spinner: s}
			p := tea.NewProgram(m)

			// Start installation in goroutine
			go func() {
				// Try each package manager
				var installErr error
				for _, pm := range mgr.GetManagers() {
					if pm.IsAvailable() {
						if err := pm.Install(pkg); err == nil {
							errChan <- nil
							return
						} else {
							installErr = err
						}
					}
				}
				if installErr != nil {
					errChan <- fmt.Errorf("no package manager could install %s: %v", pkg, installErr)
				} else {
					errChan <- fmt.Errorf("no available package manager found")
				}
			}()

			// Start spinner
			go func() {
				if err := p.Start(); err != nil {
					fmt.Printf("Error starting spinner: %v\n", err)
				}
			}()

			// Wait for installation result
			err := <-errChan
			p.Quit()

			if err != nil {
				return fmt.Errorf("installation failed: %v", err)
			}

			fmt.Printf("\n✓ Successfully installed %s\n", pkg)
			return nil
		},
	}

	return cmd
}
