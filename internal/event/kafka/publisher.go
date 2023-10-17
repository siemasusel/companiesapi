package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	kafka "github.com/segmentio/kafka-go"
	"github.com/siemasusel/companiesapi/internal/event"
	"golang.org/x/exp/slog"
)

type EventPublisher struct {
	writer *kafka.Writer
}

func NewEventPublisher(writer *kafka.Writer) *EventPublisher {
	return &EventPublisher{
		writer: writer,
	}
}

func (e *EventPublisher) PushEvents(ctx context.Context, events ...event.Event) error {
	msgs, err := marshalEvents(events)
	if err != nil {
		return err
	}

	if err := e.writer.WriteMessages(ctx, msgs...); err != nil {
		return errors.Wrapf(err, "unable to send %d messages to kafka", len(msgs))
	}

	slog.InfoContext(ctx, "Sucessfully sent messages", "count", len(msgs))

	return nil
}

func marshalEvents(events []event.Event) ([]kafka.Message, error) {
	msgs := make([]kafka.Message, 0, len(events))
	for _, e := range events {
		value, err := json.Marshal(e)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to marshall event '%s' for aggregate id '%s'", e.Type, e.AggregateID)
		}

		msgs = append(msgs, kafka.Message{
			Key:   []byte(fmt.Sprintf("%s:%s", e.AggregateName, e.AggregateID)),
			Value: value,
		})
	}

	return msgs, nil
}
