#!/bin/bash

# Demo Runner Script
set -e

echo "=== glog Demo Runner ==="
echo ""
echo "Available demos:"
echo "  1. basic       - Basic functionality demo"
echo "  2. webapp      - Web application simulation"
echo "  3. performance - Performance and stress tests"
echo ""

if [ $# -eq 0 ]; then
    echo "Usage: ./run_demos.sh [basic|webapp|performance|all]"
    exit 0
fi

case $1 in
    basic|1)
        echo "Running Basic Demo..."
        go run demo/basic/main.go
        ;;
    webapp|web|2)
        echo "Running Web App Demo..."
        go run demo/webapp/main.go
        ;;
    performance|perf|3)
        echo "Running Performance Demo..."
        go run demo/performance/main.go
        ;;
    all)
        echo "Running all demos..."
        echo ""
        echo "=== Basic Demo ==="
        go run demo/basic/main.go
        echo ""
        echo "=== Web App Demo ==="
        go run demo/webapp/main.go
        echo ""
        echo "=== Performance Demo ==="
        go run demo/performance/main.go
        ;;
    *)
        echo "Unknown demo: $1"
        exit 1
        ;;
esac

echo ""
echo "Demo completed!"
