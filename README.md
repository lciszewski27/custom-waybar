# 🚀 Waybar Go Modules

[![Go Version](https://img.shields.io/github/go-mod/go-version/lciszewski27/custom-waybar)](https://golang.org)
[![License](https://img.shields.io/github/license/lciszewski27/custom-waybar)](LICENSE)

A single, high-performance binary written in Go that handles multiple custom modules for Waybar. All modules feature clean, well-formatted tooltips using the built-in `ui.DrawTable` utility.

> _Sure, you could do all of this with Bash scripts—but where's the fun in that?_ 😉

---

## ✨ Features

### ☁️ Weather (via wttr.in)

Fetches real-time weather data with a beautifully framed tooltip. Supports multiple weather conditions with Nerd Font icons.

**CLI Usage:**

```bash
./custom-waybar weather [city]
```

_Default city:_ Gdansk

**Tooltip Preview:**

```
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
    "exec": "/path/to/custom-waybar weather Warsaw",
    "on-click": "ghostty -e sh -c 'curl v2.wttr.in/Warsaw; echo Done - Press enter to exit; read'",
    "tooltip": true
}
```

**CSS Classes:** Uses `alt` field with temperature value.

---

### 📦 Updates (Pacman, AUR & Flatpak)

Monitors pending updates from official repositories, AUR (via `yay`), and Flatpak.

**CLI Usage:**

```bash
./custom-waybar updates [--ttl=5m]
```

**Tooltip Preview:**

```
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
    "format": "󰅢 {}",
    "interval": 1800,
    "return-type": "json",
    "exec": "/path/to/custom-waybar updates",
    "on-click": "ghostty -e sh -c 'yay -Syu; pkill -SIGRTMIN+8 waybar'",
    "signal": 8,
    "format-empty": ""
}
```

**CSS Classes:** `none`, `few`, `some`, `many` — based on total update count.

---

### ⏳ Countdown

Countdown to a target date/time with customizable output format. Supports absolute dates, relative durations, and recurring weekdays.

**CLI Usage:**

```bash
# Absolute date/time
./custom-waybar countdown --target="2025-12-25 00:00" --label="Christmas"

# Relative duration
./custom-waybar countdown --target="+2h30m" --format="{HH}:{MM}:{SS}"

# Named weekday (next occurrence)
./custom-waybar countdown --target="friday 17:00" --format="clock" --label="Weekend"

# Full custom format
./custom-waybar countdown --target="2026-01-01T00:00:00" --format="{d}d {H}h {M}m {S}s" --label="New Year"
```

**Format Tokens:**

| Token        | Description                    | Example  |
| ------------ | ------------------------------ | -------- |
| `{w}`        | Weeks (no padding)             | `2`      |
| `{ww}`       | Weeks (zero-padded)            | `02`     |
| `{d}`        | Days (no padding)              | `5`      |
| `{dd}`       | Days (zero-padded)             | `05`     |
| `{H}`        | Hours (no padding)             | `12`     |
| `{HH}`       | Hours (zero-padded)            | `12`     |
| `{M}`        | Minutes (no padding)           | `30`     |
| `{MM}`       | Minutes (zero-padded)          | `30`     |
| `{S}`        | Seconds (no padding)           | `45`     |
| `{SS}`       | Seconds (zero-padded)          | `45`     |
| `{label-w}`  | "week" (always singular)       | `week`   |
| `{label-ws}` | "week"/"weeks" (plural-aware)  | `weeks`  |
| `{total-d}`  | Total days as decimal          | `7.5`    |
| `{total-h}`  | Total hours as decimal         | `180.0`  |
| `{total-m}`  | Total minutes (rounded)        | `10800`  |
| `{total-s}`  | Total seconds (rounded)        | `648000` |
| `{status}`   | "active", "done", or "overdue" | `active` |

**Preset Formats:**

| Preset     | Format                                | Example Output       |
| ---------- | ------------------------------------- | -------------------- |
| `default`  | `{d}d {H}h {M}m {S}s`                 | `5d 12h 30m 45s`     |
| `compact`  | `{dd}:{HH}:{MM}:{SS}`                 | `05:12:30:45`        |
| `clock`    | `{HH}:{MM}:{SS}`                      | `12:30:45`           |
| `long`     | `{w} {label-ws}, {d} {label-ds}, ...` | `2 weeks, 5 days...` |
| `days-hrs` | `{d}d {H}h`                           | `5d 12h`             |
| `weeks`    | `{w}w {d}d`                           | `2w 5d`              |
| `total-hr` | `{total-h}h`                          | `180.0h`             |
| `total-d`  | `{total-d}d`                          | `7.5d`               |

**Tooltip Preview:**

```
╭─────── 󰃰 Countdown ────────╮
│ 󰔛 Christmas                │
│ 󰀵 2 weeks                  │
│ 󰧟 5 days                   │
│ 󰥔 12 hours                 │
│ 󰋚 30 minutes               │
│ 󱫯 45 seconds               │
│                            │
│ 󰅐 Thu 2025-12-25 00:00:00  │
╰────────────────────────────╯
```

**CSS Classes:** `active` (normal), `soon` (< 60 min), `overdue` (past).

**Waybar Configuration (`config.jsonc`):**

```jsonc
"custom/countdown": {
    "format": "{}",
    "interval": 1,
    "return-type": "json",
    "exec": "/path/to/custom-waybar countdown --target=\"2025-12-25 00:00\" --label=\"Christmas\" --format=default",
    "tooltip": true
}
```

---

## 🛠️ Installation & Building

Ensure you have **Go 1.2x+** installed.

1. **Clone the repository:**

   ```bash
   git clone https://github.com/lciszewski27/custom-waybar.git
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

4. **Add CSS classes** (optional, for styling):
   ```css
   /* weather */
   #custom-weather {
     color: #89b4fa;
   }

   /* updates */
   #custom-pacman {
     color: #a6e3a1;
   }
   #custom-pacman.none {
     color: #585b70;
   }
   #custom-pacman.few {
     color: #a6e3a1;
   }
   #custom-pacman.some {
     color: #f9e2af;
   }
   #custom-pacman.many {
     color: #f38ba8;
   }

   /* countdown */
   #custom-countdown {
     color: #cba6f7;
   }
   #custom-countdown.soon {
     color: #f9e2af;
   }
   #custom-countdown.overdue {
     color: #f38ba8;
   }
   ```

## 🗒 TODO

- Stock / Crypto prices
- Docker container status
- Uptime
- Pomodoro timer
- World clocks
- Calendar events (ical)
