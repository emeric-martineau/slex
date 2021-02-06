package lexer

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	LexerLogLevel = LexerLogNone
	code := m.Run()

	os.Exit(code)
}

func Test_Lexer_One_Token_HardValue(t *testing.T) {
	tokensList := []TokenEntry{
		NewHardValueToken("MODULE", "module", 1),
	}

	tokens, err := Lexer("module", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 1 {
		t.Errorf("Only one token normaly return. It return %d tokens", len(tokens))
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 1 || tk.StartPos != 1 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:1 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}
}

func Test_Lexer_One_Token_RegexValue(t *testing.T) {
	tokensList := []TokenEntry{
		NewRegexValueToken("MODULE", "(m[a-z]+)", 1),
	}

	tokens, err := Lexer("module", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 1 {
		t.Errorf("Only one token normaly return. It return %d tokens", len(tokens))
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 1 || tk.StartPos != 1 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:1 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}
}

func Test_Lexer_One_Token_FunctionCall(t *testing.T) {
	tokensList := []TokenEntry{
		NewFunctionCallToken("COMMENT", skipComment, 1),
	}

	tokens, err := Lexer("/* comment */", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 1 {
		t.Errorf("Only one token normaly return. It return %d tokens", len(tokens))
	}

	tk := tokens[0]

	if tk.Name != "COMMENT" || tk.IDValue != 1 || tk.LineNumber != 1 || tk.StartPos != 1 || tk.Lenght != 13 || tk.Data != "/* comment */" {
		t.Errorf("Expected {Name:COMMENT IDValue:1 LineNumber:1 StartPos:1 Lenght:13 Data:/* comment */} found %+v ", tk)
	}
}

func Test_Lexer_One_Token_Skip(t *testing.T) {
	tokensList := []TokenEntry{
		NewFunctionCallToken("_COMMENT", skipComment, -1),
		NewHardValueToken("MODULE", "module", 1),
	}

	tokens, err := Lexer("/* comment */module", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 1 {
		t.Errorf("Only one token normaly return. It return %d tokens", len(tokens))
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 1 || tk.StartPos != 14 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:1 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}
}

func Test_Lexer_One_Token_Multiline_In_Skip_Token(t *testing.T) {
	tokensList := []TokenEntry{
		NewFunctionCallToken("_COMMENT", skipComment, -1),
		NewHardValueToken("MODULE", "module", 1),
	}

	tokens, err := Lexer("/* comment \n */module", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 1 {
		t.Errorf("Only one token normaly return. It return %d tokens", len(tokens))
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 2 || tk.StartPos != 4 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:2 StartPos:4 Lenght:6 Data:module} found %+v ", tk)
	}
}

func Test_Lexer_One_Token_Multiline_After_Skip_Token(t *testing.T) {
	tokensList := []TokenEntry{
		NewFunctionCallToken("_COMMENT", skipComment, -1),
		NewHardValueToken("MODULE", "module", 1),
		NewRegexValueToken("_NEWLINE", "(\\r|\\n|\\r\\n)", -1),
	}

	tokens, err := Lexer("/* comment */\nmodule", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 1 {
		t.Errorf("Only one token normaly return. It return %d tokens", len(tokens))
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 2 || tk.StartPos != 1 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:2 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}
}

func Test_Lexer_Two_Token_Multiline(t *testing.T) {
	tokensList := []TokenEntry{
		NewFunctionCallToken("_COMMENT", skipComment, -1),
		NewHardValueToken("MODULE", "module", 1),
		NewRegexValueToken("_NEWLINE", "(\\r|\\n|\\r\\n)", -1),
		NewRegexValueToken("_SPACE", "\\s", -1),
	}

	tokens, err := Lexer("/* comment\n */\nmodule\n  \n      module", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 2 {
		t.Errorf("Only two tokens normaly return. It return %d tokens", len(tokens))
		return
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 3 || tk.StartPos != 1 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:2 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}

	tk = tokens[1]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 5 || tk.StartPos != 7 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:5 StartPos:7 Lenght:6 Data:module} found %+v ", tk)
	}
}

func Test_Lexer_Two_Token_Multiline_With_Another_Order(t *testing.T) {
	tokensList := []TokenEntry{
		NewFunctionCallToken("_COMMENT", skipComment, -1),
		NewRegexValueToken("_NEWLINE", "(\\r|\\n|\\r\\n)", -1),
		NewRegexValueToken("_SPACE", "\\s", -1),
		NewHardValueToken("MODULE", "module", 1),
	}

	tokens, err := Lexer("/* comment\n */\nmodule\n  \n      module", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 2 {
		t.Errorf("Only two tokens normaly return. It return %d tokens", len(tokens))
		return
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 3 || tk.StartPos != 1 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:2 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}

	tk = tokens[1]

	if tk.Name != "MODULE" || tk.IDValue != 1 || tk.LineNumber != 5 || tk.StartPos != 7 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:1 LineNumber:5 StartPos:7 Lenght:6 Data:module} found %+v ", tk)
	}
}

func Test_Lexer_Token_Not_Found(t *testing.T) {
	tokensList := []TokenEntry{
		NewHardValueToken("MODULE", "module", 1),
	}

	tokens, err := Lexer("tttt", tokensList)

	if err == nil {
		t.Errorf("A token was found! %+v", tokens)
	}

	if err.Error() != "Invalid token found at 1:1\nttt\n^" {
		t.Errorf("Wrong error message:'%s'", err.Error())
	}
}

func Test_Lexer_Token_With_SubPattern_String(t *testing.T) {
	tokensList := []TokenEntry{
		NewRegexValueToken("_NEWLINE", "(\\r|\\n|\\r\\n)", -1),
		NewRegexWithSubValueToken(
			"IDENTIFIER",
			"([a-zA-Z]+)",
			[]SubPattern{
				{"MODULE", 2, "module"},
				{"TRUC", 3, "truc"},
			},
			0,
		),
	}

	tokens, err := Lexer("module\ntruc\n", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 2 {
		t.Errorf("Only two tokens normaly return. It return %d tokens", len(tokens))
		return
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 2 || tk.LineNumber != 1 || tk.StartPos != 1 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:2 LineNumber:1 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}

	tk = tokens[1]

	if tk.Name != "TRUC" || tk.IDValue != 3 || tk.LineNumber != 2 || tk.StartPos != 1 || tk.Lenght != 4 || tk.Data != "truc" {
		t.Errorf("Expected {Name:MODULE IDValue:3 LineNumber:2 StartPos:1 Lenght:4 Data:truc} found %+v ", tk)
	}
}

func Test_Lexer_Token_With_SubPattern_Callback(t *testing.T) {
	tokensList := []TokenEntry{
		NewRegexValueToken("_NEWLINE", "(\\r|\\n|\\r\\n)", -1),
		NewRegexWithSubValueFnToken(
			"IDENTIFIER",
			"([a-zA-Z]+)",
			searchSubValue,
			0,
		),
	}

	tokens, err := Lexer("module\ntruc\n", tokensList)

	if err != nil {
		t.Errorf("An error occure %+v", err)
	}

	if len(tokens) != 2 {
		t.Errorf("Only two tokens normaly return. It return %d tokens", len(tokens))
		return
	}

	tk := tokens[0]

	if tk.Name != "MODULE" || tk.IDValue != 2 || tk.LineNumber != 1 || tk.StartPos != 1 || tk.Lenght != 6 || tk.Data != "module" {
		t.Errorf("Expected {Name:MODULE IDValue:2 LineNumber:1 StartPos:1 Lenght:6 Data:module} found %+v ", tk)
	}

	tk = tokens[1]

	if tk.Name != "TRUC" || tk.IDValue != 3 || tk.LineNumber != 2 || tk.StartPos != 1 || tk.Lenght != 4 || tk.Data != "truc" {
		t.Errorf("Expected {Name:MODULE IDValue:3 LineNumber:2 StartPos:1 Lenght:4 Data:truc} found %+v ", tk)
	}
}

func skipComment(text string, token TokenEntry) (Token, bool) {
	lenText := len(text)

	if lenText >= 4 && text[0] == '/' && text[1] == '*' {
		for i := 2; i < lenText; i++ {
			if text[i] == '*' && text[i+1] == '/' {
				i += 2

				return Token{
					Name:    token.Name,
					IDValue: token.IDValue,
					Lenght:  i,
					Data:    text[:i],
				}, true
			}
		}
	}

	return Token{}, false
}

func searchSubValue(text string, token TokenEntry) (Token, bool) {
	switch text {
	case "module":
		return Token{
			Name:    "MODULE",
			IDValue: 2,
			Data:    text,
			Lenght:  len(text),
		}, true
	case "truc":
		return Token{
			Name:    "TRUC",
			IDValue: 3,
			Data:    text,
			Lenght:  len(text),
		}, true
	}

	return Token{
		Name:    token.Name,
		IDValue: token.IDValue,
		Data:    text,
		Lenght:  len(text),
	}, true
}
