package auth

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Username string
	password string
	Name     string
}

type LoginAttempt struct {
	Username string
	Password string
}

type SignupAttempt struct {
	Username string
	Password string
	Name     string
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

func (ur *UserRepo) getUserById(ctx context.Context, userId int) (*User, error) {
	return nil, nil
}

func (ur *UserRepo) signIn(ctx context.Context, la LoginAttempt) (*User, error) {
	row := ur.db.QueryRowContext(ctx, `
    SELECT id, username, password, name
    FROM users
    WHERE username = ?
  `, la.Username)
	if row != nil {
		var u User
		err := row.Scan(&u.Id, &u.Username, &u.password, &u.Name)
		if err != nil {
			return nil, err
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.password), []byte(la.Password))
		if err == nil {
			return &u, nil
		}
	}
	return nil, UnauthenticatedError{
		ErrorInfo: "Unauthenticated",
	}
}

func (ur *UserRepo) signUp(ctx context.Context, sa SignupAttempt) error {
	_, err := ur.db.ExecContext(ctx, `
    INSERT INTO users (username, password, name)
    VALUES(
      ?,
      ?,
      ?
    )
    `, sa.Username, sa.Password, sa.Name)
	if err != nil {
		return DatabaseError{Err: err}
	}
	return nil
}

func (ur *UserRepo) usernameInUse(ctx context.Context, username string) (bool, error) {
	row := ur.db.QueryRowContext(ctx, `
      SELECT count(*)
      FROM users
      WHERE username = ?
    `,
		username,
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
