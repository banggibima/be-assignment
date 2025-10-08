package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/banggibima/be-assignment/internal/dto"
	"github.com/banggibima/be-assignment/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SettlementRepository interface {
	UpsertJob(ctx context.Context, tx pgx.Tx, runID, merchantID, date string, gross, fee, net, txnCount int) error
}

type JobRepository interface {
	Create(ctx context.Context, tx pgx.Tx, job *models.Job) error
	UpdateProgress(ctx context.Context, tx pgx.Tx, jobID string, processed, progress int) error
	MarkDone(ctx context.Context, tx pgx.Tx, jobID, resultPath string) error
	MarkCancelled(ctx context.Context, tx pgx.Tx, jobID string) error
	GetByID(ctx context.Context, jobID string) (*models.Job, error)
}

type JobService struct {
	db              *pgxpool.Pool
	jobRepo         JobRepository
	transactionRepo TransactionRepository
	settlementRepo  SettlementRepository
	jobQueue        chan string
	cancelSignals   map[string]chan struct{}
	workers         int
	mu              sync.Mutex
}

func NewJobService(
	db *pgxpool.Pool,
	jobRepo JobRepository,
	transactionRepo TransactionRepository,
	settlementRepo SettlementRepository,
) *JobService {
	numCPU := runtime.NumCPU()
	return &JobService{
		db:              db,
		jobRepo:         jobRepo,
		transactionRepo: transactionRepo,
		settlementRepo:  settlementRepo,
		jobQueue:        make(chan string, 10),
		cancelSignals:   make(map[string]chan struct{}),
		workers:         numCPU,
	}
}

func (s *JobService) StartWorkerPool(ctx context.Context) {
	for i := 0; i < s.workers; i++ {
		go s.worker(ctx, i)
	}
	fmt.Printf("[JobWorker] %d workers started\n", s.workers)
}

func (s *JobService) worker(ctx context.Context, workerID int) {
	for jobID := range s.jobQueue {
		s.mu.Lock()
		cancelChan := make(chan struct{})
		s.cancelSignals[jobID] = cancelChan
		s.mu.Unlock()

		fmt.Printf("[Worker-%d] Start processing job %s\n", workerID, jobID)

		job, err := s.jobRepo.GetByID(ctx, jobID)
		if err != nil {
			fmt.Printf("job %s not found: %v\n", jobID, err)
			continue
		}
		if job.From == "" || job.To == "" {
			fmt.Printf("job %s has empty from/to: %q/%q\n", jobID, job.From, job.To)
			s.jobRepo.MarkCancelled(ctx, nil, jobID)
			continue
		}

		folder := "/tmp/settlements"
		if err := os.MkdirAll(folder, 0o755); err != nil {
			fmt.Printf("failed to create folder: %v\n", err)
			s.jobRepo.MarkCancelled(ctx, nil, jobID)
			continue
		}

		path := filepath.Join(folder, job.JobID+".csv")
		file, err := os.Create(path)
		if err != nil {
			fmt.Printf("failed to create CSV: %v\n", err)
			s.jobRepo.MarkCancelled(ctx, nil, jobID)
			continue
		}
		writer := csv.NewWriter(file)
		writer.Write([]string{"merchant_id", "date", "gross", "fee", "net", "txn_count"})

		settlementsMap := make(map[string]*models.Settlement)

		total, _ := s.transactionRepo.CountByDateRange(ctx, job.From, job.To)
		processed := 0
		limit := 5000
		offset := 0

		for {
			select {
			case <-cancelChan:
				fmt.Println("[Worker] Job cancelled")
				s.jobRepo.MarkCancelled(ctx, nil, jobID)
				file.Close()
				return
			default:
			}

			rows, err := s.transactionRepo.FetchBatch(ctx, job.From, job.To, limit, offset)
			if err != nil {
				fmt.Printf("failed to fetch batch: %v\n", err)
				s.jobRepo.MarkCancelled(ctx, nil, jobID)
				file.Close()
				return
			}

			count := 0
			for rows.Next() {
				var txnID, orderID, merchantID, status string
				var amount, fee int
				var paidAt time.Time

				if err := rows.Scan(&txnID, &orderID, &merchantID, &amount, &fee, &status, &paidAt, new(interface{}), new(interface{})); err != nil {
					fmt.Printf("row scan error: %v\n", err)
					continue
				}

				if status != "PAID" {
					continue
				}

				date := paidAt.Format("2006-01-02")
				key := merchantID + "|" + date
				if settlement, ok := settlementsMap[key]; ok {
					settlement.GrossAmount += amount
					settlement.FeeAmount += fee
					settlement.NetAmount += amount - fee
					settlement.TxnCount++
				} else {
					settlementsMap[key] = &models.Settlement{
						GrossAmount: amount,
						FeeAmount:   fee,
						NetAmount:   amount - fee,
						TxnCount:    1,
					}
				}

				processed++
				count++

				progress := (processed * 100) / total
				s.jobRepo.UpdateProgress(ctx, nil, jobID, processed, progress)
			}
			rows.Close()
			if count < limit {
				break
			}
			offset += limit
		}

		for key, settlement := range settlementsMap {
			parts := strings.Split(key, "|")
			record := []string{
				parts[0], parts[1],
				strconv.Itoa(settlement.GrossAmount),
				strconv.Itoa(settlement.FeeAmount),
				strconv.Itoa(settlement.NetAmount),
				strconv.Itoa(settlement.TxnCount),
			}
			writer.Write(record)
		}
		writer.Flush()
		file.Close()

		tx, err := s.db.Begin(ctx)
		if err == nil {
			for key, settlement := range settlementsMap {
				parts := strings.Split(key, "|")
				merchantID := parts[0]
				date := parts[1]
				runID := uuid.New().String()
				_ = s.settlementRepo.UpsertJob(ctx, tx, runID, merchantID, date, settlement.GrossAmount, settlement.FeeAmount, settlement.NetAmount, settlement.TxnCount)
			}
			tx.Commit(ctx)
		}

		s.jobRepo.MarkDone(ctx, nil, jobID, path)
		fmt.Printf("[Worker-%d] Job %s DONE, CSV path: %s\n", workerID, jobID, path)

		s.mu.Lock()
		delete(s.cancelSignals, jobID)
		s.mu.Unlock()
	}
}

func (s *JobService) CreateJob(ctx context.Context, req dto.CreateSettlementJobRequest) (*dto.CreateSettlementJobResponse, error) {
	if req.From == "" || req.To == "" {
		return nil, fmt.Errorf("invalid job request: From and To must be set")
	}

	total, err := s.transactionRepo.CountByDateRange(ctx, req.From, req.To)
	if err != nil {
		fmt.Printf("failed to count total transactions: %v\n", err)
		total = 0
	}

	job := &models.Job{
		ID:        uuid.New().String(),
		JobID:     uuid.New().String(),
		Status:    "QUEUED",
		Processed: 0,
		Total:     total,
		Progress:  0,
		From:      req.From,
		To:        req.To,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.jobRepo.Create(ctx, tx, job); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	select {
	case s.jobQueue <- job.JobID:
		fmt.Println("[JobService] Job queued:", job.JobID)
	default:
		fmt.Println("[JobService] Job queue is full, retry later")
	}

	res := &dto.CreateSettlementJobResponse{
		JobID:   job.JobID,
		Status:  "QUEUED",
		Message: "Job created and queued successfully",
	}

	return res, nil
}

func (s *JobService) CancelJob(jobID string) (*dto.CancelJobResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cancel, ok := s.cancelSignals[jobID]; ok {
		close(cancel)
		delete(s.cancelSignals, jobID)
		s.jobRepo.MarkCancelled(context.Background(), nil, jobID)
		fmt.Println("[JobService] Job cancelled:", jobID)
		return &dto.CancelJobResponse{
			JobID:   jobID,
			Status:  "CANCELLED",
			Message: "Job cancelled successfully",
		}, nil
	}

	return &dto.CancelJobResponse{
		JobID:   jobID,
		Status:  "NOT_FOUND",
		Message: "Job not found or already finished",
	}, nil
}

func (s *JobService) GetJobStatus(ctx context.Context, jobID string) (*dto.JobStatusResponse, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	res := &dto.JobStatusResponse{
		JobID:      job.JobID,
		Status:     job.Status,
		Processed:  job.Processed,
		Total:      job.Total,
		Progress:   job.Progress,
		ResultPath: job.ResultPath,
	}

	if job.Status == "DONE" && job.ResultPath != nil && *job.ResultPath != "" {
		url := fmt.Sprintf("/tmp/settlements/%s.csv", job.JobID)
		res.DownloadURL = &url
	}

	return res, nil
}
