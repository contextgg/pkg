package aggregates

import (
	"time"

	"github.com/contextgg/pkg/es"
)

// Entry entry stuff
type Entry struct {
	es.BaseAggregateHolder

	LedgerId    string    `json:"ledger_id"`
	Book        string    `json:"book"`
	Description string    `json:"description"`
	At          time.Time `json:"at"`
}

func NewEntry(id string) es.Entity {
	return &Entry{
		BaseAggregateHolder: es.NewBaseAggregateHolder(id, "Entry"),
	}
}
