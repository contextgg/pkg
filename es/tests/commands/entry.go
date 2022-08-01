package commands

import (
	"time"

	"github.com/contextgg/pkg/es"
)

// LedgerAddEntryCommand adds a ledger entry
type LedgerAddEntryCommand struct {
	es.BaseCommand

	LedgerId    string
	Book        string
	Description string
	At          time.Time
	Items       []*LedgerLineItem
}

type LedgerLineItem struct {
	AccountId    string
	SubAccountId string
	Credit       int64
	Debit        int64
}
