package updates

import (
	"custom-waybar/pkg/waybar"
	"custom-waybar/utils/ui"
	"fmt"
	"os/exec"
	"strings"
)

type Module struct {
}

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

func (m *Module) Run(args []string) (waybar.WaybarOutput, error) {
	pacmanPkgs := getPackages("checkupdates")
	aurPkgs := getPackages("yay", "-Qum")

	total := len(pacmanPkgs) + len(aurPkgs)
	if total == 0 {
		return waybar.WaybarOutput{}, nil
	}

	var tooltips []string
	if len(pacmanPkgs) > 0 {
		tooltips = append(tooltips, ui.DrawTable("󰮯 Pacman", pacmanPkgs))
	}
	if len(aurPkgs) > 0 {
		tooltips = append(tooltips, ui.DrawTable("󰣇 AUR", aurPkgs))
	}

	out := waybar.WaybarOutput{
		Text:    fmt.Sprintf("%d", total),
		Tooltip: strings.Join(tooltips, "\n"),
	}
	return out, nil
}
