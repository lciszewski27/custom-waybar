package updates

import (
	"custom-waybar/pkg/waybar"
	"custom-waybar/utils/ui"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Module monitors pending system updates.
type Module struct{}

// cache holds update results for a configurable TTL to avoid hammering pkg managers.
var (
	cacheStore   []updateGroup
	cacheTime    time.Time
	cacheTTL     = 5 * time.Minute
	cacheMu      sync.Mutex
)

type updateGroup struct {
	Title    string   // e.g. "󰮯 Pacman"
	Packages []string
}

// getPackages runs a command and returns the first field of each output line.
func getPackages(command string, args ...string) []string {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return []string{}
	}
	lines := strings.Split(trimmed, "\n")
	pkgs := make([]string, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			pkgs = append(pkgs, fields[0])
		}
	}
	return pkgs
}

// refreshCache runs all package managers and caches the results.
func refreshCache() {
	pacmanPkgs := getPackages("checkupdates")
	aurPkgs := getPackages("yay", "-Qum")
	flatpakPkgs := getPackages("flatpak", "remote-ls", "--updates")

	var groups []updateGroup
	if len(pacmanPkgs) > 0 {
		groups = append(groups, updateGroup{Title: "󰮯 Pacman", Packages: pacmanPkgs})
	}
	if len(aurPkgs) > 0 {
		groups = append(groups, updateGroup{Title: "󰣇 AUR", Packages: aurPkgs})
	}
	if len(flatpakPkgs) > 0 {
		groups = append(groups, updateGroup{Title: "󰏖 Flatpak", Packages: flatpakPkgs})
	}
	cacheMu.Lock()
	cacheStore = groups
	cacheTime = time.Now()
	cacheMu.Unlock()
}

// getCachedUpdates returns cached results, refreshing if TTL expired.
func getCachedUpdates() []updateGroup {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if len(cacheStore) == 0 || time.Since(cacheTime) > cacheTTL {
		// Release lock while refreshing to avoid deadlock
		cacheMu.Unlock()
		refreshCache()
		cacheMu.Lock()
	}
	return cacheStore
}

// countUpdates returns total update count across all groups.
func countUpdates(groups []updateGroup) int {
	total := 0
	for _, g := range groups {
		total += len(g.Packages)
	}
	return total
}

// buildTooltip renders all update groups as a single tooltip.
func buildTooltip(groups []updateGroup) string {
	if len(groups) == 0 {
		return "Up to date"
	}
	var parts []string
	for _, g := range groups {
		parts = append(parts, ui.DrawTable(g.Title, g.Packages))
	}
	return strings.Join(parts, "\n")
}

// buildAltText returns a summary like "P:5 A:2 F:1".
func buildAltText(groups []updateGroup) string {
	var tags []string
	for _, g := range groups {
		var prefix string
		switch {
		case strings.Contains(g.Title, "Pacman"):
			prefix = "P"
		case strings.Contains(g.Title, "AUR"):
			prefix = "A"
		case strings.Contains(g.Title, "Flatpak"):
			prefix = "F"
		}
		tags = append(tags, fmt.Sprintf("%s:%d", prefix, len(g.Packages)))
	}
	return strings.Join(tags, " ")
}

// buildClass returns a CSS class based on update severity.
func buildClass(total int) string {
	switch {
	case total == 0:
		return "none"
	case total < 5:
		return "few"
	case total < 15:
		return "some"
	default:
		return "many"
	}
}

// Run checks for available system updates.
func (m *Module) Run(args []string) (waybar.WaybarOutput, error) {
	// Parse optional --ttl flag
	flags, _ := waybar.FormatFlags(args)
	if ttlStr, ok := flags["ttl"]; ok {
		if d, err := time.ParseDuration(ttlStr); err == nil && d > 0 {
			cacheMu.Lock()
			cacheTTL = d
			cacheMu.Unlock()
		}
	}

	groups := getCachedUpdates()
	total := countUpdates(groups)

	if total == 0 {
		return waybar.OutputEmpty(), nil
	}

	return waybar.OutputFull(
		fmt.Sprintf("%d", total),
		buildAltText(groups),
		buildTooltip(groups),
		buildClass(total),
		0,
	), nil
}