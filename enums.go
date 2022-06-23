package jsonwindow

//go:generate stringer -type TokenType,ValueType -output enums_string.go

type TokenType int

const (
	StringToken TokenType = iota + 1
	NumberToken
	CommaToken
	ColonToken
	OpenObjectToken
	CloseObjectToken
	OpenArrayToken
	CloseArrayToken
	TrueToken
	FalseToken
	NullToken
)

type ValueType int

const (
	StringValue ValueType = iota + 1
	NumberValue
	BooleanValue
	NullValue
	ArrayValue
	ObjectValue
)
