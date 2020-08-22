package session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kielkow/Post-Service/shared/database"
)

func getPassword(email string) (*Token, error) {
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

	token := &Token{}

	err := row.Scan(&token.Token)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return token, nil
}
