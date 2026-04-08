#!/bin/sh
set -e

echo "Starting HydraDNS dataplane..."
/app/dataplane &
DATAPLANE_PID=$!

# Wait for gRPC port to be ready using /proc/net/tcp6
echo "Waiting for dataplane gRPC on port 50051..."
GRPC_PORT_HEX=$(printf '%X' 50051)
for i in $(seq 1 30); do
    if ! kill -0 $DATAPLANE_PID 2>/dev/null; then
        echo "Dataplane process died during startup"
        exit 1
    fi
    if grep -q ":${GRPC_PORT_HEX} " /proc/net/tcp6 2>/dev/null || grep -q ":${GRPC_PORT_HEX} " /proc/net/tcp 2>/dev/null; then
        echo "Dataplane gRPC ready"
        break
    fi
    if [ "$i" -eq 30 ]; then
        echo "Timeout waiting for dataplane gRPC"
        exit 1
    fi
    sleep 1
done

echo "Starting HydraDNS controlplane..."
/app/controlplane &
CONTROLPLANE_PID=$!

# Trap signals for graceful shutdown
trap 'echo "Shutting down..."; kill $DATAPLANE_PID $CONTROLPLANE_PID 2>/dev/null; wait' TERM INT

# Wait for either process to exit
wait -n $DATAPLANE_PID $CONTROLPLANE_PID 2>/dev/null || true
EXIT_CODE=$?

echo "A process exited (code=$EXIT_CODE), shutting down..."
kill $DATAPLANE_PID $CONTROLPLANE_PID 2>/dev/null || true
wait
exit $EXIT_CODE
