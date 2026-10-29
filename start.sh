    echo "Installing dependencies..."
    npm install --prefix frontend/
    echo "Dependencies installed."
    
# Start both frontend and backend servers concurrently
(
    echo "Starting backend server..."
    cd ./backend
    go run .
) & (
    echo "Starting frontend server..."
    cd ./frontend
    npm run build
    go run serveFrontend.go 
    
)
# Wait for both servers to finish
wait