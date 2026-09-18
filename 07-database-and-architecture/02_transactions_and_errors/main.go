package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TransferFunds demonstrates the canonical Go transaction pattern.
// Notice the 'defer tx.Rollback()' idiom!
// If tx.Commit() succeeds, the subsequent deferred tx.Rollback() returns sql.ErrTxDone and is safely ignored.
func TransferFunds(ctx context.Context, db *sql.DB, fromAccount, toAccount int, amount float64) error {
	// 1. Begin atomic transaction with explicit isolation level and context
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Canonical safety net: will rollback if any step returns early with error!
	defer tx.Rollback()

	// 2. Deduct from source account
	res, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, fromAccount)
	if err != nil {
		return fmt.Errorf("debit failed: %w", err)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("insufficient funds or account not found")
	}

	// 3. Credit destination account
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toAccount)
	if err != nil {
		return fmt.Errorf("credit failed: %w", err)
	}

	// 4. Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("=== Canonical Go Transaction Pattern ===")
	fmt.Println(`
Key Rules for database/sql Transactions:
1. Always call 'defer tx.Rollback()'.
   - If you return with an error anywhere in the function, it rolls back.
   - If 'tx.Commit()' succeeds, 'defer tx.Rollback()' is a safe no-op.
2. Use 'tx.ExecContext' and 'tx.QueryRowContext' (NOT 'db.ExecContext' inside a transaction!).
3. Set transaction isolation levels when needed (e.g. sql.LevelSerializable).
`)

	// Context with timeout prevents hung transactions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("Transaction context configured with 5s timeout. Done: %v\n", ctx.Err())
}
