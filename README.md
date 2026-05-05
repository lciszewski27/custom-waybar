# Waybar go modules
A single, high-performance binary written in Go that handles multiple custom modules for Waybar. All modules feature clean, well-formatted tooltips. I know all of this could be achieved with Bash scripts—but where's the fun in that?

## Usage / Features
- Weather(using wttr.in)
```bash
./custom-waybar weather [city]
```
Example output:
`  20°C`
`╭──  Warszawa ──╮
│ 󰖐 Clear        │
│  20 °C        │
│ 󰖝 ↙ 9 km/h     │
│ 󰈈 10 km        │
│ 󰖗 0.0 mm       │
╰────────────────╯`
Usage:
```json
"custom/weather": {
  "format": "{text}",
  "interval": 1800,
  "return-type": "json",
  "exec": "[path to binary]/custom-waybar weather",
  "on-click": "ghostty -e sh -c 'curl v2.wttr.in/[your city]; echo Done - Press enter to exit; read'",
  "tooltip": true,
}
```

- Updates(working with pacman and yay)
```bash
./custom-waybar updates
```
`
╭────── 󰮯 Pacman ──────╮
│ linux                │
│ ghoustty             │
│ waybar               │
╰──────────────────────╯
╭─────── 󰣇 AUR ────────╮
│ some-think           │
│ other                │
╰──────────────────────╯
`


```json
"custom/pacman": {
    "format": "󰅢  {}",
    "interval": 1800,
    "return-type": "json",
    "exec": "~/.config/waybar/scripts/waybar-modules checkupdates",
    "on-click": "ghostty -e sh -c 'yay -Syu; echo Done - Press enter to exit; read'; pkill -SIGRTMIN+8 waybar",
    "signal": 8,
    "format-empty": ""
}
```
