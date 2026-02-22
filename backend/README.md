# openMSG — Backend

Go + Fiber API server. Stores users, conversations, and encrypted messages.
**Untrusted by design** — never sees plaintext, never decrypts anything.
All data is publicly readable. Encryption lives entirely in the browser.

## API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/users` | Register user (nickname, public_key, avatar) |
| `GET` | `/users` | List all users |
| `GET` | `/users/:id` | Get user by ID |
| `GET` | `/users/:userId/conversations` | List conversations for a user |
| `POST` | `/conversations` | Get or create conversation between two users |
| `GET` | `/conversations/:id` | Get single conversation by ID |
| `POST` | `/messages` | Send encrypted message (ciphertext + nonce only) |
| `GET` | `/messages/:conversationId` | List messages in a conversation |
| `GET` | `/public/messages` | Public feed: sender, recipient, ciphertext, timestamp |

## Stack

| Layer | Tech |
|-------|------|
| HTTP | Fiber v2 |
| DB driver | pgx v5 + pgxpool |
| Query layer | sqlc (generated in `internal/repository/sqlcgen/`) |
| Testing | testify + minimock |
| Config | env vars (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `SERVER_PORT`) |

## Commands

```sh
task build      # compile binary → backend/bin/server
task test       # run unit tests
task sqlc       # regenerate DB query code after editing queries.sql
task generate   # regenerate minimock mocks after changing repo interfaces
```
