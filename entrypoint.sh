#!/bin/sh
set -e

echo "Running migration..."
./migration

echo "Starting server..."
./server
