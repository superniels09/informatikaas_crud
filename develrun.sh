#!/bin/bash
# Build backend
echo "Building..."
./build.sh

# Start processes
setsid sh -c './informatikaas_crud'&
backend_pid=$!
setsid sh -c 'php -S localhost:9001'&
frontend_pid=$!

# Setup trap
trap 'kill -TERM -$backend_pid -$frontend_pid 2>/dev/null; exit' SIGINT SIGTERM

# Wait
wait $backend_pid $frontend_pid
