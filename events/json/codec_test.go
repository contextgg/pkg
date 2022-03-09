package json

import (
	"context"
	"testing"

	"github.com/contextgg/pkg/types"
)

type MassCustomerPayoutsProcessed struct{}

func Test_It(t *testing.T) {
	var body = `
	{"aggregate_id":"482af783-c743-524d-96c4-e5a0a91b7a7f","aggregate_type":"MassCustomerPayouts","version":3,"type":"MassCustomerPayoutsProcessed","timestamp":"2022-03-04T15:25:46.004196471Z","data":{},"metadata":{},"context":{"namespace":"demo","user":"{\"id\":\"7093b1bd-c7d9-5b1f-8dd0-c436b18965d6\",\"connection\":\"Standard\",\"username\":\"admin\",\"roles\":null,\"registered\":false,\"display_name\":\"Admin\",\"metadata\":null,\"identities\":[{\"user_id\":\"7093b1bd-c7d9-5b1f-8dd0-c436b18965d6\",\"username\":\"\",\"provider\":\"Standard\",\"metadata\":null}],\"audience\":\"demo\"}"}}
	`

	types.Upsert(&MassCustomerPayoutsProcessed{})

	ctx := context.Background()

	c := EventCodec{}
	evt, ctx, err := c.UnmarshalEvent(ctx, []byte(body))
	if err != nil {
		t.Error(err)
		t.Log(evt)
		return
	}
}
