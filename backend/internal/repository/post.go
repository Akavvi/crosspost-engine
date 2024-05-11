package repository

import (
	"blog-app/internal/models"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewPostRepository(db *sqlx.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

type PostRepository struct {
	db *sqlx.DB
}

func (r *PostRepository) Save(ctx context.Context, timeout time.Duration, post *models.Post) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	p := &models.Post{}

	query := "INSERT INTO posts (title, content) VALUES ($1, $2) returning *;"
	tx := r.db.MustBegin()
	tx.QueryRowxContext(ctx, query, post.Title, post.Content).StructScan(p)
	err := tx.Commit()
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PostRepository) Delete(ctx context.Context, timeout time.Duration, id int) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	query := "DELETE FROM posts WHERE id = $1"
	tx := r.db.MustBegin()
	tx.QueryRowxContext(ctx, query, id)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) GetAll(ctx context.Context, timeout time.Duration) ([]*models.Post, error) {
	var posts []*models.Post
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	query := "SELECT * FROM posts"
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := &models.Post{}
		err := rows.StructScan(p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) Get(ctx context.Context, timeout time.Duration, id int) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	dest := &models.Post{}
	query := "SELECT * FROM posts WHERE id = $1"
	err := r.db.QueryRowxContext(ctx, query, id).StructScan(dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
