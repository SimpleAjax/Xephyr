# Xephyr Infrastructure Setup

This directory contains the Docker Compose infrastructure for the Xephyr project.

## Quick Start

### 1. Copy Environment Variables

```bash
cp .env.example .env
```

Edit `.env` with your preferred settings.

### 2. Start Infrastructure Services

```bash
# Start PostgreSQL and Redis
docker-compose up -d

# Or with explicit file
docker-compose -f docker-compose.yml up -d
```

### 3. Verify Services

```bash
# Check running containers
docker-compose ps

# Check PostgreSQL logs
docker-compose logs -f postgres

# Check Redis logs
docker-compose logs -f redis
```

### 4. Connect to Database

```bash
# Using psql from the container
docker-compose exec postgres psql -U xephyr -d xephyr

# Or from host (if psql is installed)
psql -h localhost -p 5432 -U xephyr -d xephyr
```

## Services

| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL | 5432 | Main database for application data |
| Redis | 6379 | Caching and session storage |

## Default Credentials

**PostgreSQL:**
- Host: `localhost`
- Port: `5432`
- User: `xephyr`
- Password: `xephyr123`
- Database: `xephyr`

## Database Migrations

The backend uses GORM AutoMigrate. When the server starts, it automatically:

1. Connects to PostgreSQL
2. Runs migrations for all models
3. Creates tables if they don't exist

To manually trigger migrations:

```bash
cd backend
go run cmd/server/main.go
```

## Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v
```

## Troubleshooting

### Port Already in Use

If you get "port already allocated" errors:

```bash
# Check what's using port 5432
lsof -i :5432  # macOS/Linux
netstat -ano | findstr 5432  # Windows

# Change ports in .env file
DB_PORT=5433
REDIS_PORT=6380
```

### Database Connection Issues

1. Ensure PostgreSQL is running:
   ```bash
   docker-compose ps
   ```

2. Check logs:
   ```bash
   docker-compose logs postgres
   ```

3. Verify credentials in `.env` match `docker-compose.yml`

### Reset Everything

```bash
# Remove all containers and volumes
docker-compose down -v

# Remove postgres data manually (if needed)
rm -rf postgres_data/

# Start fresh
docker-compose up -d
```

## Production Considerations

For production deployment:

1. **Change default passwords** in `.env`
2. **Enable SSL** for PostgreSQL connections
3. **Use external volume** for persistent storage
4. **Set up backups** for PostgreSQL data
5. **Use Redis with persistence** enabled
6. **Add monitoring** (Prometheus/Grafana)

## Architecture

```
┌─────────────────┐
│   Xephyr API    │
│   (Go/Gin)      │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐  ┌──▼────┐
│PostgreSQL│  │ Redis  │
│  (5432) │  │ (6379) │
└─────────┘  └────────┘
```
