//nolint:all
package audite

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
)

const (
	batchSize   = 5
	timelim     = 5000
	countWorker = 4
)

// httpWorker воркер для хедлеров
func (a *Audit) httpWorker(ctx context.Context) {
	ticker := time.NewTicker(timelim * time.Millisecond)
	defer ticker.Stop()
	batch := make([]interface{}, 0, batchSize)

	for {
		select {
		case <-ctx.Done():

			return
		case helper, ok := <-a.HTTP:
			if !ok {
				if len(batch) > 0 {
					a.ProcessBatch(batch)
				}
				return
			}
			if len(batch) < batchSize {
				batch = append(batch, helper)
			}
			if len(batch) == batchSize {
				a.ProcessBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) != 0 {
				a.ProcessBatch(batch)
				batch = batch[:0]
			}
			ticker = time.NewTicker(timelim * time.Millisecond)
		}
	}
}

// StatusWorker воркер для статусов
func (a *Audit) StatusWorker(ctx context.Context) {
	ticker := time.NewTicker(timelim * time.Millisecond)
	defer ticker.Stop()
	batch := make([]interface{}, 0, batchSize)

	for {
		select {
		case <-ctx.Done():

			return
		case helper, ok := <-a.Status:
			if !ok {
				if len(batch) > 0 {
					a.ProcessBatch(batch)
				}
				return
			}
			if len(batch) < batchSize {
				batch = append(batch, helper)
			}
			if len(batch) == batchSize {
				a.ProcessBatch(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) != 0 {
				a.ProcessBatch(batch)
				batch = batch[:0]
			}
			ticker = time.NewTicker(timelim * time.Millisecond)
		}
	}
}

// StdoutWorker воркер для записи в stdout
func (a *Audit) StdoutWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():

			return
		case val, ok := <-a.Stdout:
			if !ok {
				return
			}
			a.StdoutLog(val)
		}
	}
}

// DBWorker воркер для записи в бд
func (a *Audit) DBWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():

			return
		case val, ok := <-a.DB:
			if !ok {
				return
			}
			if err := a.TX.RunSerializable(ctx, func(ctxTx context.Context) error {
				switch v := val.(type) {
				case entity.HTTPInfo:
					if err := a.Repo.AddHandlerLog(ctxTx, v); err != nil {
						return fmt.Errorf("a.Repo.AddHandlerLog: %w", err)
					}
				case entity.UpdateStatus:
					if err := a.Repo.AddStatusLog(ctxTx, v); err != nil {
						return fmt.Errorf("a.Repo.AddStatusLog: %w", err)
					}
				}

				jsonLog, err := json.Marshal(val)
				if err != nil {
					return fmt.Errorf("json.Marshal: %w", err)
				}
				statusID, err := a.Repo.GetStatusID(ctxTx, consts.CreatedTask)
				if err != nil {
					return fmt.Errorf("a.Repo.GetStatusID: %w", err)
				}
				if err := a.Repo.AddTasks(ctxTx, jsonLog, statusID); err != nil {
					return fmt.Errorf("a.Repo.AddTasks: %w", err)
				}
				return nil
			}); err != nil {
				log.Printf("a.TX.RunSerializable: %v", err)
			}

		}
	}
}

// Run запускает воркеры
func (a *Audit) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(countWorker)
	go func() {
		defer wg.Done()
		a.httpWorker(ctx)
	}()
	go func() {
		defer wg.Done()
		a.StatusWorker(ctx)
	}()
	go func() {
		defer wg.Done()
		a.StdoutWorker(ctx)
	}()
	go func() {
		defer wg.Done()
		a.DBWorker(ctx)
	}()
	wg.Wait()
}
