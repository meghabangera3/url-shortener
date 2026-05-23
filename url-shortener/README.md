# URL Shortener

A high-performance URL shortening service built in Go, containerized with Docker, and orchestrated with Kubernetes.

## Tech Stack
- Go — REST API server using net/http
- Docker — Multi-stage build, ~10MB final image
- Kubernetes — 2-replica deployment
- GitHub Actions — CI pipeline on every push

## Features
- Shorten any URL to a 6-character code
- Thread-safe in-memory store using sync.Mutex
- /health endpoint for Kubernetes liveness probes

## Run Locally
go run .

## Run with Docker
docker build -t url-shortener:v1 .
docker run -p 8080:8080 url-shortener:v1

## Deploy to Kubernetes
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl port-forward service/url-shortener-service 8080:80

## API Endpoints
- GET  /          Home page
- POST /shorten   Shorten a URL
- GET  /r/{code}  Redirect to original
- GET  /health    Health check
