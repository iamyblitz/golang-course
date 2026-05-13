package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"repo-stat/collector/internal/domain"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/kafkamessage"

	"github.com/segmentio/kafka-go"
)

type Worker struct {
	log    *slog.Logger
	reader *kafka.Reader
	writer *kafka.Writer
	repo   RepositoryProvider
}

type RepositoryProvider interface {
	GetInfo(ctx context.Context, owner string, repo string) (domain.RepositoryInfo, error)
}

func NewWorker(brokers []string, taskTopic string, resultTopic string, groupID string, repo RepositoryProvider, log *slog.Logger) *Worker {
	return &Worker{
		log: log,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   taskTopic,
			GroupID: groupID,
		}),
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    resultTopic,
			Balancer: &kafka.Hash{},
		},
		repo: repo,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		msg, err := w.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("fetch collect task: %w", err)
		}

		var task kafkamessage.CollectTask
		if err := json.Unmarshal(msg.Value, &task); err != nil {
			w.log.Error("failed to unmarshal collect task", "error", err)
			_ = w.reader.CommitMessages(ctx, msg)
			continue
		}

		w.log.Debug("received collect task", "owner", task.Owner, "repo", task.Repo)

		result := w.collect(ctx, task)
		if err := w.publishResult(ctx, result); err != nil {
			return err
		}
		if err := w.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit collect task: %w", err)
		}
	}
}

func (w *Worker) Close() error {
	readerErr := w.reader.Close()
	writerErr := w.writer.Close()
	if readerErr != nil {
		return readerErr
	}
	return writerErr
}

func (w *Worker) collect(ctx context.Context, task kafkamessage.CollectTask) kafkamessage.CollectResult {
	info, err := w.repo.GetInfo(ctx, task.Owner, task.Repo)
	if err != nil {
		errorCode := "github_unavailable"
		if errors.Is(err, usecase.ErrRepositoryNotFound) {
			errorCode = "not_found"
		}
		return kafkamessage.CollectResult{
			Owner: task.Owner,
			Repo:  task.Repo,
			Error: errorCode,
		}
	}

	return kafkamessage.CollectResult{
		Owner:       task.Owner,
		Repo:        task.Repo,
		FullName:    info.FullName,
		Description: info.Description,
		Stars:       info.Stars,
		Forks:       info.Forks,
		CreatedAt:   info.CreatedAt,
	}
}

func (w *Worker) publishResult(ctx context.Context, result kafkamessage.CollectResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal collect result: %w", err)
	}

	if err := w.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(result.Owner + "/" + result.Repo),
		Value: payload,
		Time:  time.Now(),
	}); err != nil {
		return err
	}

	w.log.Debug("published collect result", "owner", result.Owner, "repo", result.Repo, "error", result.Error)
	return nil
}

type TaskProducer struct {
	writer *kafka.Writer
}

func NewTaskProducer(brokers []string, topic string) *TaskProducer {
	return &TaskProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
		},
	}
}

func (p *TaskProducer) PublishCollectTask(ctx context.Context, owner string, repo string) error {
	payload, err := json.Marshal(kafkamessage.CollectTask{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return fmt.Errorf("marshal collect task: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(owner + "/" + repo),
		Value: payload,
		Time:  time.Now(),
	})
}

func (p *TaskProducer) Close() error {
	return p.writer.Close()
}
