package es

// Command will find its way to an aggregate
type CommandVersion interface {
	GetVersion() int
}

// BaseCommandVersion to make it easier to get the ID
type BaseCommandVersion struct {
	AggregateId string `json:"aggregate_id"`
	Version     int    `json:"int"`
}

// GetAggregateID return the aggregate id
func (c BaseCommandVersion) GetAggregateId() string {
	return c.AggregateId
}

// GetAggregateID return the aggregate id
func (c BaseCommandVersion) GetVersion() int {
	return c.Version
}
