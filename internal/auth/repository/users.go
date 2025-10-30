package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mSulimenko/dev-blog-platform/internal/auth/models"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewUsersRepository(pool *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{
		db: pool,
	}
}

func (u *UsersRepository) CreateUser(ctx context.Context, user *models.User) error {
	q := `
		INSERT INTO users (username, email, password_hash, verification_token) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
`

	err := u.db.QueryRow(ctx, q, user.Username, user.Email, user.PasswordHash, user.VerificationToken).
		Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (u *UsersRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	q := `
        SELECT id, email, username, password_hash, role, verification_token, created_at 
        FROM users WHERE id = $1`
	err := u.db.QueryRow(ctx, q, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.VerificationToken,
		&user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user %v: %w", id, err)
	}
	return &user, nil
}

func (u *UsersRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	q := `
        SELECT id, email, username, password_hash, role, verification_token, created_at 
        FROM users WHERE email = $1`
	err := u.db.QueryRow(ctx, q, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.VerificationToken,
		&user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user %v: %w", email, err)
	}
	return &user, nil
}

func (u *UsersRepository) FindByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	user := models.User{}
	q := `
        SELECT id, email, username, password_hash, role, verification_token, created_at 
        FROM users WHERE verification_token = $1`
	err := u.db.QueryRow(ctx, q, token).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.VerificationToken,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user with token %v: %w", token, err)
	}
	return &user, nil
}

func (u *UsersRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	q := `SELECT id, email, username, password_hash, role, verification_token, created_at 
        FROM users 
`
	rows, err := u.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.Role,
			&user.VerificationToken,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("db rows fail: %w", err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users, nil

}

func (u *UsersRepository) UpdateUser(ctx context.Context, user *models.User) error {
	q := `UPDATE users SET 
                 email = $1, 
                 username = $2, 
                 password_hash = $3,
                 role = $4,
                 verification_token = $5
             WHERE id = $6`

	_, err := u.db.Exec(
		ctx, q, user.Email, user.Username, user.PasswordHash, user.Role, user.VerificationToken, user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user %s: %w", user.ID, err)
	}
	return nil
}
func (u *UsersRepository) DeleteUser(ctx context.Context, id string) error {
	q := `DELETE FROM users WHERE id = $1`

	result, err := u.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete user %s: %w", id, err)
	}

	if result.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	return nil
}
