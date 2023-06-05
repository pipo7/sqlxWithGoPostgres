package sqlxpq // this is any name ...this is not related to sqlx

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// constant query to be created
const createUsersTableQuery = `CREATE TABLE IF NOT EXISTS users (
   id SERIAL PRIMARY KEY,
   username VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL
);`

// user struct to pass and save data in it
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
}

func connectToDatabase() (*sqlx.DB, error) {
	host := "localhost"  // please replace it with your database host
	port := 5432         // please replace it with your database port
	user := "pipo"       // please replace it with your database user
	pwd := "mypwdtest"   // please replace it with your database password
	dbName := "practice" // please replace it with your database name
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user,
		pwd,
		host,
		port,
		dbName,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createUsersTableQuery)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	return db, nil
}

func insertUser(tx *sqlx.Tx, user *User) error {
	query := `INSERT INTO users (username, email) VALUES (:username, :email) RETURNING id`
	rows, err := tx.NamedQuery(query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user.ID)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}
