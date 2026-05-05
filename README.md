# 🚀 Waybar Go Modules

[![Go Version](https://img.shields.io/github/go-mod/go-version/lciszewski27/custom-waybar)](https://golang.org)
[![License](https://img.shields.io/github/license/lciszewski27/custom-waybar)](LICENSE)

A single, high-performance binary written in Go that handles multiple custom modules for Waybar. All modules feature clean, well-formatted tooltips.

> *Sure, you could do all of this with Bash scripts—but where's the fun in that?* 😉

---

## ✨ Features

### ☁️ Weather (via wttr.in)
Fetches real-time weather data with a beautifully framed tooltip.

**CLI Usage:**
```bash
./custom-waybar weather [city]
```

**Tooltip Preview:**
```text
╭──  Warsaw ───╮
│ 󰖐 Clear       │
│  20 °C       │
│ 󰖝 ↙ 9 km/h    │
│ 󰈈 10 km       │
│ 󰖗 0.0 mm      │
╰───────────────╯
```

**Waybar Configuration (`config.jsonc`):**
```jsonc
"custom/weather": {
    "format": "{}",
    "interval": 1800,
    "return-type": "json",
    "exec": "/path/to/binary/custom-waybar weather [city]",
    "on-click": "ghostty -e sh -c 'curl v2.wttr.in/[your_city]; echo Done - Press enter to exit; read'",
    "tooltip": true
}
```

---

### 📦 Updates (Pacman & AUR)
Monitors pending updates for both official repositories and the AUR (via `yay`).

**CLI Usage:**
```bash
./custom-waybar updates
```

**Tooltip Preview:**
```text
╭────── 󰮯 Pacman ──────╮
│ linux                │
│ ghostty              │
│ waybar               │
╰──────────────────────╯
╭─────── 󰣇 AUR ────────╮
│ some-package         │
│ other-app            │
╰──────────────────────╯
```

**Waybar Configuration (`config.jsonc`):**
```jsonc
"custom/pacman": {
    "format": "󰅢  {}",
    "interval": 1800,
    "return-type": "json",
    "exec": "/path/to/binary/custom-waybar updates",
    "on-click": "ghostty -e sh -c 'yay -Syu; pkill -SIGRTMIN+8 waybar'",
    "signal": 8,
    "format-empty": ""
}
```

---

## 🛠️ Installation & Building

Ensure you have **Go 1.2x+** installed.

1. **Clone the repository:**
   ```bash
   git clone [https://github.com/lciszewski27/custom-waybar.git](https://github.com/lciszewski27/custom-waybar.git)
   cd custom-waybar
   ```

2. **Build the binary:**
   ```bash
   go build -o custom-waybar main.go
   ```

3. **Move to your config folder:**
   ```bash
   cp custom-waybar ~/.config/waybar/scripts/
   ```
