# Project Variables
PROJECT_NAME=ship_line
GO_SERVICE=backend
FRONTEND_SERVICE=frontend
REGION=europe-central2
CLUSTER_NAME=ship-line-cluster
DOCKER_REPO=gcr.io/$(GCP_PROJECT_ID)

# Default target
.DEFAULT_GOAL := help

## 🏗️ Build Go Backend
build-backend:
	@echo "🚀 Building Go backend..."
	cd backend && go build -o main .

## 🎨 Build React Frontend
build-frontend:
	@echo "📦 Building React frontend..."
	cd frontend && npm install && npm run build

## 🐳 Build Docker Images
build-docker:
	@echo "🐳 Building Docker images..."
	docker build -t $(DOCKER_REPO)/$(GO_SERVICE) -f backend/Dockerfile.backend ./backend
	docker build -t $(DOCKER_REPO)/$(FRONTEND_SERVICE) -f frontend/Dockerfile.frontend ./frontend

## 📝 Generate Swagger Models
openapi-models:
	@echo "Generating OpenAPI models from docs/swagger.yaml..."
	oapi-codegen -generate types -package swagger -o backend/swagger/models.gen.go docs/swagger.yaml

## 🚀 Run Backend
run-backend:
	@echo "🔧 Running Go backend..."
	cd backend && go run main.go

## 🌐 Run Frontend
run-frontend:
	@echo "🌍 Running React frontend..."
	cd frontend && npm start

## 🐳 Run Docker Compose
run-docker:
	@echo "🐳 Running services with Docker Compose..."
	cd deployment && docker-compose up --build

## 🛠️ Stop Docker Compose
stop-docker:
	@echo "🛑 Stopping Docker Compose..."
	cd deployment && docker-compose down

## 🧪 Run Backend Tests
test-backend:
	@echo "🧪 Running tests for Go backend..."
	cd backend && go test ./...

## 📦 Push Docker Images to GCP
push-docker:
	@echo "🚀 Pushing Docker images to Google Container Registry..."
	docker push $(DOCKER_REPO)/$(GO_SERVICE)
	docker push $(DOCKER_REPO)/$(FRONTEND_SERVICE)

## ☁️ Deploy to Google Kubernetes Engine (GKE)
deploy-gke:
	@echo "🚀 Deploying to GKE..."
	gcloud container clusters get-credentials $(CLUSTER_NAME) --region $(REGION)
	kubectl apply -f deployment/k8s/

## 🔍 Get Kubernetes Service Info
get-service:
	@echo "🔍 Getting external IP of services..."
	kubectl get services

## 📌 Help Menu
help:
	@echo "Available commands:"
	@awk '/^## /{help=$$0; sub(/^## /,"",help); next} /^[a-zA-Z0-9_-]+:/ && help { \
	  split($$1, target, ":"); \
	  printf "\033[36m%-20s\033[0m %s\n", target[1], help; \
	  help=""; \
	}' $(MAKEFILE_LIST)