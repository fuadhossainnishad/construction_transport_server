---

## 🚛 Construction Transport Server

**Event‑driven, real‑time marketplace for trucking & logistics**  
Built with Go, Gin, PostgreSQL, Redis, RabbitMQ, WebSocket, Stripe Connect.

---

## 📖 Table of Contents

- [Overview](#overview)
- [System Architecture](#system-architecture)
- [Tech Stack](#tech-stack)
- [Folder Structure](#folder-structure)
- [Core Flows](#core-flows)
- [Event‑Driven Design](#event‑driven-design)
- [API Documentation](#api-documentation)
- [Setup & Installation](#setup--installation)
- [Environment Variables](#environment-variables)
- [Running the Server](#running-the-server)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

This platform connects **customers** (shippers) with **transporters** (truck owners/drivers).  
Key features:

- **User roles** – `USER`, `TRANSPORTER`, `ADMIN`
- **Authentication** – JWT + OTP email verification (RabbitMQ + SMTP)
- **Vehicle management** – transporters can add/update their trucks
- **Booking** – customers create a job; automatically assigned to a transporter
- **Real‑time job tracking** – WebSocket push notifications on status changes
- **Event‑driven** – all async tasks (email, notifications, payments) via RabbitMQ
- **Payment** – automatic Stripe Connect transfers when a job is completed
- **Clean architecture** – domain, repository, usecase, delivery layers

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client (Mobile/Web)                      │
└─────────────┬───────────────────────────┬───────────────────────┘
              │ REST (Gin)                 │ WebSocket
              ▼                            ▼
┌─────────────────────────┐      ┌─────────────────────────┐
│      API Gateway        │      │    WebSocket Hub        │
│      (Gin Router)       │◄────►│    (Real-time pushes)    │
└─────────────┬───────────┘      └─────────────┬───────────┘
              │                                  │
              ▼                                  │
┌─────────────────────────────────────────────┐  │
│           Application Layer                  │  │
│  Handlers → Usecases → Repositories          │  │
└─────────────┬───────────────────────────────┘  │
              │                                   │
    ┌─────────┼─────────┬─────────────┐          │
    ▼         ▼         ▼             ▼          │
┌──────┐ ┌──────┐ ┌──────────┐ ┌────────────┐    │
│Postgres│Redis│ │RabbitMQ │ │ Stripe API │    │
└───┬──┘ └──┬──┘ └────┬─────┘ └────────────┘    │
    │       │         │                          │
    │       │    ┌────▼─────┐                    │
    │       │    │  Worker  │ (email consumer)   │
    │       │    └──────────┘                    │
    │       │                                    │
    └───────┴────────────────────────────────────┘
```

**Data flow example (job status update):**

1. Transporter `PATCH /jobs/:id/status` → Gin handler
2. Usecase updates DB, adds timeline entry
3. Usecase publishes `job.status.updated` event to RabbitMQ
4. Usecase sends WebSocket message to customer (real‑time)
5. If status = `completed`, usecase calls Stripe transfer (async)

---

## Tech Stack

| Component      | Technology                               |
| -------------- | ---------------------------------------- |
| Language       | Go 1.25+                                 |
| Web Framework  | Gin                                      |
| Database       | PostgreSQL 16 (pgxpool)                  |
| Cache / OTP    | Redis 7                                  |
| Message Broker | RabbitMQ 3 (with management plugin)      |
| Real‑time      | Gorilla WebSocket                        |
| Payments       | Stripe Connect API v79                   |
| Email          | SMTP (Gmail / any) via RabbitMQ consumer |
| Migrations     | Manual SQL (will add tool later)         |
| Dev tool       | Air (live reload)                        |
| Container      | Docker + Docker Compose                  |

---

## Folder Structure

```
construction_transport_server/
├── api/
│   └── rest/v1/
│       ├── delivery/          # Gin handlers (auth, account, vehicle, booking, job)
│       ├── dto/               # request/response structs
│       ├── middleware/        # JWT auth, role checks
│       └── routes.go          # route registration
├── cmd/
│   ├── api/main.go            # server entry point
│   └── worker/main.go         # (future) separate worker process
├── config/                    # configuration loaders (DB, app)
├── infrastructure/
│   ├── cache/redis/           # OTP store, rate limiter
│   ├── database/postgres/     # connection, retry, pool config
│   └── messaging/rabbitmq/    # publisher, consumer base
├── internal/                  # domain logic (no external deps)
│   ├── account/               # user profile
│   ├── auth/                  # authentication, sessions, OTP
│   ├── booking/               # booking domain & repository
│   ├── events/                # event types & interfaces
│   ├── job/                   # job status management
│   ├── notification/          # email consumer (SMTP)
│   ├── vehicle/               # truck CRUD
│   └── websocket/             # hub & client handling
├── pkg/                       # shared utilities (logger, JWT, metrics)
├── docker-compose.yml
├── Dockerfile
├── .air.toml
├── .env.example
├── Makefile
└── README.md
```

---

## Core Flows

### 1. User Registration & OTP Verification

```
User → POST /auth/register → DB insert (pending) → RabbitMQ "user.registered"
→ email consumer → sends OTP → User POST /auth/verify-otp → email verified.
```

### 2. Create a Booking (Customer)

```
Customer (USER) → POST /api/v1/bookings → validate vehicle → create booking (pending)
→ publish "booking.created" → transporter sees new job in GET /jobs.
```

### 3. Job Status Lifecycle (Transporter)

```
assigned → heading_to_pickup → arrived_at_pickup → loaded → in_transit → delivered → completed
```

Each transition:

- Updates DB and timeline
- Sends WebSocket push to customer
- Publishes `job.status.updated` event
- On `completed` → Stripe transfer to transporter

### 4. Real‑time Tracking

Customer opens WebSocket connection: `ws://server/ws?user_id=<id>`  
Transporter updates status → server pushes JSON message to customer’s WebSocket.

---

## Event‑Driven Design

All async side effects go through RabbitMQ. **No blocking calls** inside HTTP handlers.

| Event                | Trigger                    | Consumer Action                      |
| -------------------- | -------------------------- | ------------------------------------ |
| `user.registered`    | User signs up              | Send OTP email via SMTP              |
| `booking.created`    | Customer creates booking   | Notify transporter (push, SMS)       |
| `job.status.updated` | Transporter changes status | Log analytics, trigger notifications |
| `job.completed`      | Job reaches `completed`    | Generate invoice, update earnings    |

Events are published **inside usecases** using an `events.EventPublisher` interface.  
The RabbitMQ implementation is injected at startup – easy to swap for Kafka or Google Pub/Sub later.

---

## API Documentation

All endpoints return JSON with consistent structure:

```json
{
  "status_code": 200,
  "message": "success",
  "data": { ... }
}
```

### Authentication (public)

| Method | Endpoint                | Description                        |
| ------ | ----------------------- | ---------------------------------- |
| POST   | `/auth/register`        | Create account, OTP sent via email |
| POST   | `/auth/login`           | Returns access + refresh tokens    |
| POST   | `/auth/verify-otp`      | Verify email with OTP              |
| POST   | `/auth/forgot-password` | Send OTP to reset password         |
| POST   | `/auth/reset-password`  | Reset password (requires OTP)      |
| POST   | `/auth/refresh`         | Get new access token               |

### Authenticated (`/api/v1`)

| Method | Endpoint               | Role        | Description                |
| ------ | ---------------------- | ----------- | -------------------------- |
| GET    | `/profile`             | ANY         | Get own profile            |
| PUT    | `/profile`             | ANY         | Update profile             |
| POST   | `/vehicles`            | TRANSPORTER | Add a truck                |
| GET    | `/vehicles`            | TRANSPORTER | List my vehicles           |
| POST   | `/bookings`            | USER        | Create a booking           |
| GET    | `/bookings`            | USER        | List my bookings           |
| GET    | `/bookings/:id`        | USER/ADMIN  | Get booking details        |
| POST   | `/bookings/:id/cancel` | USER        | Cancel a pending booking   |
| GET    | `/jobs`                | TRANSPORTER | List assigned jobs         |
| GET    | `/jobs/:id`            | TRANSPORTER | Get job details + timeline |
| PATCH  | `/jobs/:id/status`     | TRANSPORTER | Update job status          |

### WebSocket

`ws://localhost:8080/ws?user_id={user_id}`  
Receives JSON messages:

```json
{
  "event": "job_status_updated",
  "booking_id": 123,
  "status": "heading_to_pickup",
  "notes": "on the way"
}
```

---

## Setup & Installation

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Make (optional)
- Stripe test account (for payouts)
- SMTP credentials (Gmail app password)

### 1. Clone the repo

```bash
git clone https://github.com/your-org/construction_transport_server.git
cd construction_transport_server
```

### 2. Environment variables

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

Edit `.env` with your database, Redis, RabbitMQ, Stripe, and SMTP settings.

### 3. Run with Docker Compose

```bash
docker-compose up --build
```

This starts:

- PostgreSQL (port 5432)
- Redis (6379)
- RabbitMQ (5672, management UI on 15672)
- Go app (port 8080)

### 4. Run database migrations

Migrations are **manual** for now. Apply `infrastructure/database/migrations/001_init.up.sql` to your PostgreSQL instance:

```bash
docker exec -i construction_transport_db psql -U admin -d construction_transport < infrastructure/database/migrations/001_init.up.sql
```

### 5. Verify server is running

```bash
curl http://localhost:8080/auth/register -X POST -d '{"email":"test@example.com","password":"123456","role":"USER"}' -H "Content-Type: application/json"
```

You should receive an OTP email (check console logs if SMTP not configured).

---

## Environment Variables

| Variable            | Description                     | Default                           |
| ------------------- | ------------------------------- | --------------------------------- |
| `PORT`              | HTTP listen port                | 8080                              |
| `DB_HOST`           | PostgreSQL host                 | postgres                          |
| `DB_PORT`           | PostgreSQL port                 | 5432                              |
| `DB_USER`           | PostgreSQL user                 | admin                             |
| `DB_PASSWORD`       | PostgreSQL password             | admin_secret                      |
| `DB_NAME`           | Database name                   | construction_transport            |
| `DB_SSLMODE`        | SSL mode (disable for local)    | disable                           |
| `REDIS_ADDR`        | Redis address                   | redis:6379                        |
| `RABBITMQ_URL`      | RabbitMQ connection URL         | amqp://guest:guest@rabbitmq:5672/ |
| `STRIPE_SECRET_KEY` | Stripe secret key (test/live)   | (required for payouts)            |
| `JWT_SECRET`        | Secret for signing JWT tokens   | (change in prod)                  |
| `SMTP_HOST`         | SMTP server host                | smtp.gmail.com                    |
| `SMTP_PORT`         | SMTP port                       | 587                               |
| `SMTP_FROM`         | Sender email address            |                                   |
| `SMTP_PASSWORD`     | SMTP password (or app password) |                                   |

---

## Running the Server (without Docker)

```bash
# Install dependencies
go mod download

# Run with Air (live reload)
make dev

# Or build and run
make build
./bin/app
```

---

## Testing

```bash
# Run all tests
make test

# Run integration tests (requires DB)
go test -tags=integration ./tests/integration/...
```

---

## Deployment

### Production checklist

1. Use **managed PostgreSQL** (RDS, Cloud SQL) with SSL
2. Use **managed Redis** (ElastiCache, Memorystore)
3. Use **managed RabbitMQ** (CloudAMQP)
4. Set `JWT_SECRET` to a strong random value
5. Use **Stripe Live keys** and verify webhook signatures
6. Run **two instances** of the API behind a load balancer
7. Use **environment‑specific** config (e.g., `config.prod.yaml`)
8. Set up **database migrations** as part of CI/CD (e.g., `golang-migrate`)

### Docker production build

```dockerfile
# Multi‑stage build example (add to Dockerfile)
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

---

## Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (conventional commits)
4. Push and open a Pull Request

---

## License

MIT © 2025 – Built for portfolio & job application.

---

## 📬 Contact

For questions or job opportunities:  
**Fuad Hossain** – [fuadhossainbk01@gmail.com] – [LinkedIn/(https://www.linkedin.com/in/fuadhossain01/)]

---
