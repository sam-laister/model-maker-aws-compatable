#!/bin/sh
set -e

# Production entrypoint script
# Handles initialization, secrets, and graceful startup

# Colors for logging
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    printf "%s[%s] INFO:%s %s\n" "$GREEN" "$(date +'%Y-%m-%d %H:%M:%S')" "$NC" "$1"
}

warn() {
    printf "%s[%s] WARN:%s %s\n" "$YELLOW" "$(date +'%Y-%m-%d %H:%M:%S')" "$NC" "$1"
}

error() {
    printf "%s[%s] ERROR:%s %s\n" "$RED" "$(date +'%Y-%m-%d %H:%M:%S')" "$NC" "$1"
}

# Initialize Google Cloud credentials if provided
init_gcp_credentials() {
    if [ -n "$GOOGLE_CREDENTIALS" ]; then
        log "Setting up Google Cloud credentials..."
        printf "%s" "$GOOGLE_CREDENTIALS" > /tmp/service-account-key.json
        export GOOGLE_APPLICATION_CREDENTIALS="/tmp/service-account-key.json"

        if command -v jq >/dev/null 2>&1; then
            echo "$GOOGLE_CREDENTIALS" | jq empty 2>/dev/null || warn "Google credentials may not be valid JSON"
        else
            warn "jq not installed, cannot validate JSON"
        fi
    elif [ -n "$GOOGLE_APPLICATION_CREDENTIALS" ] && [ -f "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
        log "Using existing Google Cloud credentials file: $GOOGLE_APPLICATION_CREDENTIALS"
    else
        warn "No Google Cloud credentials provided"
    fi
}

# Main execution
main() {
    log "Starting application with PID: $$"
    log "Running as user: $(whoami)"
    log "Working directory: $(pwd)"
    log "Environment: ${GIN_MODE:-development}"

    # Run initialization steps
    init_gcp_credentials

    # Execute the command passed as arguments
    if [ $# -gt 0 ]; then
        log "Executing command: $*"
        exec "$@"
    else
        error "No command provided to execute"
        exit 1
    fi
}

main "$@"
