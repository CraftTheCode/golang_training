# Project: Layered PostgreSQL REST API

This project demonstrates clean layered architecture using the **Repository Pattern**:
- Domain interfaces are decoupled from underlying database implementations.
- Includes both a real `PostgresRepository` using parameterized SQL queries and an in-memory `MemoryRepository` for instant unit tests.

## Database Migration
Execute `db/migrations/001_init.sql` to initialize the PostgreSQL schema.

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## Running with Docker PostgreSQL
```bash
docker run --name go-postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=gotraining -p 5432:5432 -d postgres:16-alpine
```
