#!/bin/bash

echo "Stopping containers"
docker-compose down

echo "Removing Docker images"
docker rmi frontend-image
docker rmi backend-image

echo "Docker images removed"
