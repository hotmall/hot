package commands

const goStart = `#!/bin/sh

# Configuration
EXEC=%s
PORT=53720
APP_HOME=$(cd "$(dirname "$0")/.." && pwd) || exit 1
CMD="$APP_HOME/bin/$EXEC"
NOHUP_LOG="$APP_HOME/var/log/nohup.out"
TIMEOUT=1  # seconds to wait for process start check

# Validate environment
if [ -z "$APP_HOME" ]; then
    echo "ERROR: Could not determine APP_HOME directory" >&2
    exit 1
fi

if [ ! -d "$APP_HOME/bin" ]; then
    echo "ERROR: Bin directory not found: $APP_HOME/bin" >&2
    exit 1
fi

if [ ! -x "$CMD" ]; then
    echo "ERROR: Executable not found or not executable: $CMD" >&2
    exit 1
fi

# Ensure log directory exists
LOG_DIR=$(dirname "$NOHUP_LOG")
if ! mkdir -p "$LOG_DIR"; then
    echo "ERROR: Failed to create log directory: $LOG_DIR" >&2
    exit 1
fi

# Start the service
echo "Starting $EXEC..."
nohup "$CMD" -p "$APP_HOME" -port "$PORT" > "$NOHUP_LOG" 2>&1 &
PID=$!

# Get the PID and verify process is running
echo "Service $EXEC started with PID $PID"

# Wait a moment then check if process is still running
sleep "$TIMEOUT"
if ! kill -0 "$PID" 2>/dev/null; then
    echo "ERROR: Process $PID did not stay running. Check logs: $NOHUP_LOG" >&2
    tail -n 20 "$NOHUP_LOG" >&2
    exit 1
fi

# Success message
echo "Log output can be found at: $NOHUP_LOG"
echo "Startup completed successfully"
exit 0
`

const goStop = `#!/bin/sh

# Configuration
EXEC=%s
APP_HOME=$(cd "$(dirname "$0")/.." && pwd) || exit 1
CMD="$APP_HOME/bin/$EXEC"
TIMEOUT=10  # seconds to wait for graceful shutdown
KILL_TIMEOUT=3  # seconds to wait after SIGTERM before SIGKILL

# Validate environment
if [ -z "$APP_HOME" ]; then
    echo "ERROR: Could not determine APP_HOME directory" >&2
    exit 1
fi

if [ ! -x "$CMD" ]; then
    echo "ERROR: Executable not found or not executable: $CMD" >&2
    exit 1
fi

# Find running process
PID=$(pgrep -f "$CMD -p $APP_HOME")
if [ -z "$PID" ]; then
    echo "WARNING: No running instances of $EXEC found"
    exit 0
fi

echo "Stopping $EXEC (PID $PID)..."

# Send graceful shutdown signal
if ! "$CMD" -p "$APP_HOME" -s stop; then
    echo "ERROR: Failed to send stop command to $EXEC" >&2
    exit 1
fi

# Wait for process to terminate gracefully
ELAPSED=0
while kill -0 "$PID" 2>/dev/null && [ "$ELAPSED" -lt "$TIMEOUT" ]; do
    sleep 1
    ELAPSED=$((ELAPSED + 1))
done

# Force kill if still running
if kill -0 "$PID" 2>/dev/null; then
    echo "WARNING: Process did not stop gracefully after $TIMEOUT seconds, forcing termination..."
    kill -TERM "$PID"
    
    sleep "$KILL_TIMEOUT"
    if kill -0 "$PID" 2>/dev/null; then
        kill -KILL "$PID"
        echo "Sent SIGKILL to $PID"
    fi
fi

# Verify process is stopped
if kill -0 "$PID" 2>/dev/null; then
    echo "ERROR: Failed to stop process $PID" >&2
    exit 1
else
    echo "Successfully stopped $EXEC (PID $PID)"
    exit 0
fi
`
