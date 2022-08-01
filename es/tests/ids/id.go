package ids

import (
	"github.com/google/uuid"
)

var (
	ledger = uuid.NewSHA1(uuid.NameSpaceURL, []byte("ledger.demo.com"))
)

func generate(id uuid.UUID, parts ...string) uuid.UUID {
	b := id
	for _, p := range parts {
		b = uuid.NewSHA1(b, []byte(p))
	}
	return b
}

func LedgerLineItem(entryId, accountId, subaccountId string) uuid.UUID {
	return generate(ledger, entryId, accountId, subaccountId)
}
