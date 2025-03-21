# Ecommerce Api

Each service has its own dedicated database to ensure scalability and minimal coupling. The structure of the project is shown below:

backend/
│── api-gateway/ # API Gateway (Handles HTTP requests and forwards to gRPC services)
│ ├── main.go
│ ├── routes.go
│ ├── handlers/
│ ├── middleware/
│ ├── Dockerfile
│── services/
│ ├── auth-service/
│ │ ├── proto/ # gRPC Proto definitions
│ │ ├── internal/
│ │ │ ├── handlers/ # gRPC Handlers
│ │ │ ├── services/ # Business logic
│ │ │ ├── repository/ # Database operations
│ │ │ ├── models/ # Data structures
│ │ ├── main.go
│ │ ├── Dockerfile
│ │ ├── config.yaml
│ ├── product-service/
│ ├── order-service/
│ ├── payment-service/
│ ├── notification-service/
│ ├── vendor-service/
│ ├── review-service/
│ ├── search-service/
│ ├── subscription-service/
│ ├── aliexpress-service/
│── shared/ # Common utilities used by all services
│ ├── grpc/ # gRPC client helper functions
│ ├── database/ # Database connection and migrations
│ ├── logger/ # Centralized logging
│ ├── config/ # Configuration loader
│ ├── middleware/ # Shared middleware (e.g., authentication, rate limiting)
│── infra/ # Infrastructure files (Docker, Kubernetes, etc.)
│ ├── docker-compose.yml # Local dev container setup
│ ├── k8s/ # Kubernetes deployment files
│── proto/ # All shared gRPC Protobuf files
│── Makefile # Automate build and run tasks
│── README.md
