#!/bin/sh
set -e

# Default config path
CONFIG_PATH=${CONFIG_PATH:-/app/config.yaml}

# Function to log messages
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1"
}

# Check if CONFIG_PATH is set
if [ -z "$CONFIG_PATH" ]; then
    log "ERROR: CONFIG_PATH environment variable is not set"
    exit 1
fi

# Check if the directory exists
CONFIG_DIR=$(dirname "$CONFIG_PATH")
if [ ! -d "$CONFIG_DIR" ]; then
    log "WARNING: Config directory $CONFIG_DIR does not exist. Creating..."
    mkdir -p "$CONFIG_DIR"
fi

# Wait for config file if it doesn't exist
log "Checking for config file at $CONFIG_PATH"
while [ ! -f "$CONFIG_PATH" ]; do
    log "Waiting for config file at $CONFIG_PATH..."
    sleep 1
done

# Check if the config file is readable
if [ ! -r "$CONFIG_PATH" ]; then
    log "ERROR: Config file at $CONFIG_PATH is not readable"
    exit 1
fi

log "Starting AegisGate with config at $CONFIG_PATH"

# Start the application
exec /app/aegisgate "$CONFIG_PATH" 