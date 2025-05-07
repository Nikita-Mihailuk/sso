package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"sso/internal/domain/models"
	"sso/internal/repository"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(storagePath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {

	stmt, err := r.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, email, passwordHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, repository.ErrUserExists
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetUser(ctx context.Context, email string) (models.User, error) {

	stmt, err := r.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email=?")
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	row := stmt.QueryRowContext(ctx, email)

	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repository.ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *Repository) IsAdmin(ctx context.Context, userID int64) (bool, error) {

	stmt, err := r.db.Prepare("SELECT is_admin FROM users WHERE id=?")
	if err != nil {
		return false, err
	}

	var isAdmin bool
	row := stmt.QueryRowContext(ctx, userID)

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.ErrUserNotFound
		}
		return false, err
	}

	return isAdmin, nil
}

func (r *Repository) GetApp(ctx context.Context, appID int) (models.App, error) {

	stmt, err := r.db.Prepare("SELECT id, name, secret FROM apps WHERE id=?")
	if err != nil {
		return models.App{}, err
	}

	var app models.App
	row := stmt.QueryRowContext(ctx, appID)

	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, repository.ErrAppNotFound
		}
		return models.App{}, err
	}

	return app, nil
}
