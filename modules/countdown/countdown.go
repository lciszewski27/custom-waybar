package countdown

import (
	"custom-waybar/pkg/waybar"
	"custom-waybar/utils/ui"
	"fmt"
	"math"
	"strings"
	"time"
)

// TimeSpan represents a human-friendly breakdown of a duration.
type TimeSpan struct {
	Weeks    int
	Days     int
	Hours    int
	Minutes  int
	Seconds  int
	Negative bool
}

func splitDuration(d time.Duration) TimeSpan {
	if d < 0 {
		ts := splitDuration(-d)
		ts.Negative = true
		return ts
	}
	totalSec := int(math.Round(d.Seconds()))
	if totalSec < 0 {
		totalSec = 0
	}
	const (
		secPerMin = 60
		secPerHr  = 3600
		secPerDay = 86400
		secPerWk  = 604800
	)
	weeks := totalSec / secPerWk
	rem := totalSec % secPerWk
	days := rem / secPerDay
	rem = rem % secPerDay
	hours := rem / secPerHr
	rem = rem % secPerHr
	mins := rem / secPerMin
	secs := rem % secPerMin
	return TimeSpan{Weeks: weeks, Days: days, Hours: hours, Minutes: mins, Seconds: secs}
}

// ─── Format tokens ───────────────────────────────────────────────────────────

type formatFunc func(ts TimeSpan) string

var tokenMap = map[string]formatFunc{
	"w":  func(ts TimeSpan) string { return fmt.Sprintf("%d", ts.Weeks) },
	"ww": func(ts TimeSpan) string { return fmt.Sprintf("%02d", ts.Weeks) },
	"d":  func(ts TimeSpan) string { return fmt.Sprintf("%d", ts.Days) },
	"dd": func(ts TimeSpan) string { return fmt.Sprintf("%02d", ts.Days) },
	"H":  func(ts TimeSpan) string { return fmt.Sprintf("%d", ts.Hours) },
	"HH": func(ts TimeSpan) string { return fmt.Sprintf("%02d", ts.Hours) },
	"M":  func(ts TimeSpan) string { return fmt.Sprintf("%d", ts.Minutes) },
	"MM": func(ts TimeSpan) string { return fmt.Sprintf("%02d", ts.Minutes) },
	"S":  func(ts TimeSpan) string { return fmt.Sprintf("%d", ts.Seconds) },
	"SS": func(ts TimeSpan) string { return fmt.Sprintf("%02d", ts.Seconds) },
	"total-d": func(ts TimeSpan) string {
		return fmt.Sprintf("%.1f", ts.TotalSeconds()/86400)
	},
	"total-h": func(ts TimeSpan) string {
		return fmt.Sprintf("%.1f", ts.TotalSeconds()/3600)
	},
	"total-m": func(ts TimeSpan) string {
		return fmt.Sprintf("%.0f", ts.TotalSeconds()/60)
	},
	"total-s": func(ts TimeSpan) string {
		return fmt.Sprintf("%.0f", ts.TotalSeconds())
	},
	"label-w": labelConst("week"), "label-d": labelConst("day"),
	"label-h": labelConst("hour"), "label-m": labelConst("minute"), "label-s": labelConst("second"),
	"label-ws": labelPlural("week"), "label-ds": labelPlural("day"),
	"label-hs": labelPlural("hour"), "label-ms": labelPlural("minute"), "label-ss": labelPlural("second"),
	"status": func(ts TimeSpan) string {
		if ts.Negative { return "overdue" }
		if ts.Weeks == 0 && ts.Days == 0 && ts.Hours == 0 && ts.Minutes == 0 && ts.Seconds == 0 { return "done" }
		return "active"
	},
}

func labelConst(s string) formatFunc {
	return func(ts TimeSpan) string { return s }
}

func labelPlural(singular string) formatFunc {
	return func(ts TimeSpan) string {
		return Pluralize(singular, ts.ValueForUnit(singular[0]))
	}
}

func (ts TimeSpan) ValueForUnit(unit byte) int {
	switch unit {
	case 'w': return ts.Weeks
	case 'd': return ts.Days
	case 'h': return ts.Hours
	case 'm': return ts.Minutes
	case 's': return ts.Seconds
	}
	return 0
}

func (ts TimeSpan) TotalSeconds() float64 {
	return float64(ts.Weeks)*604800 + float64(ts.Days)*86400 +
		float64(ts.Hours)*3600 + float64(ts.Minutes)*60 + float64(ts.Seconds)
}

func Pluralize(singular string, n int) string {
	if n == 1 { return singular }
	return singular + "s"
}

// ApplyFormat renders a TimeSpan using the given format string.
// Tokens in {} braces are replaced. Use {{ and }} for literal braces.
// Examples: "{d}d {H}h {M}m" -> "5d 12h 30m", "{dd}:{HH}:{MM}:{SS}" -> "05:12:30:45"
func ApplyFormat(format string, ts TimeSpan) string {
	if format == "" {
		format = "{d}d {H}h {M}m {S}s"
	}
	var sb strings.Builder
	i := 0
	for i < len(format) {
		if format[i] == '{' {
			if i+1 < len(format) && format[i+1] == '{' {
				sb.WriteByte('{'); i += 2; continue
			}
			j := i + 1
			for j < len(format) && format[j] != '}' { j++ }
			if j >= len(format) { sb.WriteByte(format[i]); i++; continue }
			token := format[i+1 : j]
			if fn, ok := tokenMap[token]; ok {
				sb.WriteString(fn(ts))
			} else {
				sb.WriteString("{" + token + "}")
			}
			i = j + 1
		} else if format[i] == '}' && i+1 < len(format) && format[i+1] == '}' {
			sb.WriteByte('}'); i += 2
		} else {
			sb.WriteByte(format[i]); i++
		}
	}
	return sb.String()
}

// buildTooltip creates a nice visual breakdown for the tooltip.
func buildTooltip(ts TimeSpan, target time.Time, label string) string {
	var rows []string
	if label != "" {
		rows = append(rows, fmt.Sprintf("\U000F015B %s", label))
	}
	if ts.Negative {
		rows = append(rows, "\u26A0 Overdue")
		return ui.DrawTable("\U000F00F0 Countdown", rows)
	}
	var parts []string
	add := func(n int, icon, singular string) {
		if n > 0 {
			parts = append(parts, fmt.Sprintf("%s %d %s", icon, n, Pluralize(singular, n)))
		}
	}
	add(ts.Weeks, "\U000F0035", "week")
	add(ts.Days, "\U000F07DF", "day")
	add(ts.Hours, "\U000F0954", "hour")
	add(ts.Minutes, "\U000F02DA", "minute")
	add(ts.Seconds, "\U000F06EF", "second")
	if len(parts) == 0 {
		parts = append(parts, "\U000F0134 Done!")
	}
	rows = append(rows, parts...)
	rows = append(rows, "")
	rows = append(rows, fmt.Sprintf("\U000F0150 %s", target.Format("Mon 2006-01-02 15:04:05")))
	return ui.DrawTable("\U000F00F0 Countdown", rows)
}

const DefaultFormat = "{d}d {H}h {M}m {S}s"

var Presets = map[string]string{
	"compact":  "{dd}:{HH}:{MM}:{SS}",
	"default":  DefaultFormat,
	"long":     "{w} {label-ws}, {d} {label-ds}, {H} {label-hs}, {M} {label-ms}",
	"clock":    "{HH}:{MM}:{SS}",
	"days-hrs": "{d}d {H}h",
	"weeks":    "{w}w {d}d",
	"total-hr": "{total-h}h",
	"total-d":  "{total-d}d",
}

// ─── Parsing target time ─────────────────────────────────────────────────────

func parseTarget(s string) (time.Time, error) {
	if s == "" { return time.Time{}, fmt.Errorf("empty target time") }
	if strings.HasPrefix(s, "+") { return parseRelativeDuration(s[1:]) }
	if t, err := time.Parse(time.RFC3339, s); err == nil { return t, nil }
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil { return t, nil }
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil { return t, nil }
	if t, err := time.Parse("2006-01-02 15:04", s); err == nil { return t, nil }
	if t, err := time.Parse("2006-01-02", s); err == nil { return t, nil }
	weekday := parseWeekday(s)
	if weekday >= 0 { return nextWeekday(weekday, 0, 0), nil }
	if parts := strings.Fields(s); len(parts) == 2 {
		weekday = parseWeekday(parts[0])
		if weekday >= 0 {
			var hour, min int
			if _, err := fmt.Sscanf(parts[1], "%d:%d", &hour, &min); err == nil {
				return nextWeekday(weekday, hour, min), nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("could not parse target time: %q", s)
}

func parseRelativeDuration(s string) (time.Time, error) {
	d, err := parseDurationV2(s)
	if err != nil { return time.Time{}, err }
	return time.Now().Add(d), nil
}

func parseDurationV2(s string) (time.Duration, error) {
	var total time.Duration
	numStr := ""
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch >= '0' && ch <= '9') || ch == '.' {
			numStr += string(ch)
			continue
		}
		if numStr == "" { return 0, fmt.Errorf("invalid duration: %q", s) }
		switch ch {
		case 'w':
			total += parseNumber(numStr) * 7 * 24 * time.Hour
			numStr = ""
		case 'd':
			total += parseNumber(numStr) * 24 * time.Hour
			numStr = ""
		default:
			part := numStr + string(ch)
			j := i + 1
			for j < len(s) && s[j] >= 'a' && s[j] <= 'z' { part += string(s[j]); j++ }
			d, err := time.ParseDuration(part)
			if err != nil { return 0, fmt.Errorf("invalid duration part %q in %q", part, s) }
			total += d
			numStr = ""
			i = j - 1
		}
	}
	if numStr != "" { return 0, fmt.Errorf("trailing number in duration: %q", s) }
	return total, nil
}

func parseNumber(s string) time.Duration {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return time.Duration(f)
}

func parseWeekday(s string) int {
	lower := strings.ToLower(s)
	mapping := map[string]time.Weekday{
		"sunday": time.Sunday, "sun": time.Sunday,
		"monday": time.Monday, "mon": time.Monday,
		"tuesday": time.Tuesday, "tue": time.Tuesday, "tues": time.Tuesday,
		"wednesday": time.Wednesday, "wed": time.Wednesday,
		"thursday": time.Thursday, "thu": time.Thursday, "thur": time.Thursday,
		"friday": time.Friday, "fri": time.Friday,
		"saturday": time.Saturday, "sat": time.Saturday,
	}
	if wd, ok := mapping[lower]; ok { return int(wd) }
	return -1
}

func nextWeekday(targetWeekday int, hour, min int) time.Time {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	diff := (targetWeekday - int(target.Weekday()) + 7) % 7
	if diff == 0 && (target.Before(now) || target.Equal(now)) { diff = 7 }
	return target.AddDate(0, 0, diff)
}

// ─── Module ──────────────────────────────────────────────────────────────────

type Module struct{}

func (m *Module) Run(args []string) (waybar.WaybarOutput, error) {
	flags, positional := waybar.FormatFlags(args)
	targetStr := flags["target"]
	if targetStr == "" && len(positional) > 0 { targetStr = positional[0] }
	if targetStr == "" {
		return waybar.WaybarOutput{}, fmt.Errorf("countdown: no target specified; use --target=<time>")
	}
	target, err := parseTarget(targetStr)
	if err != nil { return waybar.WaybarOutput{}, fmt.Errorf("countdown: %w", err) }
	format := flags["format"]
	if format == "" { format = flags["f"] }
	if format == "" { format = DefaultFormat } else if preset, ok := Presets[format]; ok { format = preset }
	label := flags["label"]
	if label == "" { label = flags["l"] }
	class := flags["class"]
	if class == "" { class = flags["c"] }
	now := time.Now()
	remaining := target.Sub(now)
	ts := splitDuration(remaining)
	text := ApplyFormat(format, ts)
	if class == "" {
		switch {
		case ts.Negative: class = "overdue"
		case ts.Weeks == 0 && ts.Days == 0 && ts.Hours == 0 && ts.Minutes < 60: class = "soon"
		default: class = "active"
		}
	}
	alt := ApplyFormat("{d}d {H}h {M}m", ts)
	out := waybar.OutputFull(text, alt, buildTooltip(ts, target, label), class, 0)
	return out, nil
}
