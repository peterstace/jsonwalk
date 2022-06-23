// Code generated by "stringer -type TokenType,ValueType -output enums_string.go"; DO NOT EDIT.

package jsonwindow

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StringToken-1]
	_ = x[NumberToken-2]
	_ = x[CommaToken-3]
	_ = x[ColonToken-4]
	_ = x[OpenObjectToken-5]
	_ = x[CloseObjectToken-6]
	_ = x[OpenArrayToken-7]
	_ = x[CloseArrayToken-8]
	_ = x[TrueToken-9]
	_ = x[FalseToken-10]
	_ = x[NullToken-11]
}

const _TokenType_name = "StringTokenNumberTokenCommaTokenColonTokenOpenObjectTokenCloseObjectTokenOpenArrayTokenCloseArrayTokenTrueTokenFalseTokenNullToken"

var _TokenType_index = [...]uint8{0, 11, 22, 32, 42, 57, 73, 87, 102, 111, 121, 130}

func (i TokenType) String() string {
	i -= 1
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StringValue-1]
	_ = x[NumberValue-2]
	_ = x[BooleanValue-3]
	_ = x[NullValue-4]
	_ = x[ArrayValue-5]
	_ = x[ObjectValue-6]
}

const _ValueType_name = "StringValueNumberValueBooleanValueNullValueArrayValueObjectValue"

var _ValueType_index = [...]uint8{0, 11, 22, 34, 43, 53, 64}

func (i ValueType) String() string {
	i -= 1
	if i < 0 || i >= ValueType(len(_ValueType_index)-1) {
		return "ValueType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ValueType_name[_ValueType_index[i]:_ValueType_index[i+1]]
}
