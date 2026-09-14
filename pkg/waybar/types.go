package waybar

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WaybarOutput represents the JSON output structure for Waybar custom modules.
// See: https://github.com/Alexays/Waybar/wiki/Module:-Custom#return-types
type WaybarOutput struct {
	Text        string `json:"text,omitempty"`
	Alt         string `json:"alt,omitempty"`
	Tooltip     string `json:"tooltip,omitempty"`
	Class       string `json:"class,omitempty"`
	Percentage  int    `json:"percentage,omitempty"`
}

// Module is the interface every waybar custom module must implement.
type Module interface {
	Run(args []string) (WaybarOutput, error)
}

// OutputEmpty returns an empty output (causes Waybar to hide the module).
func OutputEmpty() WaybarOutput {
	return WaybarOutput{}
}

// OutputSimple creates a basic output with just text.
func OutputSimple(text string) WaybarOutput {
	return WaybarOutput{Text: text}
}

// OutputTooltip creates output with text and a tooltip.
func OutputTooltip(text, tooltip string) WaybarOutput {
	return WaybarOutput{
		Text:    text,
		Tooltip: tooltip,
	}
}

// OutputClass creates output with text and a CSS class for styling.
func OutputClass(text, class string) WaybarOutput {
	return WaybarOutput{
		Text:  text,
		Class: class,
	}
}

// OutputFull creates a fully specified output.
func OutputFull(text, alt, tooltip, class string, percentage int) WaybarOutput {
	return WaybarOutput{
		Text:       text,
		Alt:        alt,
		Tooltip:    tooltip,
		Class:      class,
		Percentage: percentage,
	}
}

// MustMarshalJSON marshals the output to JSON, panicking on error.
// Use in main.go when the module already handled errors.
func (o WaybarOutput) MustMarshalJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(fmt.Sprintf("waybar: failed to marshal output: %v", err))
	}
	return data
}

// MergeArgs merges default args with override args from the command line.
// Defaults come first; any arg from cmdline that is non-empty overrides the default.
func MergeArgs(defaults, overrides []string) []string {
	if len(overrides) == 0 {
		return defaults
	}
	merged := make([]string, len(defaults))
	copy(merged, defaults)
	for i, o := range overrides {
		if o != "" && i < len(merged) {
			merged[i] = o
		} else if o != "" {
			merged = append(merged, o)
		}
	}
	return merged
}

// FormatFlags parses command-line flags of the form --key=value or --key value
// into a map. Remaining positional args are returned separately.
func FormatFlags(args []string) (flags map[string]string, positional []string) {
	flags = make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "--") {
			positional = append(positional, arg)
			continue
		}
		trim := strings.TrimPrefix(arg, "--")
		if k, v, ok := strings.Cut(trim, "="); ok {
			flags[k] = v
		} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			flags[trim] = args[i+1]
			i++
		} else {
			flags[trim] = "true"
		}
	}
	return flags, positional
}