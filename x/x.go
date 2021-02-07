package x

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
	"slex/lexer"
	"strings"
)

const (
	hardValueToken = iota
	regexValueToken
	functionCallToken
	dataToken
	idToken
	valueToken
)

// ParseParameters convert parameter in file into parameter code
func ParseParameters(data string, packageName string) (string, error) {
	filterTokens, errFilter := filterComment(data)

	if errFilter != nil {
		return "", errFilter
	}

	tokens, errMerge := mergeContinueLine(filterTokens)

	if errMerge != nil {
		return "", errMerge
	}

	if len(packageName) > 0 {
		packageName = packageName + "."
	}

	result := []string{fmt.Sprintf("[]%sTokenEntry{", packageName)}

	for _, token := range tokens {
		lineTokens, errLine := parseOneLine(token)

		if errLine != nil {
			return "", errLine
		}

		line, errGenerate := generateOneLine(lineTokens, packageName)

		if errGenerate != nil {
			return "", errGenerate
		}

		result = append(result, line)
	}

	result = append(result, "}", "") // Empty string to have return line at end

	return strings.Join(result, "\n"), nil
}

// Remove all comments
func filterComment(data string) ([]lexer.Token, error) {
	tokensList := []lexer.TokenEntry{
		lexer.NewRegexValueToken("_COMMENT", "(\\s*//.*)", -1),
		lexer.NewRegexValueToken("_NEWLINE", "(\\r|\\n|(\\r\\n))", -1),
		lexer.NewRegexValueToken("DATA", "([^\\r\\n]+)", dataToken),
	}

	return lexer.Lexer(data, tokensList)
}

// A line can be in multi line with a '\' at end.
// We merge a multi line in one line to be more easier to manage.
func mergeContinueLine(tokens []lexer.Token) ([]lexer.Token, error) {
	// Merge line
	var currentToken lexer.Token
	tokenCount := len(tokens)

	for index := 0; index < tokenCount; {
		currentToken = tokens[index]

		if strings.HasSuffix(currentToken.Data, "\\") {
			// Replace last char by space
			currentToken.Data = currentToken.Data[:len(currentToken.Data)-2] + " "
			// Merge nextline

			// Check if last item to avoid error
			if tokenCount <= index+1 {
				return nil, fmt.Errorf("Continue line '\\' without newline")
			}

			currentToken.Data = currentToken.Data + tokens[index+1].Data
			// Remove item
			tokens = remove(tokens, index+1)
			tokenCount--
			tokens = update(tokens, currentToken, index)
		} else {
			index++
		}
	}

	return tokens, nil
}

// We parse one line of identifiant, type, value....
func parseOneLine(token lexer.Token) ([]lexer.Token, error) {
	tokensList := []lexer.TokenEntry{
		lexer.NewRegexValueToken("_SPACE", "(\\s)", -1),
		lexer.NewRegexValueToken("IDENTIFIANT", "([a-zA-Z_0-9]+)", idToken),
		lexer.NewHardValueToken("HARD_VALUE", "==", hardValueToken),
		lexer.NewHardValueToken("REGEX_VALUE", "~=", regexValueToken),
		lexer.NewHardValueToken("FN_CALL", "=>", functionCallToken),
		lexer.NewRegexValueToken("DATA", "([^\\r\\n]+)", dataToken),
	}

	return lexer.Lexer(token.Data, tokensList)
}

// Remove item in array.
func remove(slice []lexer.Token, index int) []lexer.Token {
	return append(slice[:index], slice[index+1:]...)
}

// Update an array.
func update(slice []lexer.Token, t lexer.Token, index int) []lexer.Token {
	newToken := append(slice[:index], t)
	return append(newToken, slice[index+1:]...)
}

// Generate the string of one line in input file to output file.
func generateOneLine(tokens []lexer.Token, packageName string) (string, error) {
	if len(tokens) < 2 {
		return "", fmt.Errorf("Synthaxe error, missing symbol after '%s'", tokens[0].Data)
	} else if len(tokens) < 3 {
		return "", fmt.Errorf("Synthaxe error, missing value after '%s%s'", tokens[0].Data, tokens[1].Data)
	}

	id := tokens[0].Data
	typeOf := tokens[1].Data
	var num string
	var value string
	var fn string
	var extra string

	// If token id start by underscore, we add special value to ignore it
	if strings.HasPrefix(id, "_") {
		num = "-1"
	} else {
		num = id
	}

	switch typeOf {
	case "=>":
		fn = fmt.Sprintf("%sNewFunctionCallToken", packageName)
		value = tokens[2].Data
		extra = ""
	case "==":
		fn = fmt.Sprintf("%sNewHardValueToken", packageName)
		// TODO escape \
		value = fmt.Sprintf("\"%s\"", tokens[2].Data)
		extra = ""
	case "~=":
		datas := strings.SplitN(tokens[2].Data, "\t", 2)
		value = fmt.Sprintf("\"%s\"", datas[0])

		if len(datas) == 1 {
			fn = fmt.Sprintf("%sNewRegexValueToken", packageName)
			// TODO escape \
			extra = ""
		} else if strings.Index(datas[1], "=") == -1 {
			// Check if = sign is found.
			// If not found, this is a line with function to call
			fn = fmt.Sprintf("%sNewRegexWithSubValueFnToken", packageName)
			extra = fmt.Sprintf("%s, ", datas[1])
		} else {
			fn = fmt.Sprintf("%sNewRegexWithSubValueToken", packageName)

			extras := []string{"", fmt.Sprintf("\t\t[]%sSubPattern{", packageName)}

			subParamsTokens, err := parseSubParameters(datas[1])

			if err != nil {
				return "", fmt.Errorf("Error when parse sub parameters")
			}

			subParameters, errGenerate := generateSubParameters(subParamsTokens)

			if errGenerate != nil {
				return "", errGenerate
			}

			extras = append(extras, subParameters...)
			extras = append(extras, "\t\t},", "\t\t")

			extra = strings.Join(extras, "\n")

			num += ",\n\t"
		}

	default:
		fn = "???"
		value = "???"
	}

	return fmt.Sprintf(
		"\t%s(\"%s\", %s, %s%s),",
		fn, id, value, extra, num), nil
}

// A line with sub parameter A=x B=y ... to be convert into token.
func parseSubParameters(data string) ([]lexer.Token, error) {
	tokensList := []lexer.TokenEntry{
		lexer.NewRegexValueToken("_SPACE", "(\\s)", -1),
		lexer.NewRegexValueToken("IDENTIFIANT", "([a-zA-Z_0-9]+)", idToken),
		lexer.NewHardValueToken("_EQUAL", "=", -1),
		lexer.NewRegexValueToken("VALUE", "([^\\s]+)", valueToken),
	}

	return lexer.Lexer(data, tokensList)
}

// Take a list of token that represent a sub parameter
func generateSubParameters(tokens []lexer.Token) ([]string, error) {
	// A value is always 3 item
	lenOfTokens := len(tokens)
	var id string
	var value string
	result := []string{}

	for index := 0; index < lenOfTokens; index += 2 {
		id = tokens[index].Data

		if lenOfTokens <= index+1 {
			return nil, fmt.Errorf("Syntax error. Missing value of sub parameter '%s'", id)
		}

		value = tokens[index+1].Data // TODO check if A= without value

		result = append(result, fmt.Sprintf("\t\t\t{\"%s\", %s, \"%s\"},", id, id, value))
	}

	return result, nil
}
