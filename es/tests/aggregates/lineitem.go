package aggregates

import "github.com/contextgg/pkg/es"

// LineItem entry stuff
type LineItem struct {
	es.BaseAggregateHolder

	EntryId   string `json:"entry_id"`
	LedgerId  string `json:"ledger_id"`
	AccountId string `json:"account_id"`
	Credit    int64  `json:"credit"`
	Debit     int64  `json:"debit"`
}

func NewLineItem(id string) es.Entity {
	return &LineItem{
		BaseAggregateHolder: es.NewBaseAggregateHolder(id, "LineItem"),
	}
}
