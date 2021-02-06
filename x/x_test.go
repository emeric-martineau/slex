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
	"testing"

	"github.com/andreyvit/diff"
)

func Test_Generate(t *testing.T) {
	data := `// Comment
	NUMBER     == 123
	_SPACE	~= (\s)
	// Another comment
_NEWLINE            ~= (\n|\r|\r\n)
				IDENTIFIER ~= ([a-z]+)	MODULE=module \
					   END=end
					   _COMMENT => skipComment	
		DASH ~= (-+)	countDash
`
	// TODO add prefix of NewXXX, SubPattern, TokenEntry
	dataToGet := `[]TokenEntry{
	NewHardValueToken("NUMBER", "123", NUMBER),
	NewRegexValueToken("_SPACE", "(\s)", -1),
	NewRegexValueToken("_NEWLINE", "(\n|\r|\r\n)", -1),
	NewRegexWithSubValueToken("IDENTIFIER", "([a-z]+)", 
		[]SubPattern{
			{"MODULE", MODULE, "module"},
			{"END", END, "end"},
		},
		IDENTIFIER,
	),
	NewFunctionCallToken("_COMMENT", skipComment, -1),
	NewRegexWithSubValueFnToken("DASH", "(-+)", countDash, DASH),
}
`
	dataToWriteInFile, err := ParseParameters(data, "")

	if err != nil {
		t.Errorf(err.Error())
	} else if dataToWriteInFile != dataToGet {
		t.Errorf("Result not as expected:\n%v", diff.LineDiff(dataToWriteInFile, dataToGet))
	}
}

func Test_Generate_With_Packagename(t *testing.T) {
	data := `// Comment
	NUMBER     == 123
	_SPACE	~= (\s)
	// Another comment
_NEWLINE            ~= (\n|\r|\r\n)
				IDENTIFIER ~= ([a-z]+)	MODULE=module \
					   END=end
					   _COMMENT => skipComment	
		DASH ~= (-+)	countDash
`
	// TODO add prefix of NewXXX, SubPattern, TokenEntry
	dataToGet := `[]x.TokenEntry{
	x.NewHardValueToken("NUMBER", "123", NUMBER),
	x.NewRegexValueToken("_SPACE", "(\s)", -1),
	x.NewRegexValueToken("_NEWLINE", "(\n|\r|\r\n)", -1),
	x.NewRegexWithSubValueToken("IDENTIFIER", "([a-z]+)", 
		[]x.SubPattern{
			{"MODULE", MODULE, "module"},
			{"END", END, "end"},
		},
		IDENTIFIER,
	),
	x.NewFunctionCallToken("_COMMENT", skipComment, -1),
	x.NewRegexWithSubValueFnToken("DASH", "(-+)", countDash, DASH),
}
`
	dataToWriteInFile, err := ParseParameters(data, "x")

	if err != nil {
		t.Errorf(err.Error())
	} else if dataToWriteInFile != dataToGet {
		t.Errorf("Result not as expected:\n%v", diff.LineDiff(dataToWriteInFile, dataToGet))
	}
}

func Test_Errors_MultiLine(t *testing.T) {
	_, err := ParseParameters("A=1 \\", "")

	if err.Error() != "Continue line '\\' without newline" {
		t.Error("No error when not found newline")
	}
}

func Test_Error_Missing_Symbol(t *testing.T) {
	_, err := ParseParameters("aaa", "")

	if err.Error() != "Synthaxe error, missing symbol after 'aaa'" {
		t.Error("No error when not found equal")
	}
}

func Test_Error_Missing_Value(t *testing.T) {
	_, err := ParseParameters("aaa=>", "")

	if err.Error() != "Synthaxe error, missing value after 'aaa=>'" {
		t.Error("No error when not found value")
	}
}

func Test_Error_Missing_SubValue(t *testing.T) {
	_, err := ParseParameters("aaa~=(aaa)\tA=", "")

	if err.Error() != "Syntax error. Missing value of sub parameter 'A'" {
		t.Error("No error when not found sub value")
	}
}
