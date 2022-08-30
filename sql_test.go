package golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
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

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()

	query := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var id, name string

		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Printf("ID: %s, Name: %s\n", id, name)
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()

	query := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool

		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("===========================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
		if email.Valid {
			fmt.Println("Email:", email.String)
		}
		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		if birthDate.Valid { //
			fmt.Println("BirthDate:", birthDate.Time)
		}
		fmt.Println("Married:", married)
		fmt.Println("CreatedAt:", createdAt)
	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()
	username := "admin'; #"
	password := "notAdmin"

	query := "SELECT username FROM user WHERE username = '" + username +
		"' AND password = '" + password + "' LIMIT 1"

	fmt.Println(query)

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Username:", username)
	} else {
		fmt.Println("Login failed")
	}
}

func TestSqlInjectionSolution(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()
	username := "admin'; #"
	password := "notAdmin"

	query := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"

	fmt.Println(query)

	rows, err := db.QueryContext(ctx, query, username, password)

	if err != nil {
		panic(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("Username:", username)
	} else {
		fmt.Println("Login failed")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()

	username := "user'; DROP TABLE user; #"
	password := "user"

	query := "INSERT INTO user(username, password) VALUES(?, ?)"
	_, err := db.ExecContext(ctx, query, username, password)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully executed query")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()

	email := "bob@gmail.com"
	comment := "Test Comment"

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, query, email, comment)

	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id:", insertId)
}

func TestPreparedStatement(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {

		}
	}(statement)

	for i := 0; i < 10; i++ {
		email := fmt.Sprintf("bob%s@gmail.com", strconv.Itoa(i))
		comment := fmt.Sprintf("Test Comment %s", strconv.Itoa(i))

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println(id)
	}
}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	statement, err := tx.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {

		}
	}(statement)

	//	Do transaction
	for i := 0; i < 10; i++ {
		email := fmt.Sprintf("bob%s@gmail.com", strconv.Itoa(i))
		comment := fmt.Sprintf("Test Comment %s", strconv.Itoa(i))

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println(id)
	}

	err = tx.Commit()
	//err = tx.Rollback() // cancel/rollback transaction
	if err != nil {
		panic(err)
	}
}
