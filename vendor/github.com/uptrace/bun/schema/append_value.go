package schema

import (
	"database/sql/driver"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/sqltype"
	"github.com/uptrace/bun/extra/bunjson"
	"github.com/uptrace/bun/internal"
)

type (
	AppenderFunc   func(fmter Formatter, b []byte, v reflect.Value) []byte
	CustomAppender func(typ reflect.Type) AppenderFunc
)

var appenders = []AppenderFunc{
	reflect.Bool:          AppendBoolValue,
	reflect.Int:           AppendIntValue,
	reflect.Int8:          AppendIntValue,
	reflect.Int16:         AppendIntValue,
	reflect.Int32:         AppendIntValue,
	reflect.Int64:         AppendIntValue,
	reflect.Uint:          AppendUintValue,
	reflect.Uint8:         AppendUintValue,
	reflect.Uint16:        AppendUintValue,
	reflect.Uint32:        AppendUintValue,
	reflect.Uint64:        AppendUintValue,
	reflect.Uintptr:       nil,
	reflect.Float32:       AppendFloat32Value,
	reflect.Float64:       AppendFloat64Value,
	reflect.Complex64:     nil,
	reflect.Complex128:    nil,
	reflect.Array:         AppendJSONValue,
	reflect.Chan:          nil,
	reflect.Func:          nil,
	reflect.Interface:     nil,
	reflect.Map:           AppendJSONValue,
	reflect.Ptr:           nil,
	reflect.Slice:         AppendJSONValue,
	reflect.String:        AppendStringValue,
	reflect.Struct:        AppendJSONValue,
	reflect.UnsafePointer: nil,
}

func FieldAppender(dialect Dialect, field *Field) AppenderFunc {
	if field.Tag.HasOption("msgpack") {
		return appendMsgpack
	}

	switch strings.ToUpper(field.UserSQLType) {
	case sqltype.JSON, sqltype.JSONB:
		return AppendJSONValue
	}

	return dialect.Appender(field.StructField.Type)
}

func Appender(typ reflect.Type, custom CustomAppender) AppenderFunc {
	switch typ {
	case bytesType:
		return appendBytesValue
	case timeType:
		return appendTimeValue
	case ipType:
		return appendIPValue
	case ipNetType:
		return appendIPNetValue
	case jsonRawMessageType:
		return appendJSONRawMessageValue
	}

	if typ.Implements(queryAppenderType) {
		return appendQueryAppenderValue
	}
	if typ.Implements(driverValuerType) {
		return driverValueAppender(custom)
	}

	kind := typ.Kind()

	if kind != reflect.Ptr {
		ptr := reflect.PtrTo(typ)
		if ptr.Implements(queryAppenderType) {
			return addrAppender(appendQueryAppenderValue, custom)
		}
		if ptr.Implements(driverValuerType) {
			return addrAppender(driverValueAppender(custom), custom)
		}
	}

	switch kind {
	case reflect.Interface:
		return ifaceAppenderFunc(typ, custom)
	case reflect.Ptr:
		if fn := Appender(typ.Elem(), custom); fn != nil {
			return PtrAppender(fn)
		}
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			return appendBytesValue
		}
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return appendArrayBytesValue
		}
	}

	if custom != nil {
		if fn := custom(typ); fn != nil {
			return fn
		}
	}
	return appenders[typ.Kind()]
}

func ifaceAppenderFunc(typ reflect.Type, custom func(reflect.Type) AppenderFunc) AppenderFunc {
	return func(fmter Formatter, b []byte, v reflect.Value) []byte {
		if v.IsNil() {
			return dialect.AppendNull(b)
		}
		elem := v.Elem()
		appender := Appender(elem.Type(), custom)
		return appender(fmter, b, elem)
	}
}

func PtrAppender(fn AppenderFunc) AppenderFunc {
	return func(fmter Formatter, b []byte, v reflect.Value) []byte {
		if v.IsNil() {
			return dialect.AppendNull(b)
		}
		return fn(fmter, b, v.Elem())
	}
}

func AppendBoolValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	return dialect.AppendBool(b, v.Bool())
}

func AppendIntValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	return strconv.AppendInt(b, v.Int(), 10)
}

func AppendUintValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	return strconv.AppendUint(b, v.Uint(), 10)
}

func AppendFloat32Value(fmter Formatter, b []byte, v reflect.Value) []byte {
	return dialect.AppendFloat32(b, float32(v.Float()))
}

func AppendFloat64Value(fmter Formatter, b []byte, v reflect.Value) []byte {
	return dialect.AppendFloat64(b, float64(v.Float()))
}

func appendBytesValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	return dialect.AppendBytes(b, v.Bytes())
}

func appendArrayBytesValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	if v.CanAddr() {
		return dialect.AppendBytes(b, v.Slice(0, v.Len()).Bytes())
	}

	tmp := make([]byte, v.Len())
	reflect.Copy(reflect.ValueOf(tmp), v)
	b = dialect.AppendBytes(b, tmp)
	return b
}

func AppendStringValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	return dialect.AppendString(b, v.String())
}

func AppendJSONValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	bb, err := bunjson.Marshal(v.Interface())
	if err != nil {
		return dialect.AppendError(b, err)
	}

	if len(bb) > 0 && bb[len(bb)-1] == '\n' {
		bb = bb[:len(bb)-1]
	}

	return dialect.AppendJSON(b, bb)
}

func appendTimeValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	tm := v.Interface().(time.Time)
	return fmter.Dialect().AppendTime(b, tm)
}

func appendIPValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	ip := v.Interface().(net.IP)
	return dialect.AppendString(b, ip.String())
}

func appendIPNetValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	ipnet := v.Interface().(net.IPNet)
	return dialect.AppendString(b, ipnet.String())
}

func appendJSONRawMessageValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	bytes := v.Bytes()
	if bytes == nil {
		return dialect.AppendNull(b)
	}
	return dialect.AppendString(b, internal.String(bytes))
}

func appendQueryAppenderValue(fmter Formatter, b []byte, v reflect.Value) []byte {
	return AppendQueryAppender(fmter, b, v.Interface().(QueryAppender))
}

func driverValueAppender(custom CustomAppender) AppenderFunc {
	return func(fmter Formatter, b []byte, v reflect.Value) []byte {
		return appendDriverValue(fmter, b, v.Interface().(driver.Valuer), custom)
	}
}

func appendDriverValue(fmter Formatter, b []byte, v driver.Valuer, custom CustomAppender) []byte {
	value, err := v.Value()
	if err != nil {
		return dialect.AppendError(b, err)
	}
	return Append(fmter, b, value, custom)
}

func addrAppender(fn AppenderFunc, custom CustomAppender) AppenderFunc {
	return func(fmter Formatter, b []byte, v reflect.Value) []byte {
		if !v.CanAddr() {
			err := fmt.Errorf("bun: Append(nonaddressable %T)", v.Interface())
			return dialect.AppendError(b, err)
		}
		return fn(fmter, b, v.Addr())
	}
}
