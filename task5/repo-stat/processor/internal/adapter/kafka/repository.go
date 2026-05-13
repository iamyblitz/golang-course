package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"repo-stat/platform/kafkamessage"
	"repo-stat/processor/internal/domain"

	"github.com/segmentio/kafka-go"
)

type TaskProducer struct {
	log    *slog.Logger
	writer *kafka.Writer
}

func NewTaskProducer(brokers []string, topic string, log *slog.Logger) *TaskProducer {
	return &TaskProducer{
		log: log,
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

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(owner + "/" + repo),
		Value: payload,
	}); err != nil {
		return err
	}

	p.log.Debug("published collect task", "owner", owner, "repo", repo)
	return nil
}

func (p *TaskProducer) Close() error {
	return p.writer.Close()
}

type ResultConsumer struct {
	log    *slog.Logger
	reader *kafka.Reader
	store  ResultStore
}

type ResultStore interface {
	SaveInfo(ctx context.Context, owner string, repo string, info domain.RepositoryInfo) error
	SaveError(ctx context.Context, owner string, repo string, message string) error
}

func NewResultConsumer(brokers []string, topic string, groupID string, store ResultStore, log *slog.Logger) *ResultConsumer {
	return &ResultConsumer{
		log: log,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		store: store,
	}
}

func (c *ResultConsumer) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("fetch kafka result: %w", err)
		}

		var result kafkamessage.CollectResult
		if err := json.Unmarshal(msg.Value, &result); err != nil {
			c.log.Error("failed to unmarshal collect result", "error", err)
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}

		c.log.Debug("received collect result", "owner", result.Owner, "repo", result.Repo, "error", result.Error)

		if result.Error != "" {
			err = c.store.SaveError(ctx, result.Owner, result.Repo, result.Error)
		} else {
			err = c.store.SaveInfo(ctx, result.Owner, result.Repo, domain.RepositoryInfo{
				FullName:    result.FullName,
				Description: result.Description,
				Stars:       result.Stars,
				Forks:       result.Forks,
				CreatedAt:   result.CreatedAt,
			})
		}
		if err != nil {
			return fmt.Errorf("save collect result: %w", err)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit collect result: %w", err)
		}
	}
}

func (c *ResultConsumer) Close() error {
	return c.reader.Close()
}
