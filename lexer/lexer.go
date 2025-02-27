package lexer

// Copyright 2021 Simple Lexer Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"fmt"
	"regexp"
	"strings"
)

// Log level
const (
	LexerLogNone = iota
	LexerLogError
	LexerLogWarning
	LexerLogInfo
	LexerLogDebug
)

// LexerLogLevel define log level of lexer
var LexerLogLevel = LexerLogError

// <identifier> == hard value
// <identifier> ~= regex [ function call | ID=value ID=value ]
// <identifier> => function call
// If <indentifier> start by '_' data is skip
const (
	HardValue = iota
	RegexValue
	FunctionCall
)

// SkipToken is value to use to ask this token must be skip
const SkipToken int = -1

// SubPattern is sub stype for SubPatternValue
type SubPattern struct {
	// Name of token
	Name string
	// IDValue is generate by yacc
	IDValue int
	// Value the value to search
	Value string
}

// TokenCallback token callback
type TokenCallback = func(string, TokenEntry) (Token, bool)

// TokenEntry the token entry to search
type TokenEntry struct {
	// Name of token
	Name string
	// TypeOf type of token (harde value, regex...)
	TypeOf int
	// Value the value
	Value string
	// FnCallback callback function. Return token and if found
	FnCallback TokenCallback
	// SubValue is string to sub qualifier
	SubValue []SubPattern
	// IDValue is generate by yacc
	IDValue int
	// Only for regex and for performance
	m *regexp.Regexp
}

var endLineRegex = regexp.MustCompile("(\n|\\r\\n)")

// NewHardValueToken create a token to search
func NewHardValueToken(name string, value string, idValue int) TokenEntry {
	return TokenEntry{
		Name:       name,
		TypeOf:     HardValue,
		Value:      value,
		FnCallback: nil,
		SubValue:   nil,
		IDValue:    idValue,
	}
}

// NewFunctionCallToken create a token to search
func NewFunctionCallToken(name string, fnCallback TokenCallback, idValue int) TokenEntry {
	return TokenEntry{
		Name:       name,
		TypeOf:     FunctionCall,
		Value:      "",
		FnCallback: fnCallback,
		SubValue:   nil,
		IDValue:    idValue,
	}
}

// NewRegexValueToken create a token to search
func NewRegexValueToken(name string, value string, idValue int) TokenEntry {
	return TokenEntry{
		Name:       name,
		TypeOf:     RegexValue,
		Value:      value,
		FnCallback: nil,
		SubValue:   nil,
		IDValue:    idValue,
		m:          regexp.MustCompile("^" + value),
	}
}

// NewRegexWithSubValueToken create a token to search
func NewRegexWithSubValueToken(name string, value string, subValue []SubPattern, idValue int) TokenEntry {
	return TokenEntry{
		Name:       name,
		TypeOf:     RegexValue,
		Value:      value,
		FnCallback: nil,
		SubValue:   subValue,
		IDValue:    idValue,
		m:          regexp.MustCompile("^" + value),
	}
}

// NewRegexWithSubValueFnToken create a token to search
func NewRegexWithSubValueFnToken(name string, value string, fnCallback TokenCallback, idValue int) TokenEntry {
	return TokenEntry{
		Name:       name,
		TypeOf:     RegexValue,
		Value:      value,
		FnCallback: fnCallback,
		SubValue:   nil,
		IDValue:    idValue,
		m:          regexp.MustCompile("^" + value),
	}
}

// FindStringIndex find an str for regex value
func (t *TokenEntry) FindStringIndex(text string) []int {
	// Token must always start at first position, cause each time of
	// NextToken() call, previous data skip.
	return t.m.FindStringIndex(text)
}

// Token token found, first pos, last pos
type Token struct {
	// The name of token
	Name string
	// IDValue to return to Lex() function
	IDValue int
	// LineNumber line number where token found
	LineNumber int
	// StartPos start position in line
	StartPos int
	// EndPos end position in line
	Lenght int
	// Token data cause regex
	Data string
}

// TokenExtraInformation extra data about token
type TokenExtraInformation struct {
	// End position of token in string in char
	EndPosInText int
	// Number of line contain in token
	LineIncludeInToken int
}

// Lexer read text and convert it in Token
func Lexer(text string, tokensList []TokenEntry) ([]Token, error) {
	var currentToken Token
	var isFound bool
	// character position in line
	charPos := 1
	// current line number
	lineNumber := 1
	// character position in text
	charPosInGlobalText := 0
	// the length of text to stop
	lenOfText := len(text)
	// Tokens list
	tokens := []Token{}

	for charPosInGlobalText < lenOfText {
		currentToken, isFound = searchToken(text[charPosInGlobalText:], tokensList)

		if isFound {
			debugLog("Lexer", "Token %+v found", currentToken)

			if currentToken.IDValue == SkipToken {
				infoLog("Lexer", "Skip token")

				debugLog("Lexer", "Length of token: %d - Current position in original text: %d'", currentToken.Lenght, charPosInGlobalText)
			} else {
				debugLog("Lexer", "Add token in list")

				currentToken.LineNumber = lineNumber
				currentToken.StartPos = charPos

				tokens = append(tokens, currentToken)
			}

			tokenStart := charPosInGlobalText
			tokenEnd := charPosInGlobalText + currentToken.Lenght

			// count number of line to have right index of token
			lineNumberInToken, lastLinePos := countLineEnd(text[tokenStart:tokenEnd])

			if lineNumberInToken == 0 {
				// No new lines
				charPos += currentToken.Lenght

				debugLog("Lexer", "New position in line %d", charPos)
			} else {
				lineNumber += lineNumberInToken
				charPos = currentToken.Lenght - lastLinePos[1] + 1 // +1 cause human position start 1

				debugLog("Lexer", "Line number %d, position in line %d", lineNumber, charPos)
			}

			// Increment to end of token to continue search
			charPosInGlobalText += currentToken.Lenght
		} else {
			indexOfChar := charPos - 1
			errorCode := extractPartOfText(text[charPosInGlobalText-indexOfChar:], indexOfChar)

			errorLog("Lexer", "No token found at %d:%d!\n%s", lineNumber, charPos, errorCode)

			return tokens, fmt.Errorf("invalid token found at %d:%d\n%s", lineNumber, charPos, errorCode)
		}
	}

	return tokens, nil
}

func extractPartOfText(text string, start int) string {
	var partOfText string
	var end int

	// FindStringIndex return a array [begin end]
	pos := endLineRegex.FindStringIndex(text)

	if len(pos) == 0 {
		end = len(text) - 1
	} else {
		end = pos[0]
	}

	partOfText = text[0:end]
	partOfText = strings.ReplaceAll(partOfText, "\t", " ")

	return partOfText + "\n" + strings.Repeat("_", start) + "^"
}

// Search a token an return if found.
func searchToken(text string, tokensList []TokenEntry) (Token, bool) {
	currentToken := Token{}
	isFound := false

	for _, token := range tokensList {
		debugLog("searchToken", "Current token %+v", token)

		switch token.TypeOf {
		case HardValue:
			currentToken, isFound = tokenHardValue(text, token)
		case RegexValue:
			currentToken, isFound = tokenRegexValue(text, token)
		default:
			debugLog("searchToken", "Call user search method")
			currentToken, isFound = token.FnCallback(text, token)
		}

		if isFound {
			break
		}
	}

	debugLog("searchToken", "Token return %+v", currentToken)

	return currentToken, isFound
}

// Check if token with hard value found.
func tokenHardValue(text string, token TokenEntry) (Token, bool) {
	lenOfSearch := len(token.Value)

	debugLog("tokenHardValue", "Search hard value '%s'", token.Value)

	if len(text) >= lenOfSearch && text[:lenOfSearch] == token.Value {
		return Token{
			Name:    token.Name,
			IDValue: token.IDValue,
			Lenght:  lenOfSearch,
			Data:    token.Value,
		}, true
	}

	return Token{}, false
}

// Check if token with regex value found.
func tokenRegexValue(text string, token TokenEntry) (Token, bool) {
	debugLog("tokenRegexValue", "Search regex value '%s'", token.Value)

	// FindStringIndex return a array [begin end]
	pos := token.FindStringIndex(text)

	debugLog("tokenRegexValue", "Regex result %+v", pos)

	if len(pos) == 0 {
		return Token{}, false
	}

	value := text[:pos[1]]

	if token.FnCallback != nil {
		debugLog("tokenSubPatternValue", "Call user search method")
		return token.FnCallback(value, token)
	}

	if token.SubValue != nil {
		for _, subValue := range token.SubValue {
			if subValue.Value == value {
				return Token{
					Name:    subValue.Name,
					IDValue: subValue.IDValue,
					Lenght:  pos[1],
					Data:    value,
				}, true
			}
		}
	}

	return Token{
		Name:    token.Name,
		IDValue: token.IDValue,
		Lenght:  pos[1],
		Data:    value,
	}, true
}

func errorLog(methodName, format string, a ...interface{}) {
	if LexerLogLevel >= LexerLogError {
		fmt.Printf("[ERROR] %s(): %s\n", methodName, fmt.Sprintf(format, a...))
	}
}

func infoLog(methodName, format string, a ...interface{}) {
	if LexerLogLevel >= LexerLogInfo {
		fmt.Printf("[INFO] %s(): %s\n", methodName, fmt.Sprintf(format, a...))
	}
}

func debugLog(methodName, format string, a ...interface{}) {
	if LexerLogLevel >= LexerLogDebug {
		fmt.Printf("[DEBUG] %s(): %s\n", methodName, fmt.Sprintf(format, a...))
	}
}

func countLineEnd(str string) (int, []int) {
	pos := endLineRegex.FindAllStringIndex(str, -1)

	debugLog("countLineEnd", "Return line found in token %+v", pos)

	l := len(pos)

	if l == 0 {
		return 0, []int{0, 0}
	}

	// FindAllStringIndex return an array of array [begin end], we take last
	// and return also number of item found that represent number of line found.
	return l, pos[l-1]
}
