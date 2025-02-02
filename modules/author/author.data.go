package author

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kielkow/Post-Service/shared/database"
)

// GetAuthor function
func GetAuthor(id int) (*Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(
		ctx,
		`SELECT 
				authors.id,
				firstname,
				lastname,
				email,
				avatars.avatar,
				authors.createdAt,
				authors.updatedAt
		FROM
				authors
						LEFT JOIN
				avatars ON avatars.authorId = authors.id
		WHERE
				authors.id = ?;`,
		id,
	)

	author := &Author{}

	err := row.Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Email,
		&author.Avatar,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		fmt.Print(err)
		return nil, err
	}

	if author.Avatar != nil {
		avatarnameJSON, err := json.Marshal(author.Avatar)

		if err != nil {
			fmt.Print(err)
			return nil, err
		}

		avatarname := string(avatarnameJSON[:])

		url := "https://" + os.Getenv("AWS_S3_BUCKET") + ".s3.amazonaws.com/" + avatarname
		url = strings.ReplaceAll(url, string('"'), "")

		author.AvatarURL = &url
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
				authors.id,
				firstname,
				lastname,
				email,
				avatars.avatar,
				authors.createdAt,
				authors.updatedAt
		FROM
				authors
						LEFT JOIN
				avatars ON avatars.authorId = authors.id`,
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
			&author.Avatar,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		if author.Avatar != nil {
			avatarnameJSON, err := json.Marshal(author.Avatar)

			if err != nil {
				fmt.Print(err)
				return nil, err
			}

			avatarname := string(avatarnameJSON[:])

			url := "https://" + os.Getenv("AWS_S3_BUCKET") + ".s3.amazonaws.com/" + avatarname
			url = strings.ReplaceAll(url, string('"'), "")

			author.AvatarURL = &url
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func updateAuthor(id int, author UpdateAuthor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if author.Password != "" {
		_, err := database.DbConn.ExecContext(
			ctx,
			`UPDATE authors SET 
				firstname = ?, 
				lastname = ?, 
				email = ?,
				password = ?
			WHERE id = ? `,
			&author.FirstName,
			&author.LastName,
			&author.Email,
			&author.Password,
			id,
		)

		if err != nil {
			fmt.Print(err)
			return err
		}
	} else {
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
				email,
				password
			) 
		VALUES (?, ?, ?, ?)`,
		newAuthor.FirstName,
		newAuthor.LastName,
		newAuthor.Email,
		newAuthor.Password,
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
				authors.id,
				firstname,
				lastname,
				email,
				avatars.avatar,
				authors.createdAt,
				authors.updatedAt
		FROM
				authors
						LEFT JOIN
				avatars ON avatars.authorId = authors.id
		ORDER BY id DESC
		LIMIT 10;`,
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
			&author.Avatar,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		if author.Avatar != nil {
			avatarnameJSON, err := json.Marshal(author.Avatar)

			if err != nil {
				fmt.Print(err)
				return nil, err
			}

			avatarname := string(avatarnameJSON[:])

			url := "https://" + os.Getenv("AWS_S3_BUCKET") + ".s3.amazonaws.com/" + avatarname
			url = strings.ReplaceAll(url, string('"'), "")

			author.AvatarURL = &url
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func searchAuthorData(authorFilter ReportFilter) ([]Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder

	queryBuilder.WriteString(
		`SELECT
			authors.id,
			firstname,
			lastname,
			email,
			avatars.avatar,
			authors.createdAt,
			authors.updatedAt
		FROM
			authors
					LEFT JOIN
			avatars ON avatars.authorId = authors.id
		WHERE
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
			&author.Avatar,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		if author.Avatar != nil {
			avatarnameJSON, err := json.Marshal(author.Avatar)

			if err != nil {
				fmt.Print(err)
				return nil, err
			}

			avatarname := string(avatarnameJSON[:])

			url := "https://" + os.Getenv("AWS_S3_BUCKET") + ".s3.amazonaws.com/" + avatarname
			url = strings.ReplaceAll(url, string('"'), "")

			author.AvatarURL = &url
		}

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

func removeAvatar(avatar string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
		`DELETE FROM avatars WHERE avatar = ?`, avatar,
	)

	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func getAuthorByEmail(email string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(
		ctx,
		`SELECT 
			email
		FROM
			authors
		WHERE
			email = ?;`,
		email,
	)

	var authorEmail *string

	err := row.Scan(&authorEmail)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return authorEmail, nil
}
