package jsonwindow

import "fmt"

type ValueType int

const (
	StringValue ValueType = iota + 1
	NumberValue
	BooleanValue
	NullValue
	ArrayValue
	ObjectValue
)

func (v ValueType) String() string {
	switch v {
	case StringValue:
		return "StringValue"
	case NumberValue:
		return "NumberValue"
	case BooleanValue:
		return "BooleanValue"
	case NullValue:
		return "NullValue"
	case ArrayValue:
		return "ArrayValue"
	case ObjectValue:
		return "ObjectValue"
	default:
		return fmt.Sprintf("ValueType(%d)", v)
	}
}
