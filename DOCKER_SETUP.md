# Docker Setup Guide

Complete guide to containerize and run the Restaurant Management API with PostgreSQL.

## Prerequisites

- Docker Desktop (latest version)
- Docker Compose (included with Docker Desktop)
- Git (to clone the repository)

## Quick Start

### 1. Clone & Navigate
```bash
cd ~/Restaurant
```

### 2. Build and Start Containers
```bash
# Build and start all containers
docker compose up -d

# View logs
docker compose logs -f

# View app logs only
docker compose logs -f app
```

### 3. Verify Setup
```bash
# Check if containers are running
docker compose ps

# Test API endpoint
curl http://localhost:8080/docs

# Expected response:
# {
#   "title": "Restaurant Management API",
#   "version": "1.0",
#   ...
# }
```

### 4. Stop Containers
```bash
# Stop all containers
docker compose down

# Stop and remove volumes (cleanup database)
docker compose down -v
```

---

## Directory Structure

```
Restaurant/
├── Dockerfile                 # Multi-stage build for Go app
├── docker-compose.yml        # Orchestration of services
├── .env                      # Environment variables
├── .dockerignore             # Files to exclude from Docker build
├── cmd/
│   └── app/
│       └── main.go           # Updated with env var support
├── migrations/
│   ├── 000001_initial.up.sql
│   └── 000001_initial.down.sql
├── internal/
└── ...
```

---

## Services Architecture

### 1. PostgreSQL Database
- **Image**: `postgres:15-alpine`
- **Container Name**: `restaurant-db`
- **Port**: `5432` (mapped to host)
- **Health Check**: Every 10s
- **Volume**: `postgres_data` (persistent storage)
- **Initialization**: Auto-runs migrations from `./migrations`

### 2. Go Application
- **Build**: Multi-stage Dockerfile
- **Container Name**: `restaurant-app`
- **Port**: `8080` (mapped to host)
- **Health Check**: `/docs` endpoint every 10s
- **Depends On**: `postgres` (healthy)
- **Network**: `restaurant-network` (internal communication)

### 3. pgAdmin (Optional - Debug Profile)
- **Image**: `dpage/pgadmin4:latest`
- **Container Name**: `restaurant-pgadmin`
- **Port**: `5050` (mapped to host)
- **Access**: http://localhost:5050
- **Profile**: `debug` (only start with `--profile debug`)

---

## Environment Variables

All environment variables are in `.env` file:

```bash
# Database Configuration
DB_HOST=postgres              # Docker service name
DB_PORT=5432
DB_USER=restaurant
DB_PASSWORD=restaurant_password
DB_NAME=restaurant
DB_SSLMODE=disable

# Application Configuration
APP_ENV=development
APP_PORT=8080
LOG_LEVEL=info

# pgAdmin (optional)
PGADMIN_EMAIL=admin@restaurant.local
PGADMIN_PASSWORD=admin
PGADMIN_PORT=5050
```

### Key Points:
- `DB_HOST=postgres` - Uses Docker service name (not localhost)
- Modify `.env` to change configuration
- Restart containers after changing `.env`

---

## Docker Commands Cheat Sheet

### Basic Operations
```bash
# Start containers (build if needed)
docker compose up -d

# Stop containers
docker compose down

# Stop and remove volumes
docker compose down -v

# Rebuild images
docker compose build --no-cache

# View logs
docker compose logs -f
docker compose logs -f app
docker compose logs -f postgres
```

### Container Management
```bash
# Execute command in container
docker compose exec app sh
docker compose exec postgres psql -U restaurant -d restaurant

# View running containers
docker compose ps

# View container details
docker inspect restaurant-app
docker inspect restaurant-db
```

### Database Access
```bash
# Connect to PostgreSQL from host (if psql installed)
psql -h localhost -U restaurant -d restaurant

# From Docker container
docker compose exec postgres psql -U restaurant -d restaurant

# Run SQL query
docker compose exec postgres psql -U restaurant -d restaurant \
  -c "SELECT * FROM sessions;"
```

---

## Troubleshooting

### Container Won't Start

**Problem**: App container exits immediately
```bash
# Check logs
docker compose logs app

# If database connection error:
# Make sure postgres is healthy first
docker compose logs postgres
```

**Solution**: 
- Ensure postgres is running and healthy
- Check `.env` file for correct credentials
- Verify port 5432 is not already in use

### Database Connection Issues

**Problem**: `failed to connect to database`

**Solution**:
```bash
# 1. Check if postgres is healthy
docker compose ps postgres

# 2. Verify environment variables
docker compose exec app env | grep DB_

# 3. Test connection
docker compose exec postgres psql -U restaurant -d restaurant -c "\dt"
```

### Port Already in Use

**Problem**: `error starting service postgres`

**Solution**:
```bash
# Find process using port 5432
lsof -i :5432

# Or change port in .env
DB_PORT=5433  # Change to different port
```

### Migrations Not Running

**Problem**: Tables not created in database

**Solution**:
```bash
# Check if migrations volume is mounted
docker compose exec postgres ls /docker-entrypoint-initdb.d/

# Manually run migrations
docker compose exec postgres psql -U restaurant -d restaurant \
  -f /docker-entrypoint-initdb.d/000001_initial.up.sql

# Or rebuild from scratch
docker compose down -v
docker compose up -d
```

---

## Advanced Usage

### Development with Live Reload

Using `air` for live reload:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with live reload
air

# Note: This runs locally, not in Docker
```

### Production Deployment

For production, modify `docker-compose.yml`:

```yaml
services:
  app:
    environment:
      APP_ENV: production
      LOG_LEVEL: error
    restart: always  # Always restart on failure
```

### Running with Debug Profile (pgAdmin)

```bash
# Start with debug profile to include pgAdmin
docker compose --profile debug up -d

# Access pgAdmin at http://localhost:5050
# Login with PGADMIN_EMAIL and PGADMIN_PASSWORD from .env
```

### Custom Configuration

Create `.env.production` for production settings:

```bash
DB_PASSWORD=secure_password_here
APP_ENV=production
LOG_LEVEL=error
PGADMIN_PASSWORD=secure_admin_password
```

Run with custom env file:
```bash
docker compose --env-file .env.production up -d
```

---

## Database Initialization

### Automatic Initialization
- Migrations in `./migrations/` auto-run on first start
- Tables created: `tables`, `sessions`, `orders`, `order_items`, `menu_items`, `categories`
- Indexes automatically created for performance

### Manual Migration

```bash
# Connect to database
docker compose exec postgres psql -U restaurant -d restaurant

# Run migration
\i /docker-entrypoint-initdb.d/000001_initial.up.sql

# Verify tables
\dt
```

### Backup Database

```bash
# Backup to file
docker compose exec postgres pg_dump -U restaurant restaurant > backup.sql

# Restore from file
cat backup.sql | docker compose exec -T postgres psql -U restaurant -d restaurant
```

---

## Network Communication

### Docker Network: `restaurant-network`

All services communicate via bridge network:

```
┌─────────────────────────────────────┐
│   restaurant-network (bridge)       │
├─────────────────────────────────────┤
│                                     │
│  ┌──────────┐      ┌────────────┐  │
│  │  app     │◄────►│ postgres   │  │
│  │:8080     │      │ :5432      │  │
│  └──────────┘      └────────────┘  │
│                                     │
│  ┌──────────┐                       │
│  │ pgadmin  │◄─────► postgres       │
│  │ :5050    │       (if --profile)  │
│  └──────────┘                       │
│                                     │
└─────────────────────────────────────┘
       ↓ (port mapping)
   Host Machine
   localhost:8080 → app:8080
   localhost:5432 → postgres:5432
   localhost:5050 → pgadmin:80
```

### Service Discovery

- App connects to `postgres:5432` (not localhost:5432)
- Docker DNS resolves `postgres` to container IP
- Automatic service discovery

---

## Performance Tuning

### Database Connection Pool

Automatically configured in main.go:
- MaxOpenConns: 25
- MaxIdleConns: 5
- ConnMaxLifetime: 5 minutes

Adjust for high traffic:
```bash
# In .env (future enhancement)
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=20
```

### Container Resource Limits

Modify `docker-compose.yml`:

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

---

## Health Checks

### App Health Check
```bash
# Check endpoint
curl http://localhost:8080/docs

# Or from Docker
docker compose exec app wget -q -O- http://localhost:8080/docs
```

### Database Health Check
```bash
# Check connectivity
docker compose exec postgres pg_isready -U restaurant
```

---

## Next Steps

1. **Monitor**: Use `docker compose logs -f` to monitor
2. **Scale**: Add multiple app instances with load balancer
3. **Backup**: Set up automated database backups
4. **CI/CD**: Integrate with GitHub Actions or GitLab CI
5. **Security**: Use secrets management in production

---

## File Reference

- **Dockerfile**: Multi-stage build (production-ready)
- **docker-compose.yml**: Service orchestration
- **.env**: Configuration defaults
- **.dockerignore**: Build context optimization
- **cmd/app/main.go**: Updated with env var support

