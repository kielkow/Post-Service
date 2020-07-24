package post

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/kielkow/Post-Service/database"
)

func getPost(id int) (*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := database.DbConn.QueryRowContext(
		ctx,
		`SELECT 
			id, 
			author, 
			description,
			createdAt,
			updatedAt
		FROM posts
		WHERE id = ?`,
		id,
	)

	post := &Post{}

	err := row.Scan(
		&post.id,
		&post.author,
		&post.description,
		&post.createdAt,
		&post.updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return post, nil
}

func removePost(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
		`DELETE FROM posts WHERE id = ?`, id,
	)

	if err != nil {
		return err
	}

	return nil
}

func getPostList() ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(
		ctx,
		`SELECT 
			id, 
			author, 
			description,
			createdAt,
			updatedAt
		from posts`,
	)

	if err != nil {
		return nil, err
	}

	defer results.Close()

	posts := make([]Post, 0)

	for results.Next() {
		var post Post

		results.Scan(
			&post.id,
			&post.author,
			&post.description,
			&post.createdAt,
			&post.updatedAt,
		)

		posts = append(posts, post)
	}

	return posts, nil
}

func updatePost(post Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
		`UPDATE posts SET 
			author = ?, 
			description = ?, 
			updatedAt = ?
		WHERE id = ? `,
		post.author,
		post.description,
		time.Now(),
		post.id,
	)

	if err != nil {
		return err
	}

	return nil
}

func insertPost(post Post) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(
		ctx,
		`INSERT INTO posts 
			(
				author, 
				description, 
				createdAt, 
				updatedAt
			) 
		VALUES (?, ?, ?, ?)`,
		post.author,
		post.description,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	insertID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(insertID), err
}

func getToptenPosts() ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := database.DbConn.QueryContext(
		ctx,
		`SELECT 
			id, 
			author, 
			description, 
			createdAt, 
			updatedAt 
		from posts ORDER BY id DESC LIMIT 10`,
	)

	if err != nil {
		return nil, err
	}

	defer results.Close()

	posts := make([]Post, 0)

	for results.Next() {
		var post Post

		results.Scan(
			&post.id,
			&post.author,
			&post.description,
			&post.createdAt,
			&post.updatedAt,
		)

		posts = append(posts, post)
	}

	return posts, nil
}

func searchPostData(postFilter ReportFilter) ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var queryArgs = make([]interface{}, 0)
	var queryBuilder strings.Builder

	queryBuilder.WriteString(`SELECT
		id,
		author,
		description,
		createdAt, 
		updatedAt
		FROM posts WHERE
	`)

	if postFilter.author != "" {
		queryBuilder.WriteString(`author LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(postFilter.author)+"%")
	}

	if postFilter.description != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}

		queryBuilder.WriteString(`description LIKE ? `)
		queryArgs = append(queryArgs, "%"+strings.ToLower(postFilter.description)+"%")
	}

	results, err := database.DbConn.QueryContext(ctx, queryBuilder.String(), queryArgs...)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer results.Close()

	posts := make([]Post, 0)

	for results.Next() {
		var post Post

		results.Scan(
			&post.id,
			&post.author,
			&post.description,
			&post.createdAt,
			&post.updatedAt,
		)

		posts = append(posts, post)
	}

	return posts, nil
}
