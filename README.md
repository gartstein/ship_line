# Ship Line (Interview Project)

Ship Line is a mono-repository containing a Go backend (Go 1.23) and a React frontend. It demonstrates how to calculate the distribution of packs to fulfill orders with minimal leftover. This project is specifically designed for interview purposes.

## Why Two Implementations of `CalculatePacks`?

1. **Dynamic Programming (DP)**
    - Guarantees an optimal solution (minimizes leftover).
    - More suitable for smaller orders due to higher accuracy but can be slower for large inputs.

2. **Greedy Approach**
    - Faster for large orders.
    - May produce slightly suboptimal results, but it scales well and avoids performance bottlenecks.

By having both implementations, the project handles orders of varying sizes efficiently and accurately.

## Repository Structure

- **backend/**
    - **db/**: BoltDB logic
    - **handlers/**: Gin route handlers
    - **services/**: Business logic (including both DP and Greedy calculations)
    - **swagger/**: API definitions
    - **utils/**: Utility functions
    - `main.go`: Application entry point
    - `go.mod`: Go modules file

- **deployment/**: Docker Compose and Kubernetes configuration
- **docs/**: Documentation
- **frontend/**: React application
- **Makefile**: Contains commands for building, testing, running, and deploying
- **README.md**: Project documentation

## Makefile Targets

- **`build-backend`**: Builds the Go backend
- **`build-frontend`**: Builds the React frontend
- **`run`**: Runs both the backend and frontend
- **`test-backend`**: Runs Go tests for the backend
- **`build-docker`**: Builds Docker images for the backend and frontend
- **`run-docker`**: Starts Docker containers via Docker Compose
- **`stop-docker`**: Stops Docker Compose services


## Getting Started

1. **Clone the repository**:
   ```bash
   git clone git@github.com:gartstein/ship_line.git
   cd ship_line

## Accessing the Website Locally

After you have built and started the application (e.g., using `make run-docker`), you can access the two main parts of the project locally:

- **Frontend:**  
  The React frontend is served through an Nginx container. Open your browser and go to:  
  [http://localhost;3000](http://localhost:3000)

- **Backend:**  
  The Go backend runs on port 8080. You can access its API endpoints directly. For example, to fetch pack sizes, visit:  
  [http://localhost:8080/v1/pack-sizes](http://localhost:8080/v1/pack-sizes)

If you prefer to run the frontend or backend separately (without Docker), use the corresponding Makefile targets to build and run them on their respective ports.