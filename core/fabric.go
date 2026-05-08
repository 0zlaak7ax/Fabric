// Package core provides the central Fabric application logic,
// including configuration management, pattern loading, and AI vendor integration.
package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	// DefaultConfigDir is the default directory for Fabric configuration files.
	DefaultConfigDir = ".config/fabric"
	// PatternsDir is the subdirectory containing pattern files.
	PatternsDir = "patterns"
	// EnvFile is the name of the environment file storing API keys and settings.
	EnvFile = ".env"
)

// Fabric is the core application struct that holds configuration and state.
type Fabric struct {
	ConfigDir   string
	PatternsDir string
	Vendors     *VendorManager
	Patterns    map[string]*Pattern
}

// NewFabric creates and initializes a new Fabric instance.
// It resolves the configuration directory and prepares the pattern registry.
func NewFabric() (*Fabric, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, DefaultConfigDir)
	patternsDir := filepath.Join(configDir, PatternsDir)

	f := &Fabric{
		ConfigDir:   configDir,
		PatternsDir: patternsDir,
		Patterns:    make(map[string]*Pattern),
		Vendors:     NewVendorManager(),
	}

	return f, nil
}

// Setup ensures the Fabric configuration directory structure exists.
func (f *Fabric) Setup() error {
	dirs := []string{f.ConfigDir, f.PatternsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// LoadPatterns reads all pattern directories from the patterns directory
// and registers them in the Fabric instance.
func (f *Fabric) LoadPatterns() error {
	entries, err := os.ReadDir(f.PatternsDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("patterns directory not found at %s; run 'fabric --setup' first", f.PatternsDir)
		}
		return fmt.Errorf("failed to read patterns directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		patternName := entry.Name()
		patternPath := filepath.Join(f.PatternsDir, patternName)

		pattern, err := LoadPattern(patternName, patternPath)
		if err != nil {
			// Log warning but continue loading other patterns
			fmt.Fprintf(os.Stderr, "warning: skipping pattern %q: %v\n", patternName, err)
			continue
		}
		f.Patterns[patternName] = pattern
	}

	return nil
}

// GetPattern retrieves a loaded pattern by name.
func (f *Fabric) GetPattern(name string) (*Pattern, error) {
	p, ok := f.Patterns[name]
	if !ok {
		return nil, fmt.Errorf("pattern %q not found; use 'fabric --list' to see available patterns", name)
	}
	return p, nil
}

// ListPatterns returns a sorted list of all available pattern names.
func (f *Fabric) ListPatterns() []string {
	names := make([]string, 0, len(f.Patterns))
	for name := range f.Patterns {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
