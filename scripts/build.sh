#!/bin/bash

echo "Building License Manager..."

# Build backend
echo "Building backend..."
cd backend
go mod tidy
go build -o ../bin/license-manager cmd/main.go
cd ..

# Build frontend
echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

echo "Build completed successfully!"
echo "Backend binary: ./bin/license-manager"
echo "Frontend dist: ./frontend/dist"