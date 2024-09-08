package repository

import (
	"database/sql"
	"fmt"
	"trello-backend/lib"
	"trello-backend/models"

	"github.com/jmoiron/sqlx"
)

type SettingsPostgres struct {
	db *sqlx.DB
}

func NewSettingsPostgres(db *sqlx.DB) *SettingsPostgres {
	return &SettingsPostgres{db: db}
}

func (r *SettingsPostgres) GetSettings(userId int) (*models.Settings, *lib.WebError) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", SettingsTable)

	settings := &models.Settings{}

	row := r.db.QueryRow(query, userId)

	if err := row.Scan(&settings.Id, &settings.UserId, &settings.BirthDay); err != nil {
		if err == sql.ErrNoRows {
			return nil, lib.NotFound("settings is not found")
		}

		return nil, lib.ServerError(err.Error())
	}

	return settings, nil
}

func (r *SettingsPostgres) CreateSettings(userId int) *lib.WebError {
	query := fmt.Sprintf("INSERT INTO %s (user_id, birth_day) VALUES ($1, $2)", SettingsTable)

	if _, err := r.db.Exec(query, userId, nil); err != nil {
		return lib.ServerError(err.Error())
	}

	return nil
}

func (r *SettingsPostgres) UpdateSettings(settings models.Settings) *lib.WebError {
	query := fmt.Sprintf("UPDATE %s SET birth_day=$1 WHERE id=$2 AND user_id=$3", SettingsTable)

	if _, err := r.db.Exec(query, settings.BirthDay, settings.Id, settings.UserId); err != nil {
		return lib.ServerError(err.Error())
	}

	return nil
}
