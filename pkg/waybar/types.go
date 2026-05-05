package waybar

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

type Module interface {
	Run(args []string) (WaybarOutput, error)
}
