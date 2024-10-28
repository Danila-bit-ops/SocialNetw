package pgx

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"newsSite/internal/model"
)

func (r Repo) RegisterUser(ctx context.Context, user model.User) error {
	// Проверка уникальности логина и email
	const checkUnique = `SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2`
	var count int
	err := r.pool.QueryRow(ctx, checkUnique, user.Username, user.Email).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("пользователь с таким логином или email уже существует")
	}

	// Если уникальны, выполняем вставку
	const q = `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err = r.pool.Exec(ctx, q, user.Username, user.Email, user.Password)
	return err
}

func (r Repo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	row := r.pool.QueryRow(ctx, query, email)
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Ошибка при получении пользователя по email: %v", err)
		return nil, err
	}
	return &user, nil
}
