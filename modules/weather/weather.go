package weather

import (
	"custom-waybar/pkg/waybar"
	"custom-waybar/utils/ui"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Module provides weather information from wttr.in.
type Module struct{}

// weatherIcon maps weather description keywords to Nerd Font icons.
func weatherIcon(desc string) string {
	d := strings.ToLower(strings.TrimSpace(desc))
	switch {
	case strings.Contains(d, "thunder") || strings.Contains(d, "lightning"):
		return "" // thunderstorm
	case strings.Contains(d, "drizzle") || strings.Contains(d, "rain") && strings.Contains(d, "light"):
		return "" // light rain
	case strings.Contains(d, "rain") && strings.Contains(d, "heavy"):
		return "" // heavy rain
	case strings.Contains(d, "rain") || strings.Contains(d, "shower"):
		return "󰼳" // rain
	case strings.Contains(d, "snow") || strings.Contains(d, "sleet") || strings.Contains(d, "blizzard"):
		return "󰙿" // snow
	case strings.Contains(d, "fog") || strings.Contains(d, "mist") || strings.Contains(d, "haze"):
		return "󰖑" // fog
	case strings.Contains(d, "cloudy") && strings.Contains(d, "partly"):
		return "" // partly cloudy
	case strings.Contains(d, "cloudy") || strings.Contains(d, "overcast"):
		return "" // cloudy
	case strings.Contains(d, "clear") && (strings.Contains(d, "night") || strings.Contains(d, "evening")):
		return "" // clear night
	case strings.Contains(d, "clear") || strings.Contains(d, "sunny"):
		return "" // sunny
	default:
		return "󰖐" // default/unknown
	}
}

// dataIcon maps wttr.in data line keywords to Nerd Font icons.
func dataIcon(content string) string {
	switch {
	case strings.Contains(content, "°C") || strings.Contains(content, "°F"):
		return ""
	case strings.Contains(content, "km/h") || strings.Contains(content, "mph") || strings.Contains(content, "m/s"):
		return "󰖝"
	case strings.Contains(content, "km") && !strings.Contains(content, "km/h"):
		return "󰈈"
	case strings.Contains(content, "mm") || strings.Contains(content, "inch"):
		return "󰖗"
	case strings.Contains(content, "%") || strings.Contains(content, "humidity"):
		return ""
	case strings.Contains(content, "hPa") || strings.Contains(content, "mbar") || strings.Contains(content, "pressure"):
		return "󰅟"
	default:
		return ""
	}
}

// formatTooltip creates the nicely formatted tooltip from wttr.in output lines.
func formatTooltip(lines []string) string {
	const artWidth = 16
	var city string
	var details []string

	for i, line := range lines {
		if i == 0 {
			city = " " + strings.TrimSpace(line)
			continue
		}
		runes := []rune(line)
		if len(runes) <= artWidth {
			continue
		}
		content := strings.TrimSpace(string(runes[artWidth:]))
		if content == "" || content == "km" || content == "mm" {
			continue
		}
		details = append(details, fmt.Sprintf("%s %s", dataIcon(content), content))
	}

	if len(details) == 0 {
		return city
	}
	return ui.DrawTable(city, details)
}

// parseWttrOutput parses the raw text output from wttr.in and returns
// temperature, unit, icon, and the raw lines for tooltip building.
func parseWttrOutput(body []byte) (temp, unit, icon string, lines []string, err error) {
	raw := string(body)
	lines = strings.Split(raw, "\n")

	if len(lines) < 4 {
		return "", "", "", nil, fmt.Errorf("wttr.in returned too few lines (%d)", len(lines))
	}

	// Line 2 (index 2) has the weather description, line 3 has the temp.
	// The first ~15 chars are ASCII art; data starts at offset 15.
	descRaw := lines[2]
	if len(descRaw) > 15 {
		descRaw = descRaw[15:]
	}
	icon = weatherIcon(strings.TrimSpace(descRaw))

	// Parse temperature line — slice raw line at offset 15 first, then clean.
	tempRaw := lines[3]
	if len(tempRaw) > 15 {
		tempRaw = tempRaw[15:]
	}
	tempFields := strings.Fields(tempRaw)
	if len(tempFields) >= 2 {
		temp = tempFields[0]
		unit = tempFields[1]
	} else if len(tempFields) == 1 {
		temp = tempFields[0]
		unit = "°C"
	} else {
		return "", "", "", nil, fmt.Errorf("could not parse temperature from line: %q", lines[3])
	}

	return temp, unit, icon, lines, nil
}

// Run fetches weather and returns Waybar output.
func (m *Module) Run(args []string) (waybar.WaybarOutput, error) {
	// Default city from args or fallback to Gdansk
	if len(args) == 0 {
		args = []string{"Gdansk"}
	}
	city := args[0]

	url := fmt.Sprintf("https://en.wttr.in/%s?0qnT", city)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return waybar.WaybarOutput{}, fmt.Errorf("weather: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "curl/8.0")

	resp, err := client.Do(req)
	if err != nil {
		return waybar.WaybarOutput{}, fmt.Errorf("weather: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return waybar.WaybarOutput{}, fmt.Errorf("weather: wttr.in returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return waybar.WaybarOutput{}, fmt.Errorf("weather: failed to read response: %w", err)
	}

	temp, unit, icon, lines, err := parseWttrOutput(body)
	if err != nil {
		return waybar.WaybarOutput{}, fmt.Errorf("weather: parse error: %w", err)
	}

	out := waybar.WaybarOutput{
		Text:    fmt.Sprintf("%s  %s%s", icon, temp, unit),
		Tooltip: formatTooltip(lines),
		Alt:     fmt.Sprintf("%s%s", temp, unit),
	}

	return out, nil
}