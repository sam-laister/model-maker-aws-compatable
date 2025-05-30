#!/bin/bash
set -e

# Production entrypoint script
# Handles initialization, secrets, and graceful startup

# Colors for logging
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] INFO:${NC} $1"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARN:${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1"
}

# Function to handle graceful shutdown
graceful_shutdown() {
    log "Received shutdown signal, performing graceful shutdown..."
    if [ ! -z "$SERVER_PID" ]; then
        kill -TERM "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
    log "Graceful shutdown completed"
    exit 0
}

# Set up signal handlers
trap graceful_shutdown SIGTERM SIGINT

# Validate required environment variables
validate_env() {
    local required_vars=("PORT")
    local missing_vars=()
    
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            missing_vars+=("$var")
        fi
    done
    
    if [ ${#missing_vars[@]} -ne 0 ]; then
        error "Missing required environment variables: ${missing_vars[*]}"
        exit 1
    fi
}

# Initialize Google Cloud credentials if provided
init_gcp_credentials() {
    if [ ! -z "$GOOGLE_CREDENTIALS" ]; then
        log "Setting up Google Cloud credentials..."
        echo "$GOOGLE_CREDENTIALS" > /tmp/service-account-key.json
        export GOOGLE_APPLICATION_CREDENTIALS="/tmp/service-account-key.json"
        
        # Validate credentials format
        if ! command -v jq >/dev/null || ! echo "$GOOGLE_CREDENTIALS" | jq empty 2>/dev/null; then
            warn "Google credentials may not be valid JSON"
        fi
    elif [ ! -z "$GOOGLE_APPLICATION_CREDENTIALS" ] && [ -f "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
        log "Using existing Google Cloud credentials file: $GOOGLE_APPLICATION_CREDENTIALS"
    else
        warn "No Google Cloud credentials provided"
    fi
}

# Health check function
health_check() {
    local max_attempts=30
    local attempt=1
    
    log "Waiting for application to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -sf "http://localhost:${PORT}/health" >/dev/null 2>&1; then
            log "Application is ready and healthy"
            return 0
        fi
        
        log "Health check attempt $attempt/$max_attempts failed, retrying in 2s..."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    error "Application failed to become healthy after $max_attempts attempts"
    return 1
}

# Pre-flight checks
preflight_checks() {
    log "Running pre-flight checks..."
    
    # Check if required binaries exist
    local required_bins=("server")
    for bin in "${required_bins[@]}"; do
        if ! command -v "$bin" >/dev/null 2>&1; then
            error "Required binary not found: $bin"
            exit 1
        fi
    done
    
    # Check OpenMVG/OpenMVS binaries
    local openmvg_bins=("openMVG_main_SfMInit_ImageListing" "openMVG_main_ComputeFeatures")
    for bin in "${openmvg_bins[@]}"; do
        if ! command -v "$bin" >/dev/null 2>&1; then
            warn "OpenMVG binary not found: $bin"
        fi
    done
    
    log "Pre-flight checks completed successfully"
}

# Main execution
main() {
    log "Starting application with PID: $$"
    log "Running as user: $(whoami)"
    log "Working directory: $(pwd)"
    log "Environment: ${GIN_MODE:-development}"
    
    # Run initialization steps
    validate_env
    init_gcp_credentials
    preflight_checks
    
    # Start the application based on the command
    case "${1:-server}" in
        "server")
            log "Starting Go server on port $PORT..."
            exec /usr/local/bin/server &
            SERVER_PID=$!
            
            # Wait a moment for server to start
            sleep 5
            
            # Optional: Run health check in background
            if command -v curl >/dev/null 2>&1; then
                health_check &
            fi
            
            # Wait for the server process
            wait $SERVER_PID
            ;;
        "health")
            # Health check command
            curl -f "http://localhost:${PORT}/health"
            ;;
        *)
            # Execute any other command passed
            log "Executing command: $*"
            exec "$@"
            ;;
    esac
}

# Execute main function with all arguments
main "$@"