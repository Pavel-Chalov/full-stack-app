package repository

import (
	"database/sql"
	"fmt"
	"time"
	"trello-backend/lib"
	"trello-backend/models"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(input *models.AuthInput) (*models.User, *lib.WebError) {
	query := fmt.Sprintf("INSERT INTO %s (name, password) VALUES ($1, $2) RETURNING *", UsersTable)

	user := &models.User{}

	row := r.db.QueryRow(query, input.Name, input.Password)

	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Role); err != nil {
		return nil, lib.ServerError(err.Error())
	}

	return user, nil
}

func (r *AuthPostgres) GetUser(name string) (*models.User, *lib.WebError) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", UsersTable)

	row := r.db.QueryRow(query, name)

	user := &models.User{}

	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, lib.NotFound("Пользователь с таким именем не найден")
		}
		return nil, lib.ServerError(err.Error())
	}

	return user, nil
}

func (r *AuthPostgres) ChangeUserData(input *models.AuthInput, id int) *lib.WebError {
	query := fmt.Sprintf("UPDATE %s SET name=$1, password=$2, updated_at=$3 WHERE id=$4 RETURNING *", UsersTable)

	user := &models.User{}

	row := r.db.QueryRow(query, input.Name, input.Password, time.Now(), id)

	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.Role); err != nil {
		fmt.Println(err)
		return lib.ServerError(err.Error())
	}

	return nil
}

func (r *AuthPostgres) GetRefreshSession(refreshToken string) (*models.RefreshSession, *lib.WebError) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE refresh_token=$1", RefreshSessionTable)

	refreshSession := &models.RefreshSession{}

	row := r.db.QueryRow(query, refreshToken)

	if err := row.Scan(&refreshSession.Id, &refreshSession.UserId, &refreshSession.RefreshToken, &refreshSession.FingerPrint); err != nil {
		if err == sql.ErrNoRows {
			return nil, lib.NotFound("Вы не авторизованы")
		}

		return nil, lib.ServerError(err.Error())
	}

	return refreshSession, nil
}

func (r *AuthPostgres) CreateRefreshSession(id int, refreshToken string, fingerPrint string) *lib.WebError {
	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token, finger_print) VALUES ($1, $2, $3) RETURNING *", RefreshSessionTable)

	refreshSession := &models.RefreshSession{}

	row := r.db.QueryRow(query, id, refreshToken, fingerPrint)

	if err := row.Scan(&refreshSession.Id, &refreshSession.UserId, &refreshSession.RefreshToken, &refreshSession.FingerPrint); err != nil {
		if err == sql.ErrNoRows {
			return lib.NotFound("Вы не авторизованы")
		}

		return lib.ServerError(err.Error())
	}

	return nil
}

func (r *AuthPostgres) DeleteRefreshSession(refreshToken string) *lib.WebError {
	if refreshToken == "" {
		return nil
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE refresh_token=$1", RefreshSessionTable)

	if _, err := r.db.Exec(query, refreshToken); err != nil {
		return lib.ServerError(err.Error())
	}

	return nil
}
