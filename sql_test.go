package golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()

	query := "INSERT INTO customer(id, name) VALUES('alice', 'Alice')"
	_, err := db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully executed query")

}
