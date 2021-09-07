package t

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TimestampPb timestamppb.Timestamp

func (b *TimestampPb) Value() (driver.Value, error) {
	return ptypes.Timestamp((*timestamppb.Timestamp)(b))
}
func (b *TimestampPb) Scan(value interface{}) error {
	var i sql.NullTime

	if err := i.Scan(value); err != nil {
		return err
	}

	tsp, err := ptypes.TimestampProto(i.Time)
	if err != nil {
		return err
	}

	*b = *(*TimestampPb)(tsp)
	return fmt.Errorf("Error converting type %T into TimestampProto", value)
}
