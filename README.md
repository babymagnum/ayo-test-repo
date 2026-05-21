## Prerequisites

| Tool | Version | Notes |
|---|---|---|
| Go | 1.24+ | Match go.mod |
| PostgreSQL | 14+ | Database backend |
| Air (optional) | latest | Hot reload (see below) |

---

## Setup & Run

### 1. Clone & enter project

```bash
git clone <repo-url>
cd ayo-test
```

### 2. Setup PostgreSQL

Make sure PostgreSQL is running locally.

```bash
# Check version
psql --version

# Create database
createdb ayo-test
```

If PostgreSQL not installed:

```bash
brew install postgresql@16
brew services start postgresql@16
```

### 3. Configure environment

```bash
cp env.example .env
```

Edit `.env` and fill:

```env
APP_ENV=development
DB_ADDR=host=localhost user=postgres password=postgres dbname=ayo-test port=5432 sslmode=disable
SECRET_KEY=your-secret-key
HTTP_PORT=8081
SHUTDOWN_TTL=30
SWAGGER_HOST=localhost:8081
SWAGGER_PATH=/v1
```

> Adjust `user` and `password` to match your local PostgreSQL credentials.

### 4. Download dependencies

```bash
go mod tidy
```

### 5. Run

```bash
go run ./cmd/api/
```

Migrations run automatically on startup. You should see:

```
Database connection successfully established with GORM.
Migrations applied successfully
starting http server
```

Server starts at `http://localhost:8081`.

---

## Optional Tools

### Hot reload with Air

[Air](https://github.com/air-verse/air) watches file changes and restarts the server automatically.

Install:

```bash
go install github.com/air-verse/air@latest
```

Make sure `.air.toml` matches your OS:

- **Local (Windows/Mac):** `bin = "./bin/api.exe"` + `cmd = "go build -o ./bin/ ./cmd/api/"`
- **Linux / Docker:** `bin = "./bin/api"` + `cmd = "go build -o ./bin/api ./cmd/api/"`

Run:

```bash
air -c .air.toml
```

### Swagger docs

Install swag CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Update docs:

1. Go to `cmd/api`
2. In `base-entity.go` — comment `DeletedAt`, uncomment the implementation below it
3. Run:

   ```bash
   swag init
   ```

4. Revert the `DeletedAt` changes after generation

---

## Database Migrations

Migrations run **automatically** on startup via `internal/db/migrate.go`. No manual step needed.

For manual migration management, install the [migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate):

```bash
# Create new migration
migrate create -seq -ext sql -dir ./migrations create_users

# Run migrations manually
export $(grep -v '^#' .env | xargs) && migrate -path ./migrations -database "${DB_ADDR}" up
```

---

## Project Structure

```
cmd/api/           # Entry point, handlers (controllers), middleware, DTOs
internal/
  config/          # (todo) App config struct
  db/              # Database connection + auto-migration
  interfaces/      # Interface definitions (service contracts)
  logger/          # Zap logger setup
  service/         # Business logic layer
  store/           # Data access / repository layer
  utils/           # Shared utilities
migrations/        # SQL migration files
```
