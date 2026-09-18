package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ConfigurePool applies production-grade pooling best practices to any *sql.DB
func ConfigurePool(db *sql.DB) {
	// Maximum number of open connections to the database.
	// Default is 0 (unlimited), which can exhaust database server connection limits.
	db.SetMaxOpenConns(25)

	// Maximum number of connections in the idle connection pool.
	// Keeps warm connections ready to handle traffic spikes.
	db.SetMaxIdleConns(10)

	// Maximum amount of time a connection may be reused.
	// Helps cleanly cycle connections to respect network and load-balancer timeouts.
	db.SetConnMaxLifetime(15 * time.Minute)

	// Maximum amount of time a connection may be idle before being closed.
	db.SetConnMaxIdleTime(5 * time.Minute)
}

func printPoolStats(stats sql.DBStats) {
	fmt.Printf("--- Connection Pool Stats ---\n")
	fmt.Printf("Open Connections : %d\n", stats.OpenConnections)
	fmt.Printf("In Use           : %d\n", stats.InUse)
	fmt.Printf("Idle             : %d\n", stats.Idle)
	fmt.Printf("Wait Count       : %d (Number of queries blocked waiting for a connection)\n", stats.WaitCount)
	fmt.Printf("Wait Duration    : %v (Total duration queries waited for connections)\n", stats.WaitDuration)
}

func main() {
	fmt.Println("=== 1. database/sql Connection Pool Architecture ===")

	// Note: We use a simulated database setup pattern here
	fmt.Println("Configuring database connection pool parameters...")
	// db, err := sql.Open("pgx", "postgres://user:pass@localhost:5432/dbname")
	// ConfigurePool(db)

	fmt.Println(`
Rules for Production Connection Pooling:
1. SetMaxOpenConns must be <= Database server max_connections / number of app replicas.
2. SetMaxIdleConns should be high enough to avoid creating new connections on brief spikes.
3. Always use context timeouts on PingContext, QueryContext, and ExecContext.
`)

	// Simulate stats output
	mockStats := sql.DBStats{
		OpenConnections: 8,
		InUse:           3,
		Idle:            5,
		WaitCount:       0,
		WaitDuration:    0,
	}
	printPoolStats(mockStats)

	// Context timeout pattern
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	fmt.Printf("\nDatabase operations must be bounded by context deadline: %v\n", ctx.Err() == nil)
}
