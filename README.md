# Restaurant Management API

A comprehensive REST API for managing restaurant operations including table sessions, menu items, orders, and more. Built with Go and Gin framework, containerized with Docker for easy deployment.

## Features

- **Session Management**: Handle customer table sessions with start/end times
- **Table Management**: Manage restaurant tables and availability
- **Menu Management**: Create and manage menu items with categories
- **Order Management**: Process customer orders with menu items
- **Category Management**: Organize menu items by categories
- **Pagination & Filtering**: Efficient data retrieval with pagination support
- **Input Validation**: Robust validation using struct tags and custom validators
- **Error Handling**: Comprehensive error handling with custom middleware
- **API Documentation**: Auto-generated Swagger documentation
- **Database Migrations**: PostgreSQL migrations for schema management

## Tech Stack

- **Backend**: Go 1.25.5
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Validation**: go-playground/validator
- **Containerization**: Docker & Docker Compose
- **API Documentation**: Swagger/OpenAPI

## Prerequisites

- Docker Desktop (latest version)
- Docker Compose (included with Docker Desktop)
- Git

## Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd restaurant
```

### 2. Start the Application
```bash
# Build and start all containers (PostgreSQL + API)
docker compose up -d

# View logs
docker compose logs -f

# View application logs only
docker compose logs -f app
```

### 3. Verify Setup
```bash
# Check running containers
docker compose ps

# Test API health
curl http://localhost:8080/docs
```

### 4. Stop the Application
```bash
docker compose down
```

## API Documentation

Once the application is running, visit:
- **Swagger UI**: http://localhost:8080/docs
- **API Base URL**: http://localhost:8080

## Testing with Postman

1. Import the Postman collection: `Restaurant_API.postman_collection.json`
2. Import the environment: `Restaurant_API.postman_environment.json`
3. Update environment variables if needed (default: localhost:8080)
4. Run the collection to test all endpoints

## Project Structure

```
restaurant/
├── cmd/app/              # Application entry point
├── internal/
│   ├── cache/           # Caching layer
│   ├── middleware/      # HTTP middleware
│   ├── optimization/    # Query optimization
│   ├── pool/           # Connection pooling
│   ├── response/       # Response formatting
│   ├── session/        # Session management module
│   │   ├── handler/    # HTTP handlers
│   │   ├── models/     # Data models
│   │   ├── repository/ # Data access layer
│   │   ├── service/    # Business logic
│   │   └── validation/ # Input validation
│   └── shutdown/       # Graceful shutdown
├── migrations/         # Database migrations
├── docs/               # Generated API docs
├── documentation/      # Project documentation
├── docker-compose.yml  # Docker services
├── Dockerfile         # Application container
└── go.mod             # Go dependencies
```

## Development

### Local Development Setup

1. **Install Go 1.25.5+**
2. **Set up PostgreSQL** (or use Docker)
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Run migrations**:
   ```bash
   # Using migrate tool or custom script
   ```
5. **Start the application**:
   ```bash
   go run ./cmd/app
   ```

### Database Migrations

Migrations are located in the `migrations/` directory. Use your preferred migration tool to apply them.

### Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...
```

## API Endpoints

### Sessions
- `GET /sessions` - List all sessions
- `POST /sessions` - Create new session
- `GET /sessions/{id}` - Get session by ID
- `PUT /sessions/{id}` - Update session
- `DELETE /sessions/{id}` - Delete session

### Tables
- `GET /tables` - List all tables
- `POST /tables` - Create new table
- `GET /tables/{id}` - Get table by ID
- `DELETE /tables/{id}` - Delete table

### Menu Items
- `GET /menu` - List menu items (with pagination)
- `POST /menu` - Create menu item
- `GET /menu/{id}` - Get menu item by ID
- `PUT /menu/{id}` - Update menu item
- `DELETE /menu/{id}` - Delete menu item

### Categories
- `GET /categories` - List all categories
- `POST /categories` - Create new category
- `GET /categories/{id}` - Get category by ID
- `PUT /categories/{id}` - Update category
- `DELETE /categories/{id}` - Delete category

### Orders
- `GET /orders` - List all orders
- `POST /orders` - Create new order
- `GET /orders/{id}` - Get order by ID
- `PUT /orders/{id}` - Update order
- `DELETE /orders/{id}` - Delete order

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This backend code for the Intellidine Restaurant App is licensed under the
Creative Commons Attribution-NonCommercial 4.0 International License.

You may not use this work for commercial purposes.  
See the LICENSE file for full details: <https://creativecommons.org/licenses/by-nc/4.0/>.
