package postgres

import (
	"SkillsRock/internal/domain/models"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v4"
)

type TaskDB struct {
	db *pgx.Conn
}

func NewTaskDb(db *pgx.Conn) *TaskDB {
	return &TaskDB{db: db}
}
func (t *TaskDB) CreateTask(ctx context.Context, task *models.Task) (int64, error) {
	query := `
		INSERT INTO tasks (title,description,status,created_at,updated_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`
	now := time.Now()
	task.Created_at = now
	task.Update_at = now
	if task.Status == "" {
		task.Status = "new"
	}

	var id int64
	err := t.db.QueryRow(ctx, query,
		task.Title, task.Description, task.Status, task.Created_at, task.Update_at).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (t *TaskDB) GetTasks(ctx context.Context) ([]models.Task, int64, error) {
	query := `
		SELECT * FROM tasks
	`
	rows, err := t.db.Query(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Created_at,
			&task.Update_at,
		)
		if err != nil {
			log.Errorf("Error scanning song row: %v", err)
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		log.Errorf("Error iterating song rows: %v", err)
		return nil, 0, err
	}
	return tasks, int64(len(tasks)), nil
}

func (t *TaskDB) DeleteTaskId(ctx context.Context, id int64) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`
	res, err := t.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		log.Debug("No task with ID:", id)
		return fmt.Errorf("task with ID %d not found", id)
	}
	log.Info("Delete task with ID:%d", id)
	return nil
}

func (t *TaskDB) UpdateTaskID(ctx context.Context, id int64, task *models.Task) error {
	count := 1
	stringsSet := "SET "
	arg := make([]interface{}, 0)
	fmt.Println(task)
	if task.Title != "" {
		stringsSet += fmt.Sprintf("title = $%v, ", count)
		arg = append(arg, task.Title)
		count++
	}
	if task.Description != "" {
		stringsSet += fmt.Sprintf("description = $%v, ", count)
		arg = append(arg, task.Description)
		count++
	}
	if task.Status != "" {
		stringsSet += fmt.Sprintf("status =$%v, ", count)
		arg = append(arg, task.Status)
		count++
	}
	now := time.Now()
	task.Update_at = now
	arg = append(arg, task.Update_at, id)
	updated := fmt.Sprintf("updated_at = $%v ", count)
	count++
	query := fmt.Sprintf(`UPDATE tasks 
	 %s %s 
	 WHERE id = $%v`, stringsSet, updated, count)

	res, err := t.db.Exec(ctx, query, arg...)
	if err != nil {
		log.Errorf("Error updating task: %v", err)
		return err
	}
	if res.RowsAffected() == 0 {
		log.Debug("No task with ID:", task.Id)
		return fmt.Errorf("task with ID %d not found", task.Id)
	}
	return nil
}
