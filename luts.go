package jsonwindow

var peekTokenLUT = [256]TokenType{
	// 0x00 to 0x07
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	// 0x08 to 0x0f
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	// 0x10 to 0x17
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	// 0x18 to 0x1f
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	// 0x20 to 0x27
	0,
	0,
	StringToken, // "
	0,
	0,
	0,
	0,
	0,
	// 0x28 to 0x2f
	0,
	0,
	0,
	0,
	CommaToken,  // ,
	NumberToken, // -
	0,
	0,
	// 0x30 to 0x37
	NumberToken, // 0
	NumberToken, // 1
	NumberToken, // 2
	NumberToken, // 3
	NumberToken, // 4
	NumberToken, // 5
	NumberToken, // 6
	NumberToken, // 7
	// 0x38 to 0x3f
	NumberToken, // 8
	NumberToken, // 9
	ColonToken,  // :
	0,
	0,
	0,
	0,
	0,
	// 0x40 to 0x47
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	// 0x48 to 0x4f
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	// 0x50 to 0x57
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	0,
	// 0x58 to 0x5f
	0,
	0,
	0,
	OpenArrayToken, // [
	0,
	CloseArrayToken, // ]
	0,
	0,
	// 0x60 to 0x67
	0,
	0,
	0,
	0,
	0,
	0,
	FalseToken, // f
	0,
	// 0x68 to 0x6f
	0,
	0,
	0,
	0,
	0,
	0,
	NullToken, // n
	0,
	// 0x70 to 0x77
	0,
	0,
	0,
	0,
	TrueToken, // t
	0,
	0,
	0,
	// 0x78 to 0x7f
	0,
	0,
	0,
	OpenObjectToken, // {
	0,
	CloseObjectToken, // }
	0,
	0,
}
