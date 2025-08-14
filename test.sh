#!/bin/bash

echo "Running Gin Simple App Tests..."
echo "================================"

# Run tests with verbose output
go test ./tests/... -v

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All tests passed!"
else
    echo ""
    echo "❌ Some tests failed!"
    exit 1
fi
