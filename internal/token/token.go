package token

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

var tokenTypeNames = map[TokenType]string{
	TokenType_LeftParen:    "LeftParen",
	TokenType_RightParen:   "RightParen",
	TokenType_LeftBrace:    "LeftBrace",
	TokenType_RightBrace:   "RightBrace",
	TokenType_Comma:        "Comma",
	TokenType_Dot:          "Dot",
	TokenType_Minus:        "Minus",
	TokenType_Plus:         "Plus",
	TokenType_Semicolon:    "Semicolon",
	TokenType_Slash:        "Slash",
	TokenType_Star:         "Star",
	TokenType_Bang:         "Bang",
	TokenType_BangEqual:    "BangEqual",
	TokenType_Equal:        "Equal",
	TokenType_EqualEqual:   "EqualEqual",
	TokenType_Greater:      "Greater",
	TokenType_GreaterEqual: "GreaterEqual",
	TokenType_Less:         "Less",
	TokenType_LessEqual:    "LessEqual",
	TokenType_Identifier:   "Identifier",
	TokenType_String:       "String",
	TokenType_Number:       "Number",
	TokenType_And:          "And",
	TokenType_Class:        "Class",
	TokenType_Else:         "Else",
	TokenType_False:        "False",
	TokenType_Fun:          "Fun",
	TokenType_For:          "For",
	TokenType_If:           "If",
	TokenType_Nil:          "Nil",
	TokenType_Or:           "Or",
	TokenType_Print:        "Print",
	TokenType_Return:       "Return",
	TokenType_Super:        "Super",
	TokenType_This:         "This",
	TokenType_True:         "True",
	TokenType_Var:          "Var",
	TokenType_While:        "While",
	TokenType_EOF:          "EOF",
}

func (t Token) String() string {
	typeName, _ := tokenTypeNames[t.Type]
	return fmt.Sprintf("%s - %q - %v", typeName, t.Lexeme, t.Literal)
}