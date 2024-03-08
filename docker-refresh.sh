echo "Starting docker-compose automation script"

echo "Building docker containers..."
docker-compose build

if [ $? -eq 0 ]; then
    echo "Build successful."
else
    echo "Build failed. Exiting..."
    exit 1
fi

echo "Taking down currently running containers..."
docker-compose down

echo "Bringing up containers in detached mode..."
docker-compose up 

echo "Docker-compose automation script finished!"