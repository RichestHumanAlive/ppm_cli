package cmd

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/RichestHumanAlive/ppm_cli/pkg/manager"
	"github.com/RichestHumanAlive/ppm_cli/pkg/manager/npm"
	"github.com/RichestHumanAlive/ppm_cli/pkg/manager/pip"
	"github.com/RichestHumanAlive/ppm_cli/pkg/manager/scoop"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF69B4")).
		Bold(true)

	providerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Italic(true)

	versionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#87CEEB"))

	descStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true)

	cellStyle = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1)
)

type searchModel struct {
	spinner  spinner.Model
	quitting bool
	query    string
}

func (m searchModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m searchModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n %s Searching for '%s' across package managers...\n", m.spinner.View(), m.query)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func renderTable(results []manager.Package) string {
	if len(results) == 0 {
		return "No packages found"
	}

	// Define columns
	headers := []string{"ID", "Name", "Version", "Provider", "Description"}
	widths := make([]int, len(headers))

	// Calculate maximum widths for each column
	for i, header := range headers {
		widths[i] = utf8.RuneCountInString(header)
	}

	for _, pkg := range results {
		name := pkg.Name
		if len(name) > widths[1] {
			widths[1] = len(name)
		}
		if len(pkg.Version) > widths[2] {
			widths[2] = len(pkg.Version)
		}
		if len(pkg.Provider) > widths[3] {
			widths[3] = len(pkg.Provider)
		}
		desc := pkg.Description
		if len(desc) > 50 { // Truncate long descriptions
			desc = desc[:47] + "..."
		}
		if len(desc) > widths[4] {
			widths[4] = len(desc)
		}
	}

	// Build the table
	var sb strings.Builder

	// Write headers
	sb.WriteString("┌")
	for i, w := range widths {
		sb.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			sb.WriteString("┬")
		}
	}
	sb.WriteString("┐\n")

	// Write header row
	sb.WriteString("│")
	for i, header := range headers {
		cell := headerStyle.Render(fmt.Sprintf("%-*s", widths[i], header))
		sb.WriteString(cellStyle.Render(cell))
		sb.WriteString("│")
	}
	sb.WriteString("\n")

	// Write separator
	sb.WriteString("├")
	for i, w := range widths {
		sb.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			sb.WriteString("┼")
		}
	}
	sb.WriteString("┤\n")

	// Write data rows
	for i, pkg := range results {
		desc := pkg.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}

		sb.WriteString("│")
		cells := []string{
			fmt.Sprintf("%d", i+1),
			titleStyle.Render(pkg.Name),
			versionStyle.Render(pkg.Version),
			providerStyle.Render(pkg.Provider),
			descStyle.Render(desc),
		}

		for j, cell := range cells {
			paddedCell := fmt.Sprintf("%-*s", widths[j], lipgloss.NewStyle().Render(cell))
			sb.WriteString(cellStyle.Render(paddedCell))
			sb.WriteString("│")
		}
		sb.WriteString("\n")
	}

	// Write bottom border
	sb.WriteString("└")
	for i, w := range widths {
		sb.WriteString(strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			sb.WriteString("┴")
		}
	}
	sb.WriteString("┘\n")

	return sb.String()
}

func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for packages across all package managers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			// Initialize spinner
			s := spinner.New()
			s.Spinner = spinner.Dot
			s.Style = s.Style.Foreground(s.Style.GetForeground())

			// Initialize manager
			mgr := manager.New()
			mgr.RegisterManager(npm.New())
			mgr.RegisterManager(pip.New())
			mgr.RegisterManager(scoop.New())

			// Create channels for each package manager's results
			npmChan := make(chan []manager.Package)
			pipChan := make(chan []manager.Package)
			scoopChan := make(chan []manager.Package)

			// Create and start spinner
			m := searchModel{spinner: s, query: query}
			p := tea.NewProgram(m)

			// Start searches in parallel
			go func() {
				if npm := npm.New(); npm.IsAvailable() {
					results, err := npm.Search(query)
					if err != nil {
						npmChan <- nil
					} else {
						npmChan <- results
					}
				} else {
					npmChan <- nil
				}
			}()

			go func() {
				if pip := pip.New(); pip.IsAvailable() {
					results, err := pip.Search(query)
					if err != nil {
						pipChan <- nil
					} else {
						pipChan <- results
					}
				} else {
					pipChan <- nil
				}
			}()

			go func() {
				if scoop := scoop.New(); scoop.IsAvailable() {
					results, err := scoop.Search(query)
					if err != nil {
						scoopChan <- nil
					} else {
						scoopChan <- results
					}
				} else {
					scoopChan <- nil
				}
			}()

			// Start spinner
			go func() {
				if err := p.Start(); err != nil {
					fmt.Printf("Error starting spinner: %v\n", err)
				}
			}()

			// Collect results
			var allResults []manager.Package
			
			// Wait for all package managers to respond
			npmResults := <-npmChan
			pipResults := <-pipChan
			scoopResults := <-scoopChan

			// Stop spinner and clear its output
			p.Quit()
			clearScreen()

			// Combine results
			if npmResults != nil {
				allResults = append(allResults, npmResults...)
			}
			if pipResults != nil {
				allResults = append(allResults, pipResults...)
			}
			if scoopResults != nil {
				allResults = append(allResults, scoopResults...)
			}

			if len(allResults) == 0 {
				fmt.Printf("No packages found matching '%s'\n", query)
				return nil
			}

			// Sort results by score
			sort.Slice(allResults, func(i, j int) bool {
				return allResults[i].Score > allResults[j].Score
			})

			// Show results table
			fmt.Printf("\nFound %d packages matching '%s'\n\n", len(allResults), query)
			fmt.Print(renderTable(allResults))
			fmt.Println("\nEnter package number to install (or 'q' to quit): ")

			// Get user selection
			var input string
			fmt.Scanln(&input)

			if input == "q" {
				return nil
			}

			var selection int
			_, err := fmt.Sscanf(input, "%d", &selection)
			if err != nil || selection < 1 || selection > len(allResults) {
				return fmt.Errorf("invalid selection")
			}

			selectedPkg := allResults[selection-1]

			// Show package details
			fmt.Println("\nPackage Details:")
			fmt.Printf("Name: %s\n", titleStyle.Render(selectedPkg.Name))
			fmt.Printf("Provider: %s\n", providerStyle.Render(selectedPkg.Provider))
			fmt.Printf("Version: %s\n", versionStyle.Render(selectedPkg.Version))
			fmt.Printf("Description: %s\n", descStyle.Render(selectedPkg.Description))
			if selectedPkg.Author != "" {
				fmt.Printf("Author: %s\n", selectedPkg.Author)
			}
			if selectedPkg.Homepage != "" {
				fmt.Printf("Homepage: %s\n", selectedPkg.Homepage)
			}
			if selectedPkg.Repository != "" {
				fmt.Printf("Repository: %s\n", selectedPkg.Repository)
			}

			// Ask if user wants to install
			fmt.Print("\nDo you want to install this package? [y/N] ")
			var answer string
			fmt.Scanln(&answer)

			if strings.ToLower(answer) == "y" {
				// Create installation spinner
				s = spinner.New()
				s.Spinner = spinner.Dot
				s.Style = s.Style.Foreground(s.Style.GetForeground())

				m = searchModel{spinner: s, query: selectedPkg.Name}
				p = tea.NewProgram(m)

				// Create channel for installation result
				installErrChan := make(chan error)

				// Start installation in goroutine
				go func() {
					// Find the appropriate package manager
					for _, pm := range mgr.GetManagers() {
						if pm.GetName() == selectedPkg.Provider {
							if err := pm.Install(selectedPkg.Name); err != nil {
								installErrChan <- err
							} else {
								installErrChan <- nil
							}
							break
						}
					}
				}()

				// Start spinner
				go func() {
					if err := p.Start(); err != nil {
						fmt.Printf("Error starting spinner: %v\n", err)
					}
				}()

				// Wait for installation result
				if err := <-installErrChan; err != nil {
					p.Quit()
					return fmt.Errorf("installation failed: %v", err)
				}

				p.Quit()
				fmt.Printf("\n✓ Successfully installed %s\n", titleStyle.Render(selectedPkg.Name))
			}

			return nil
		},
	}

	return cmd
}
