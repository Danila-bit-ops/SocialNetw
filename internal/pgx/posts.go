package pgx

import (
	"context"
	"fmt"
	"newsSite/internal/model"
)

func (r Repo) GetAllPosts(ctx context.Context) (model.PostsList, error) {
	const q = `SELECT post_id, title, post_text, post_author, author_id, created_at, updated_at FROM posts WHERE 1=1 ORDER BY created_at DESC LIMIT 10`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts model.PostsList
	for rows.Next() {
		var post model.Posts
		if err := rows.Scan(&post.ID, &post.Title, &post.PostText, &post.Author, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			fmt.Println("scan error", err)
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("rows error", err)
		return nil, err
	}
	return posts, nil
}

func (r Repo) AddPost(ctx context.Context, post model.Posts) error {
	const q = `INSERT INTO posts (title, post_text, post_author, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, q, post.Title, post.PostText, post.Author, post.AuthorID, post.CreatedAt, post.UpdatedAt)
	return err
}
