package repositories

import (
	"context"
	"time"

	"github.com/banggibima/be-assignment/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseJobRepository struct {
	db *pgxpool.Pool
}

func NewDatabaseJobRepository(db *pgxpool.Pool) *DatabaseJobRepository {
	return &DatabaseJobRepository{db: db}
}

func (r *DatabaseJobRepository) Create(ctx context.Context, tx pgx.Tx, job *models.Job) error {
	job.ID = uuid.New().String()

	query := "INSERT INTO jobs (id, job_id, status, processed, total, progress, from_date, to_date, created_at, updated_at) "
	query += "VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())"

	if tx != nil {
		_, err := tx.Exec(ctx, query, job.ID, job.JobID, job.Status, job.Processed, job.Total, job.Progress, job.From, job.To)
		return err
	}

	_, err := r.db.Exec(ctx, query, job.ID, job.JobID, job.Status, job.Processed, job.Total, job.Progress, job.From, job.To)
	return err
}

func (r *DatabaseJobRepository) UpdateProgress(ctx context.Context, tx pgx.Tx, jobID string, processed, progress int) error {
	query := "UPDATE jobs "
	query += "SET processed = $1, progress = $2, status = 'RUNNING', updated_at = NOW() "
	query += "WHERE job_id = $3"

	if tx != nil {
		_, err := tx.Exec(ctx, query, processed, progress, jobID)
		return err
	}

	_, err := r.db.Exec(ctx, query, processed, progress, jobID)
	return err
}

func (r *DatabaseJobRepository) MarkDone(ctx context.Context, tx pgx.Tx, jobID, resultPath string) error {
	query := "UPDATE jobs "
	query += "SET status = 'DONE', progress = 100, result_path = $1, updated_at = NOW() "
	query += "WHERE job_id = $2"

	if tx != nil {
		_, err := tx.Exec(ctx, query, resultPath, jobID)
		return err
	}

	_, err := r.db.Exec(ctx, query, resultPath, jobID)
	return err
}

func (r *DatabaseJobRepository) MarkCancelled(ctx context.Context, tx pgx.Tx, jobID string) error {
	query := "UPDATE jobs "
	query += "SET status = 'CANCELLED', updated_at = NOW() "
	query += "WHERE job_id = $1"

	if tx != nil {
		_, err := tx.Exec(ctx, query, jobID)
		return err
	}

	_, err := r.db.Exec(ctx, query, jobID)
	return err
}

func (r *DatabaseJobRepository) GetByID(ctx context.Context, jobID string) (*models.Job, error) {
	query := "SELECT id, job_id, status, processed, total, progress, from_date, to_date, result_path, created_at, updated_at "
	query += "FROM jobs WHERE job_id = $1"

	row := r.db.QueryRow(ctx, query, jobID)

	var j models.Job
	var fromDate, toDate time.Time
	if err := row.Scan(&j.ID, &j.JobID, &j.Status, &j.Processed, &j.Total, &j.Progress, &fromDate, &toDate, &j.ResultPath, &j.CreatedAt, &j.UpdatedAt); err != nil {
		return nil, err
	}

	j.From = fromDate.Format("2006-01-02")
	j.To = toDate.Format("2006-01-02")

	return &j, nil
}
