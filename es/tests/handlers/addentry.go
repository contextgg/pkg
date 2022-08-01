package handlers

import (
	"context"

	"github.com/contextgg/pkg/es"
	"github.com/contextgg/pkg/es/tests/aggregates"
	"github.com/contextgg/pkg/es/tests/commands"
	"github.com/contextgg/pkg/es/tests/ids"
)

func createEntry(cmd *commands.LedgerAddEntryCommand, namespace string) *aggregates.Entry {
	n := aggregates.NewEntry(cmd.AggregateId).(*aggregates.Entry)
	n.SetNamespace(namespace)
	n.LedgerId = cmd.LedgerId
	n.Book = cmd.Book
	n.Description = cmd.Description
	n.At = cmd.At
	return n
}
func createLineItems(cmd *commands.LedgerAddEntryCommand, namespace string) []*aggregates.LineItem {
	// build the line items!
	lineItems := make([]*aggregates.LineItem, len(cmd.Items))
	for i, item := range cmd.Items {
		id := ids.LedgerLineItem(cmd.AggregateId, item.AccountId, item.SubAccountId)
		lineItem := aggregates.NewLineItem(id.String()).(*aggregates.LineItem)
		lineItem.SetNamespace(namespace)
		lineItem.EntryId = cmd.AggregateId
		lineItem.LedgerId = item.AccountId
		lineItem.AccountId = item.SubAccountId
		lineItem.Credit = item.Credit
		lineItem.Debit = item.Debit
		lineItems[i] = lineItem
	}
	return lineItems
}

type addEntryHandler struct {
	entryQuery    es.Query[*aggregates.Entry]
	lineItemQuery es.Query[*aggregates.LineItem]
}

func (h *addEntryHandler) HandleCommand(ctx context.Context, command es.Command) error {
	cmd, ok := command.(*commands.LedgerAddEntryCommand)
	if !ok {
		return es.ErrInvalidCommand
	}

	existing, err := h.entryQuery.Load(ctx, cmd.AggregateId)
	if err != nil {
		return err
	}

	entry := createEntry(cmd, existing.Namespace)
	lineItems := createLineItems(cmd, entry.Namespace)

	if err := h.entryQuery.Save(ctx, entry); err != nil {
		return err
	}
	if err := h.lineItemQuery.Save(ctx, lineItems...); err != nil {
		return err
	}
	return nil
}

func NewAddEntryHandler() es.CommandHandler {
	return &addEntryHandler{
		entryQuery:    es.NewQuery[*aggregates.Entry](),
		lineItemQuery: es.NewQuery[*aggregates.LineItem](),
	}
}
