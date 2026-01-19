#!/bin/bash

# Check if services are running
services=("nginx" "postgres" "backend" "frontend" "xray")

for service in "${services[@]}"; do
    if docker ps | grep -q $service; then
        echo "✅ $service is running"
    else
        echo "❌ $service is NOT running"
    fi
done

# Check HTTP status
if curl -s -o /dev/null -w "%{http_code}" https://localhost > /dev/null 2>&1; then
    echo "✅ Web server is responding"
else
    echo "❌ Web server is not responding"
fi
