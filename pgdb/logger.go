package pgdb

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	out, err := q.FormattedQuery()
	if err != nil {
		fmt.Println(string(out))
	}
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	out, err := q.FormattedQuery()
	if err != nil {
		fmt.Println(string(out))
	}
	return nil
}
