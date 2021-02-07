// Copyright 2021 Emeric MARTINEAU.
%{

package main

import (
	"fmt"
	"strconv"
)

var variable map[string]int

var tokensList []TokenEntry = %%%TOKEN_LIST%%%

var currentToken Token

%}

%start expr

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
	intValue    int
	stringValue string
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <intValue> number addition
%type <stringValue> expr assign print

// same for terminals
%token <stringValue> PRINT IDENTIFIER EQUAL NUMBER ADD

%left ADD  EQUAL

%%

expr	:    assign
	|    print
	;

assign  :    IDENTIFIER EQUAL number
        { variable[$1] = $3 }
	;

print   :    PRINT number
        { fmt.Printf("%+v\n", $2) }
	|        PRINT addition
		{ fmt.Printf("%+v\n", $2) }
	;

number	:    NUMBER
		{
			if s, err := strconv.ParseInt(currentToken.Data, 10, 32); err == nil {
				//fmt.Println("Parse OK")
				$$ = int(s)
			} else {
				//fmt.Println("Parse FAIL")
				//fmt.Printf("%+v\n", currentToken)
				$$ = 0
			} 
		}
	;

addition : number ADD number
		{ $$ = $1 + $3 }
	|		addition ADD number
		{ $$ = $1 + $3 }
	;

%%      /*  start  of  programs  */

type BasicLex struct {
	Tokens []Token
	Index int
}

func (l *BasicLex) Lex(lval *BasicSymType) int {
	if len(l.Tokens) <= l.Index {	
		// Stop
		return 0	
	}

	currentToken = l.Tokens[l.Index]
	l.Index++

	return currentToken.IDValue
}

func (l *BasicLex) Error(s string) {
	fmt.Printf("syntax error: %s\n", s)
}

func main() {
	BasicDebug = 0
	BasicErrorVerbose = true

    tokens, err := Lexer("print 123 + 2 + 3", tokensList)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	lex := BasicLex {
		Tokens: tokens,
		Index: 0,
	}

	BasicParse(&lex)
}

