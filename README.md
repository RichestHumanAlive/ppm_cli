# Panoramic Package Manager (PPM)

A platform-agnostic CLI tool that unifies package management across multiple package managers like npm, pip, scoop, and winget.

## Features

- Unified command interface for multiple package managers
- Cross-platform support (Windows, macOS, Linux)
- Advanced package search across multiple package managers
- Environment export/import functionality
- Parallel operations support
- Interactive UI with Gum library

## Installation

```bash
go install github.com/yourusername/ppm@latest
```

## Usage

```bash
# Install a package
ppm install <package-name>

# Search for a package
ppm search <package-name>

# Update packages
ppm update [package-name]

# Remove a package
ppm remove <package-name>
```

## Development

### Prerequisites

- Go 1.21 or higher
- Git

### Building from source

```bash
git clone https://github.com/yourusername/ppm.git
cd ppm
go build
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
