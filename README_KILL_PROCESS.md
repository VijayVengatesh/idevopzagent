# How to Kill iDevopzAgent Process

This document explains how to manually stop the iDevopzAgent process on different operating systems.

## Windows

### Method 1: Command Line (Recommended)
```cmd
taskkill /F /IM "metrics-agent-windows-amd64.exe"
```

### Method 2: Find Process First
```cmd
# Find the process
tasklist | findstr "metrics-agent"

# Kill by PID (replace XXXX with actual PID)
taskkill /F /PID XXXX
```

### Method 3: Task Manager
1. Press `Ctrl+Shift+Esc` to open Task Manager
2. Find "metrics-agent-windows-amd64.exe" in Processes tab
3. Right-click and select "End Task"

### Method 4: PowerShell
```powershell
Get-Process -Name "*metrics-agent*" | Stop-Process -Force
```

---

## Linux

### Method 1: Kill by Name (Recommended)
```bash
pkill -f metrics-agent
```

### Method 2: Find and Kill by PID
```bash
# Find the process
ps aux | grep metrics-agent

# Kill by PID (replace XXXX with actual PID)
kill -9 XXXX
```

### Method 3: Kill All Related Processes
```bash
killall metrics-agent-linux-amd64
```

### Method 4: Using System Monitor
1. Open System Monitor (Activity Monitor)
2. Search for "metrics-agent"
3. Select and click "End Process"

---

## macOS

### Method 1: Kill by Name (Recommended)
```bash
pkill -f metrics-agent
```

### Method 2: Find and Kill by PID
```bash
# Find the process
ps aux | grep metrics-agent

# Kill by PID (replace XXXX with actual PID)
kill -9 XXXX
```

### Method 3: Activity Monitor
1. Open Activity Monitor (Applications > Utilities)
2. Search for "metrics-agent"
3. Select the process and click "Force Quit"

### Method 4: Force Quit Menu
1. Press `Cmd+Option+Esc`
2. Select the metrics-agent process
3. Click "Force Quit"

---

## Quick Kill Scripts

### Windows (save as kill_agent.bat)
```batch
@echo off
taskkill /F /IM "metrics-agent-windows-amd64.exe" 2>nul
if %errorlevel% == 0 (
    echo Process killed successfully
) else (
    echo Process not found or already stopped
)
pause
```

### Linux/macOS (save as kill_agent.sh)
```bash
#!/bin/bash
pkill -f metrics-agent
if [ $? -eq 0 ]; then
    echo "Process killed successfully"
else
    echo "Process not found or already stopped"
fi
```

Make the script executable on Linux/macOS:
```bash
chmod +x kill_agent.sh
./kill_agent.sh
```

---

## Notes

- Use `-F` flag on Windows to force termination
- Use `-9` signal on Linux/macOS for immediate termination
- The process name may vary based on the binary you're running
- If the process is running as a service, you may need administrator/root privileges