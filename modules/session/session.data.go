package session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kielkow/Post-Service/shared/database"
)

func getPassword(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(
		ctx,
		`SELECT 
			password
		FROM authors
		WHERE email = ?`,
		email,
	)

	var password string

	err := row.Scan(&password)

	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		fmt.Print(err)
		return "", err
	}

	return password, nil
}
