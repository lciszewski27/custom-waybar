package weather

import (
	"custom-waybar/pkg/waybar"
	"custom-waybar/utils/ui"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Module struct {
}

func getWeatherCondition(weatherDesc string) string {
	desc := strings.ToLower(strings.TrimSpace(weatherDesc))
	switch {
	case strings.Contains(desc, "clear"), strings.Contains(desc, "sunny"):
		return ""
	case strings.Contains(desc, "cloudy"):
		return ""
	case strings.Contains(desc, "rain"):
		return "󰼳"
	case strings.Contains(desc, "snow"):
		return "󰙿"
	case strings.Contains(desc, "thunder"):
		return ""
	default:
		return "󰖐"
	}
}

func formatWeatherTooltip(lines []string) string {
	const artWidth = 16
	var city string
	var details []string

	for i, line := range lines {
		if i == 0 {
			city = " " + strings.TrimSpace(line)
			continue
		}

		runes := []rune(line)
		if len(runes) > artWidth {
			content := strings.TrimSpace(string(runes[artWidth:]))
			if content == "" || content == "km" || content == "mm" {
				continue
			}

			var icon string
			switch {
			case strings.Contains(content, "°C"):
				icon = ""
			case strings.Contains(content, "km/h"):
				icon = "󰖝"
			case strings.Contains(content, "km"):
				icon = "󰈈"
			case strings.Contains(content, "mm"):
				icon = "󰖗"
			default:
				icon = "󰖐"
			}
			details = append(details, fmt.Sprintf("%s %s", icon, content))
		}
	}
	return ui.DrawTable(city, details)
}

func (m *Module) Run(args []string) (waybar.WaybarOutput, error) {
	url := "https://en.wttr.in/Gdansk?0qnT"
	if len(args) != 0 {
		url = "https://en.wttr.in/" + args[0] + "?0qnT"
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "curl")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return waybar.WaybarOutput{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	lines := strings.Split(string(body), "\n")

	if len(lines) < 4 {
		return waybar.WaybarOutput{}, err
	}

	temp := strings.Fields(lines[3][15:])[0]
	unit := strings.Fields(lines[3][15:])[1]
	icon := getWeatherCondition(strings.Join(strings.Fields(lines[2][15:]), " "))

	out := waybar.WaybarOutput{
		Text:    fmt.Sprintf("%s  %s%s", icon, temp, unit),
		Tooltip: formatWeatherTooltip(lines),
	}
	return out, nil
}
