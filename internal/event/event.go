package event

import "github.com/google/uuid"

type Event struct {
	AggregateID   uuid.UUID     `json:"aggregate_id"`
	AggregateName AggregateName `json:"aggregate_name"`
	Type          EventType     `json:"type"`
	Payload       interface{}   `json:"payload"`
}

type AggregateName string

const (
	AggregateNameCompany = "company"
)

type EventType string

const (
	EventTypeCompanyCreated = "company.created"
	EventTypeCompanyUpdated = "company.updated"
	EventTypeCompanyDeleted = "company.deleted"
)
