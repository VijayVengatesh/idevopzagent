#!/bin/bash

# Get device key from arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --device-key) USER_ID="$2"; shift ;;
        --device-key=*) USER_ID="${1#*=}" ;;
        *) echo "Unknown parameter passed: $1. Use --device-key"; exit 1 ;;
    esac
    shift
done

if [ -z "$USER_ID" ]; then
    read -rp "Enter device key: " USER_ID
fi

if [ -z "$USER_ID" ]; then
    echo "Error: Device key is required. Use --device-key <device_key>"
    exit 1
fi

OS=$(uname -s)
ARCH=$(uname -m)
echo "Step 1: Detected OS: $OS"
echo "Step 2: Detected architecture: $ARCH"

case "$OS" in
    Linux)
        case "$ARCH" in
            x86_64) BIN_NAME="idevopzagent-linux-amd64" ;;
            i386|i686) BIN_NAME="idevopzagent-linux-386" ;;
            aarch64|arm64) BIN_NAME="idevopzagent-linux-arm64" ;;
            armv7l|armv6l) BIN_NAME="idevopzagent-linux-arm" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    Darwin)
        case "$ARCH" in
            x86_64) BIN_NAME="idevopzagent-darwin-amd64" ;;
            arm64|aarch64) BIN_NAME="idevopzagent-darwin-arm64" ;;
            *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

AGENT_URL="https://github.com/VijayVengatesh/goscript/releases/download/v1.0.0/${BIN_NAME}"

echo "Step 3: Downloading agent from $AGENT_URL"
curl -L --fail -o idevopzagent "$AGENT_URL"
if [ $? -ne 0 ]; then
    echo "Error: Failed to download the agent binary. Please check the URL or internet connection."
    exit 1
fi

echo "Step 4: Installing agent binary..."
chmod +x idevopzagent || { echo "Failed to make binary executable"; exit 1; }
sudo mv idevopzagent /usr/local/bin/idevopzagent || { echo "Failed to move binary to /usr/local/bin"; exit 1; }
echo "Binary installed at /usr/local/bin/idevopzagent"

echo "Step 5: Starting agent..."
nohup /usr/local/bin/idevopzagent --device-key "$USER_ID" >/dev/null 2>&1 &
sleep 2

if [ "$OS" = "Darwin" ]; then
    PRIMARY_CONFIG_DIR="/Library/Application Support/idevopzagent"
    FALLBACK_CONFIG_DIR="$HOME/Library/Application Support/idevopzagent"
else
    PRIMARY_CONFIG_DIR="/var/lib/idevopzagent"
    FALLBACK_CONFIG_DIR="$HOME/.local/share/idevopzagent"
fi

CONFIG_DIR="$PRIMARY_CONFIG_DIR"
if ! mkdir -p "$CONFIG_DIR" >/dev/null 2>&1; then
    CONFIG_DIR="$FALLBACK_CONFIG_DIR"
    mkdir -p "$CONFIG_DIR" >/dev/null 2>&1 || true
fi
LOG_FILE="$CONFIG_DIR/logs.txt"

if [ "$OS" = "Linux" ]; then
    if command -v systemctl >/dev/null 2>&1; then
        echo "Step 6: Enabling autostart (systemd)..."
        SERVICE_PATH="/etc/systemd/system/idevopzagent.service"
        sudo tee "$SERVICE_PATH" >/dev/null <<EOF
[Unit]
Description=IDevopzAgent
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/idevopzagent --device-key $USER_ID
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
        sudo systemctl daemon-reload
        sudo systemctl enable --now idevopzagent
    else
        echo "Warning: systemctl not found. Autostart not configured."
    fi
else
    echo "Step 6: Enabling autostart (launchd)..."
    PLIST_PATH="/Library/LaunchDaemons/com.idevopzagent.agent.plist"
    sudo tee "$PLIST_PATH" >/dev/null <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.idevopzagent.agent</string>
  <key>ProgramArguments</key>
  <array>
    <string>/usr/local/bin/idevopzagent</string>
    <string>--device-key</string>
    <string>$USER_ID</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
</dict>
</plist>
EOF
    sudo launchctl load -w "$PLIST_PATH"
fi

pgrep -f idevopzagent > /dev/null
if [ $? -eq 0 ]; then
    echo "Installation complete. Agent is running in background."
    echo "Config: $CONFIG_DIR/config.json (encrypted device key + machine ID)"
    echo "Logs: $LOG_FILE"
else
    echo "Agent failed to start. Check $LOG_FILE"
    exit 1
fi
