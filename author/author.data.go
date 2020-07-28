package author

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kielkow/Post-Service/database"
)

// GetAuthor function
func GetAuthor(id int) (*Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(
		ctx,
		`SELECT 
			id, 
			firstname, 
			lastname, 
			email,
			createdAt,
			updatedAt
		FROM authors
		WHERE id = ?`,
		id,
	)

	author := &Author{}

	err := row.Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Email,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return author, nil
}

func removeAuthor(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
		`DELETE FROM authors WHERE id = ?`, id,
	)

	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func getAuthorList() ([]Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(
		ctx,
		`SELECT 
			id, 
			firstname, 
			lastname, 
			email,
			createdAt,
			updatedAt
		from authors`,
	)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	defer results.Close()

	authors := make([]Author, 0)

	for results.Next() {
		var author Author

		results.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Email,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		authors = append(authors, author)
	}

	return authors, nil
}

func updateAuthor(id int, author Author) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
		`UPDATE authors SET 
			firstname = ?, 
			lastname = ?, 
			email = ? 
		WHERE id = ? `,
		&author.FirstName,
		&author.LastName,
		&author.Email,
		id,
	)

	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func insertAuthor(newAuthor CreateAuthor) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(
		ctx,
		`INSERT INTO authors 
			(
				firstname, 
				lastname,
				email 
			) 
		VALUES (?, ?, ?)`,
		newAuthor.FirstName,
		newAuthor.LastName,
		newAuthor.Email,
	)

	if err != nil {
		fmt.Print(err)
		return 0, err
	}

	insertID, err := result.LastInsertId()

	if err != nil {
		fmt.Print(err)
		return 0, err
	}

	return int(insertID), err
}

func getToptenAuthors() ([]Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(
		ctx,
		`SELECT 
			id, 
			firstname, 
			lastname, 
			email, 
			createdAt, 
			updatedAt 
		from authors ORDER BY id DESC LIMIT 10`,
	)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	defer results.Close()

	authors := make([]Author, 0)

	for results.Next() {
		var author Author

		results.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Email,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		authors = append(authors, author)
	}

	return authors, nil
}

func searchAuthorData(authorFilter ReportFilter) ([]Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder

	queryBuilder.WriteString(`SELECT
		id,
		firstname,
		lastname,
		email,
		createdAt,
		updatedAt
		FROM authors WHERE
	`)

	if authorFilter.FirstName != "" {
		queryBuilder.WriteString(`firstname LIKE ? `)
		queryArgs = append(queryArgs, "%"+authorFilter.FirstName+"%")
	}

	if authorFilter.LastName != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}

		queryBuilder.WriteString(`lastname LIKE ? `)
		queryArgs = append(queryArgs, "%"+authorFilter.LastName+"%")
	}

	results, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer results.Close()

	authors := make([]Author, 0)

	for results.Next() {
		var author Author

		results.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Email,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		authors = append(authors, author)
	}

	return authors, nil
}

func createAvatar(newAvatar CreateAuthorAvatar) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(
		ctx,
		`INSERT INTO avatars 
			(
				authorId, 
				avatar
			) 
		VALUES (?, ?)`,
		newAvatar.AuthorID,
		newAvatar.Avatar,
	)

	if err != nil {
		fmt.Print(err)
		return 0, err
	}

	insertID, err := result.LastInsertId()

	if err != nil {
		fmt.Print(err)
		return 0, err
	}

	return int(insertID), err
}
