package post

import (
	"context"
	"database/sql"
	"fmt"
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
			posts.id, 
			authors.id AS authorId, 
			authors.firstname, 
			authors.lastname, 
			authors.email,
			authors.createdAt AS authorCreatedAt, 
			authors.updatedAt AS authorUpdatedAt, 
			description,
			posts.createdAt,
			posts.updatedAt
		FROM posts
		LEFT JOIN authors ON authors.id = posts.authorId
		WHERE posts.id = ?`,
		id,
	)

	post := &Post{}

	err := row.Scan(
		&post.ID,
		&post.Author.ID,
		&post.Author.FirstName,
		&post.Author.LastName,
		&post.Author.Email,
		&post.Author.CreatedAt,
		&post.Author.UpdatedAt,
		&post.Description,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		fmt.Print(err)
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
		fmt.Print(err)
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
			posts.id, 
			authors.id AS authorId, 
			authors.firstname, 
			authors.lastname, 
			authors.email,
			authors.createdAt AS authorCreatedAt, 
			authors.updatedAt AS authorUpdatedAt, 
			description,
			posts.createdAt,
			posts.updatedAt
		FROM posts
		LEFT JOIN authors ON authors.id = posts.authorId`,
	)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	defer results.Close()

	posts := make([]Post, 0)

	for results.Next() {
		var post Post

		results.Scan(
			&post.ID,
			&post.Author.ID,
			&post.Author.FirstName,
			&post.Author.LastName,
			&post.Author.Email,
			&post.Author.CreatedAt,
			&post.Author.UpdatedAt,
			&post.Description,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		posts = append(posts, post)
	}

	return posts, nil
}

func updatePost(id int, updatePost UpdatePost) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := database.DbConn.ExecContext(
		ctx,
		`UPDATE posts SET 
			authorId = ?, 
			description = ? 
		WHERE id = ?`,
		&updatePost.AuthorID,
		&updatePost.Description,
		id,
	)

	if err != nil {
		fmt.Print(err)
		return err
	}

	return nil
}

func insertPost(newPost CreatePost) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := database.DbConn.ExecContext(
		ctx,
		`INSERT INTO posts 
			(
				authorId, 
				description 
			) 
		VALUES (?, ?)`,
		newPost.AuthorID,
		newPost.Description,
	)

	if err != nil {
		return 0, err
	}

	insertID, err := result.LastInsertId()

	if err != nil {
		fmt.Print(err)
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
			posts.id, 
			authors.id AS authorId, 
			authors.firstname, 
			authors.lastname, 
			authors.email,
			authors.createdAt AS authorCreatedAt, 
			authors.updatedAt AS authorUpdatedAt, 
			description,
			posts.createdAt,
			posts.updatedAt
		FROM posts
		LEFT JOIN authors ON authors.id = posts.authorId
		ORDER BY id DESC LIMIT 10`,
	)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	defer results.Close()

	posts := make([]Post, 0)

	for results.Next() {
		var post Post

		results.Scan(
			&post.ID,
			&post.Author.ID,
			&post.Author.FirstName,
			&post.Author.LastName,
			&post.Author.Email,
			&post.Author.CreatedAt,
			&post.Author.UpdatedAt,
			&post.Description,
			&post.CreatedAt,
			&post.UpdatedAt,
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
			posts.id, 
			authors.id AS authorId, 
			authors.firstname, 
			authors.lastname, 
			authors.email,
			authors.createdAt AS authorCreatedAt, 
			authors.updatedAt AS authorUpdatedAt, 
			description,
			posts.createdAt,
			posts.updatedAt
		FROM posts
		LEFT JOIN authors ON authors.id = posts.authorId
		WHERE
	`)

	if postFilter.Author != "" {
		queryBuilder.WriteString(`authors.firstname LIKE ? `)
		queryArgs = append(queryArgs, "%"+postFilter.Author+"%")
	}

	if postFilter.Description != "" {
		if len(queryArgs) > 0 {
			queryBuilder.WriteString(" AND ")
		}

		queryBuilder.WriteString(`description LIKE ? `)
		queryArgs = append(queryArgs, "%"+postFilter.Description+"%")
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
			&post.ID,
			&post.Author.ID,
			&post.Author.FirstName,
			&post.Author.LastName,
			&post.Author.Email,
			&post.Author.CreatedAt,
			&post.Author.UpdatedAt,
			&post.Description,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		posts = append(posts, post)
	}

	return posts, nil
}
