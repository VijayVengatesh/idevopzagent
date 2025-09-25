#!/bin/bash

echo "🧹 Starting Metrics Agent Uninstallation..."

# Step 1: Stop any running agent processes
echo "⛔️ Stopping running agent (if any)..."
pkill -f /usr/local/bin/metrics-agent
sleep 1

# Step 2: Remove binary
if [ -f "/usr/local/bin/metrics-agent" ]; then
    echo "🗑 Removing binary from /usr/local/bin..."
    sudo rm -f /usr/local/bin/metrics-agent
else
    echo "ℹ️ Binary not found in /usr/local/bin (already removed or moved)."
fi

# Step 3: Remove config directory
if [ -d "/etc/metrics-agent" ]; then
    echo "🗑 Removing config directory /etc/metrics-agent..."
    sudo rm -rf /etc/metrics-agent
else
    echo "ℹ️ Config directory /etc/metrics-agent does not exist."
fi

# Step 4: Remove log file
if [ -f "/var/log/metrics-agent.log" ]; then
    echo "🗑 Removing log file /var/log/metrics-agent.log..."
    sudo rm -f /var/log/metrics-agent.log
else
    echo "ℹ️ Log file /var/log/metrics-agent.log not found."
fi

# Step 5: Remove logrotate configuration
if [ -f "/etc/logrotate.d/metrics-agent" ]; then
    echo "🗑 Removing logrotate config /etc/logrotate.d/metrics-agent..."
    sudo rm -f /etc/logrotate.d/metrics-agent
else
    echo "ℹ️ Logrotate config not found."
fi

# Step 6: Final confirmation
echo "✅ Uninstallation complete."
