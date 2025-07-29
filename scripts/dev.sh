#!/bin/bash

echo "Starting development environment..."

# Start backend in background
echo "Starting backend server..."
cd backend
go run cmd/main.go &
BACKEND_PID=$!
cd ..

# Start frontend
echo "Starting frontend server..."
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

echo "Development servers started:"
echo "Backend: http://localhost:8080"
echo "Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all servers"

# Wait for interrupt
trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT
wait