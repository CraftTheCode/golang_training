# Module 07: Database Integration & Application Architecture

## 1. Conceptual Foundation

In Go, database operations are standardized through the `database/sql` package in the standard library. The standard library provides the connection pool, query interfaces, and transaction lifecycle, while specific database drivers (like `github.com/jackc/pgx/v5/stdlib` or `github.com/lib/pq` for PostgreSQL) register themselves under the hood.

Key tenets:
1. **Always use parameterized queries (`$1`, `$2`)**: Never concatenate strings to build SQL queries.
2. **Always pass `context.Context`**: Use `QueryContext`, `ExecContext`, and `BeginTx` with context timeouts to prevent deadlocks and runaway queries.
3. **Always close `*sql.Rows`**: Failing to close rows leaks database pool connections.
4. **Tune the connection pool**: Default settings in Go allow unlimited open connections, which can overwhelm PostgreSQL.

---

## 2. Architecture: Repository Pattern & Clean Architecture

```text
HTTP Transport (Handlers)
        |
        v
Application Service Layer (Domain Logic & Validation)
        |
        v
Repository Interface (TaskRepository, UserRepository)
       / \
      /   \
Postgres Implementation    In-Memory Implementation (For fast unit testing)
```

By coding against repository interfaces, services can be tested with instantaneous in-memory fakes without spinning up Docker or cloud databases.

---

## 3. Database Connection Pool Tuning

```go
db, err := sql.Open("pgx", connString)

// 1. Max Open Connections: Limits total simultaneous database sockets
db.SetMaxOpenConns(25)

// 2. Max Idle Connections: Number of warm idle connections retained in pool
db.SetMaxIdleConns(10)

// 3. Max Lifetime: Closes connection after duration to prevent stale load-balancer routing
db.SetConnMaxLifetime(15 * time.Minute)

// 4. Max Idle Time: Closes connections that have been unused for too long
db.SetConnMaxIdleTime(5 * time.Minute)
```

---

## 4. Module Directory Structure

```text
07-database-and-architecture/
├── 01_connection_pools/main.go         # Pool sizing, lifecycle, connection testing
├── 02_transactions_and_errors/main.go  # Atomic BeginTx, Commit, Rollback idiom
└── project-postgres-api/               # Layered Architecture PostgreSQL Project
    ├── db/migrations/001_init.sql      # Schema migration
    ├── internal/
    │   ├── domain/models.go            # Domain entities and Repository interface
    │   ├── repository/postgres.go      # PostgreSQL database implementation
    │   └── repository/memory.go        # Fast In-memory fake for tests
    └── README.md                       # Setup with Docker Postgres
```
