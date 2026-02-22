# Detailed Implementation Plan

**Project:** Public Encrypted Messaging App
**Stack:** Go (Fiber) + Next.js + TanStack Query + PostgreSQL (Supabase)
**Infra:** Docker Compose (local), Vercel (frontend), Alwaysdata (Go API)
**Crypto:** Browser-only E2EE (WebCrypto)
**Public Mode:** Metadata + Ciphertext visible

---

# 0. Architectural Summary

### Core Principle

- Database is fully readable.
- Backend is untrusted.
- Only clients decrypt.
- Metadata (sender, recipient, timestamp) is public.
- Message content remains encrypted.

### High-Level Components

| Layer     | Technology                                            |
| --------- | ----------------------------------------------------- |
| Frontend  | Next.js (App Router) + TanStack Query                 |
| Backend   | Go + Fiber                                            |
| DB        | PostgreSQL (Supabase prod / local Postgres in Docker) |
| Crypto    | WebCrypto API                                         |
| Testing   | testify + minimock                                    |
| DB Driver | pgx v5                                                |
| Config    | yaml.v3                                               |

---

# 1. Installation Phase (Local Development Setup)

This must be completed before any development.

## 1.1 Required Software

Install:

- Go 1.22+
- Node 20+
- pnpm or npm
- Docker
- Docker Compose v2
- Make (optional but recommended)

## 1.2 Go Dependencies (Do Not Implement Business Code)

Required Go modules:

```
github.com/gofiber/fiber/v2
github.com/jackc/pgx/v5
github.com/stretchr/testify
github.com/gojuno/minimock/v3
gopkg.in/yaml.v3
```

Initialize backend:

```
go mod init app
go get ...
```

## 1.3 Frontend Setup

```
npx create-next-app@latest frontend
```

Install:

```
@tanstack/react-query
```

---

# 2. Docker Compose (MANDATORY)

Create a root-level `docker-compose.yml`.

## 2.1 Services Required

### postgres

- image: postgres:15
- expose 5432
- volume for persistence
- env:
  - POSTGRES_USER
  - POSTGRES_PASSWORD
  - POSTGRES_DB

### backend

- build: ./backend
- depends_on: postgres
- env:
  - DB_HOST=postgres
  - DB_PORT=5432
  - DB_USER
  - DB_PASSWORD
  - DB_NAME

- expose 8080

### frontend

- build: ./frontend
- depends_on: backend
- expose 3000
- env:
  - NEXT_PUBLIC_API_URL=[http://backend:8080](http://backend:8080)

## 2.2 Backend Dockerfile Requirements

- Multi-stage build
- Build static Go binary
- Use minimal runtime image (distroless or alpine)
- Expose 8080

## 2.3 Frontend Dockerfile

- Node base
- Install deps
- Build
- Start with `next start`

---

# 3. Database Design

## 3.1 Tables

### users

- id (uuid, pk)
- nickname (unique, indexed)
- public_key (text)
- created_at (timestamp)

### conversations

- id (uuid, pk)
- user_a_id (uuid, indexed)
- user_b_id (uuid, indexed)
- created_at
- UNIQUE(user_a_id, user_b_id)

Store sorted user IDs to enforce uniqueness.

### messages

- id (uuid)
- conversation_id (uuid, indexed)
- sender_id (uuid)
- ciphertext (text)
- nonce (text)
- created_at (timestamp)

## 3.2 Public Exposure Policy

All read endpoints:

- No authentication
- Fully readable

---

# 4. Backend Structure (Fiber)

## 4.1 Project Structure

```
backend/
  cmd/
    server/
  internal/
    config/
    db/
    repository/
    service/
    handler/
    models/
    middleware/
  test/
```

## 4.2 Configuration

Use `yaml.v3`:

config.yaml:

```
server:
  port: 8080

database:
  host: postgres
  port: 5432
  user: ...
  password: ...
  name: ...
```

Load on startup.

---

# 5. Database Layer (pgx v5)

## 5.1 Use pgxpool

- Initialize pool on startup
- Inject into repositories

## 5.2 Repository Pattern Required

Interfaces must exist:

```
UserRepository
ConversationRepository
MessageRepository
```

Implementations:

- PostgresUserRepository
- PostgresConversationRepository
- PostgresMessageRepository

No SQL in handlers.

---

# 6. Testing Requirements

## 6.1 Use testify

For:

- Assertions
- Integration tests

## 6.2 Use minimock

Generate mocks for:

- Repository interfaces
- Service layer tests

Tests required for:

- User creation logic
- Conversation uniqueness logic
- Message creation

No business logic inside handlers.

---

# 7. API Design (Fiber)

All JSON.

## 7.1 Users

POST /users

- nickname
- public_key

GET /users

- return all users

GET /users/:id

---

## 7.2 Conversations

POST /conversations

- user_a_id
- user_b_id

Server:

- Sort IDs
- Create if not exists
- Return existing if exists

GET /conversations/:userId

- return conversations for user

---

## 7.3 Messages

POST /messages

- conversation_id
- sender_id
- ciphertext
- nonce

GET /messages/:conversationId

---

## 7.4 Public Discover Endpoint

GET /public/messages

Returns:

- sender nickname
- recipient nickname
- timestamp
- ciphertext

No plaintext ever.

---

# 8. Frontend Architecture

## 8.1 Structure

```
frontend/
  app/
    register/
    dashboard/
    chat/[conversationId]/
    discover/
  lib/
    crypto.ts
    api.ts
  providers/
```

---

# 9. TanStack Query Integration

## 9.1 QueryClient at Root

Wrap app with QueryClientProvider.

## 9.2 Required Queries

- useUsers
- useConversations
- useMessages
- usePublicMessages

Refetch:

- 3–5 seconds polling

No websockets.

---

# 10. Crypto Design (Browser Only)

## 10.1 Key Generation

On registration:

- Generate X25519 keypair
- Store private key in IndexedDB
- Send public key to backend

## 10.2 Shared Secret

On send:

- Fetch recipient public key
- Derive shared secret
- HKDF derive AES key
- Encrypt with AES-GCM

## 10.3 Decryption

On fetch:

- Derive shared secret
- Decrypt

Backend never decrypts.

---

# 11. Registration Model

Nickname-only.

Flow:

1. Enter nickname.
2. Generate keypair.
3. POST /users.
4. Store user_id locally.

No auth.
No JWT.
No sessions.

---

# 12. Discover Page Behavior

Shows:

- sender nickname
- recipient nickname
- timestamp
- ciphertext

No decryption attempted.

This intentionally leaks metadata.

---

# 13. Deployment Plan

## 13.1 Supabase (Production DB)

- Use only Postgres.
- No Supabase auth.
- Fiber connects via pgx.

## 13.2 Alwaysdata

- Deploy built Go binary.
- Provide env variables.
- Ensure outbound connection to Supabase allowed.

## 13.3 Vercel

- Deploy Next.js.
- Set NEXT_PUBLIC_API_URL to Alwaysdata endpoint.

---

# 14. Constraints and Explicit Non-Goals

Not implemented:

- Authentication
- Key rotation
- Forward secrecy
- Multi-device sync
- Message signatures
- Rate limiting

---

# 15. Deliverables Expected from Other LLM

The implementing LLM must:

1. Generate:
   - docker-compose.yml
   - Dockerfiles

2. Scaffold backend with:
   - Config loading
   - Fiber setup
   - pgx pool
   - Repository pattern

3. Implement tests using:
   - testify
   - minimock

4. Scaffold Next.js app with:
   - TanStack Query
   - Crypto helpers

5. Implement API contracts only (no advanced business logic).

No production hardening required.

---

# 16. Final Architectural Position

This project demonstrates:

- Public database model
- Zero-trust backend
- Client-side encryption
- Metadata transparency
- Clean Go layered architecture
- Testable service layer with mocks

The system is intentionally simple but architecturally disciplined.

This plan is now sufficient for another LLM to execute without requiring additional clarification.
