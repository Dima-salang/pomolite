# PomoLite

<div align="center">

```text
██████╗  ██████╗ ███╗   ███╗ ██████╗ ██╗     ██╗████████╗███████╗
██╔══██╗██╔═══██╗████╗ ████║██╔═══██╗██║     ██║╚══██╔══╝██╔════╝
██████╔╝██║   ██║██╔████╔██║██║   ██║██║     ██║   ██║   █████╗  
██╔═══╝ ██║   ██║██║╚██╔╝██║██║   ██║██║     ██║   ██║   ██╔══╝  
██║     ╚██████╔╝██║ ╚═╝ ██║╚██████╔╝███████╗██║   ██║   ███████╗
╚═╝      ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚══════╝
```

**A high-performance, lightweight CLI Pomodoro application designed for maximum productivity and minimal distraction.**

[![Go Report Card](https://goreportcard.com/badge/github.com/Dima-salang/pomolite)](https://goreportcard.com/report/github.com/Dima-salang/pomolite)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Dima-salang/pomolite)](https://go.dev/)

</div>

---

## 🚀 Overview

**PomoLite** is a minimalist yet powerful Pomodoro timer built for the terminal. Leveraging the power of **Go** and the **Bubble Tea** TUI framework, it provides an immersive, distraction-free environment to help you stay focused on your deep work.

Whether you're a software engineer, student, or creative professional, PomoLite helps you manage your time effectively and gain insights into your productivity patterns with built-in persistence and analytics.

## ✨ Key Features

- 🕒 **Interactive TUI**: A beautiful, real-time timer with a dynamic progress bar.
- 📊 **Persistence & Analytics**: Automatic session logging to a local SQLite database.
- 📋 **Flexible Sessions**: Start sessions with custom labels and durations directly from the CLI.
- 📈 **Productivity Insights**: Comprehensive statistics (sessions count, total duration, averages) for any timeframe.
- 🔔 **Desktop Notifications**: Native alerts ensure you never miss a transition between work and break.
- 🖱️ **Intuitive Controls**: Keyboard-centric interface for pausing, resuming, and navigating.

---

## 🛠️ Installation

### 1. Automatic Installation (Recommended)
Our installation script handles building and setting up the binary:
```bash
chmod +x scripts/install.sh
./scripts/install.sh
```

### 2. Using `go install`
```bash
go install github.com/Dima-salang/pomolite/cmd/pomo@latest
```

### 3. Manual Build
```bash
git clone https://github.com/Dima-salang/pomolite.git
cd pomolite
go build -o pomo cmd/pomo/main.go
```

---

## 📖 Usage Guide

PomoLite is designed to be flexible. You can use the interactive dashboard or trigger sessions directly with flags.

### Interactive Dashboard
Launch the full TUI experience to manage your sessions, history, and stats.
```bash
pomo
```

### Direct Session Start
Launch a specific session immediately:
```bash
pomo start -m 25 -b 5 -l "Deep Work"
```
**Available Flags:**
- `-m, --minutes`: Duration of the work session (default: 30)
- `-b, --break`: Duration of the break (default: 5)
- `-l, --label`: A descriptive label for the session (default: "Work")

### In-Timer Controls
- `p` or `Space`: Toggle Pause/Resume
- `r`: Resume/Restart timer
- `q` or `Esc`: Quit and save session progress
- `ctrl+c`: Terminate without saving

### Session History
List your past work sessions:
```bash
pomo sessions --limit 10
```

### Advanced Statistics
Analyze your focus patterns across different timeframes:
```bash
pomo stat --timeframe week
```
*Timeframe options: `today`, `week`, `month`, `year`, `all`.*

---

## 🏗️ Architecture & Stack

PomoLite is built with a modern Go stack following the **Model-View-Update (MVU)** pattern:

- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (TUI engine)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) (Terminal CSS)
- **CLI Logic**: [Cobra](https://github.com/spf13/cobra) (Command-line parsing)
- **Storage**: SQLite3 (Local persistence)
- **Notifications**: [Beeep](https://github.com/gen2brain/beeep) (Cross-platform alerts)

---

## 🤝 Contributing

Contributions make the open-source community an amazing place!
1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📜 License

This project is licensed under the terms of the [LICENSE](LICENSE) file.

---

<div align="center">
  Developed by <b>PUTAN LUIS GABRIELLE</b><br>
  <a href="mailto:luisgabrielle1026@gmail.com">luisgabrielle1026@gmail.com</a>
</div>