# BracketForge API

A backend-first tournament management engine supporting Single Elimination, Round Robin, and Group Stage + Knockout formats.

## Project Structure

```
bracketforge-api/
├── cmd/api/            # Entrypoint
├── internal/
│   ├── auth/           # JWT helpers
│   ├── config/         # Environment config
│   ├── database/
│   ├── engine/         # Tournament algorithms (core)
│   │   ├── single_elimination.go
│   │   ├── round_robin.go
│   │   └── group_knockout.go
│   ├── handler/        # HTTP handlers (thin layer)
│   ├── middleware/     # Auth, logging
│   ├── model/          # Domain types + DTOs
│   ├── service/        # Business logic
│   └── server/         # Router setup
├── migrations/         # SQL up/down files
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Quick Start

### 1. Prerequisites

- Go 1.22+
- Docker & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI

### 2. Clone and configure

```bash
cp .env.example .env
# Edit .env — at minimum change JWT_SECRET
```

### 3. Start the database

```bash
docker compose up db -d
```

### 4. Run migrations

```bash
migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/bracketforge?sslmode=disable" up
```

### 5. Run the API

```bash
go run ./cmd/api
```

### 6. Verify

```bash
curl http://localhost:8080/health
```

---

## API Reference

### Auth

| Method | Endpoint                | Description             |
| ------ | ----------------------- | ----------------------- |
| POST   | `/api/v1/auth/register` | Create org + admin user |
| POST   | `/api/v1/auth/login`    | Get JWT token           |

### Tournaments (requires `Authorization: Bearer <token>`)

| Method | Endpoint                           | Description           |
| ------ | ---------------------------------- | --------------------- |
| POST   | `/api/v1/tournaments`              | Create tournament     |
| GET    | `/api/v1/tournaments`              | List your tournaments |
| GET    | `/api/v1/tournaments/:id`          | Get tournament        |
| POST   | `/api/v1/tournaments/:id/generate` | Generate bracket      |
| POST   | `/api/v1/tournaments/:id/players`  | Add player            |
| GET    | `/api/v1/tournaments/:id/players`  | List players          |
| GET    | `/api/v1/tournaments/:id/matches`  | List matches          |
| POST   | `/api/v1/matches/:id/result`       | Submit match result   |

### Example: Create and run a tournament

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"org_name":"Lagos Squash Club","email":"admin@example.com","password":"password123"}'

# Login → copy token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'

# Create tournament
curl -X POST http://localhost:8080/api/v1/tournaments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Open Championship 2026","format":"single_elimination"}'

# Add players (repeat for each)
curl -X POST http://localhost:8080/api/v1/tournaments/<id>/players \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Player One"}'

# Generate bracket
curl -X POST http://localhost:8080/api/v1/tournaments/<id>/generate \
  -H "Authorization: Bearer <token>"

# Submit a result
curl -X POST http://localhost:8080/api/v1/matches/<match_id>/result \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"player1_score":3,"player2_score":1}'
```

## Tournament Formats

**Single Elimination** — Players randomly seeded. BYEs auto-handled for non-powers-of-2. Winners auto-advance.

**Round Robin** — Every player plays every other player. N\*(N-1)/2 matches generated. Standings updated after each result.

**Group Stage + Knockout** — Players divided into groups. Round robin within each group. Top players advance to a knockout bracket (to be triggered after group play).

## Next Steps

- [ ] Add `POST /tournaments/:id/advance-to-knockout` to move group stage winners into the knockout bracket
- [ ] Add standings endpoint `GET /tournaments/:id/standings`
- [ ] Add bulk player import
- [ ] Add pagination to list endpoints
- [ ] Write integration tests for each engine
