package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"trello-backend/lib"
	"trello-backend/models"

	"github.com/jmoiron/sqlx"
)

type TimeBlockPostgres struct {
	db *sqlx.DB
}

func NewTimeBlockPostgres(db *sqlx.DB) *TimeBlockPostgres {
	return &TimeBlockPostgres{db: db}
}

func (r *TimeBlockPostgres) GetTimeBlocks(id int) ([]models.TimeBlock, *lib.WebError) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", TimeBlocksTable)

	rows, err := r.db.Query(query, id)

	if err != nil {
		return nil, lib.ServerError(err.Error())
	}

	defer rows.Close()

	timeBlocks := []models.TimeBlock{}

	for rows.Next() {
		timeBlock := models.TimeBlock{}
		var durationStr string

		err := rows.Scan(
			&timeBlock.Id,
			&timeBlock.UserId,
			&timeBlock.Name,
			&timeBlock.Color,
			&timeBlock.Order,
			&durationStr,
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		timeBlock.Duration, err = parseDuration(durationStr)
		if err != nil {
			fmt.Println("Error parsing duration:", err)
			continue
		}

		fmt.Println(timeBlock.Duration)

		timeBlocks = append(timeBlocks, timeBlock)
	}

	return timeBlocks, nil
}

func (r *TimeBlockPostgres) CreateTimeBlock(input *models.TimeBlockInput) (int, *lib.WebError) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, name, color, ordering, duration) VALUES ($1, $2, $3, $4, $5) RETURNING id", TimeBlocksTable)

	row := r.db.QueryRow(query, input.UserId, input.Name, input.Color, input.Order, input.Duration.String())

	var id int

	if err := row.Scan(&id); err != nil {
		return 0, lib.ServerError(err.Error())
	}

	return id, nil
}

func (r *TimeBlockPostgres) DeleteTimeBlock(userId, id int) *lib.WebError {
	query := fmt.Sprintf("DELETE from %s WHERE id = $1 AND user_id = $2", TimeBlocksTable)

	_, err := r.db.Exec(query, id, userId)

	if err != nil {
		return lib.ServerError(err.Error())
	}

	return nil
}

func (r *TimeBlockPostgres) UpdateTimeBlock(input *models.TimeBlock) *lib.WebError {
	query := fmt.Sprintf("UPDATE %s SET name=$1, color=$2, duration=$3 WHERE id=$4 AND user_id=$5", TimeBlocksTable)

	fmt.Println(input)

	_, err := r.db.Exec(query, input.Name, input.Color, input.Duration.String(), input.Id, input.UserId)

	if err != nil {
		return lib.ServerError(err.Error())
	}

	return nil
}

type ChangeOrderProps struct {
	TimeBlocks []models.TimeBlock `json:"timeBlocks"`
}

func (r *TimeBlockPostgres) ChangeOrder(input *ChangeOrderProps, userId int) *lib.WebError {
	query := fmt.Sprintf("UPDATE %s SET ordering=$1 WHERE id=$2 AND user_id=%d", TimeBlocksTable, userId)

	for i := 0; i < len(input.TimeBlocks); i++ {
		_, err := r.db.Exec(query, input.TimeBlocks[i].Order, input.TimeBlocks[i].Id)

		if err != nil {
			return lib.ServerError(err.Error())
		}
	}

	return nil
}

func parseDuration(durationStr string) (time.Duration, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	hours, err := strconv.ParseFloat(parts[0], 32)
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.ParseFloat(parts[2], 32)
	if err != nil {
		return 0, err
	}

	totalSeconds := hours*3600 + minutes*60 + seconds
	return time.Duration(totalSeconds) * time.Second, nil
}
