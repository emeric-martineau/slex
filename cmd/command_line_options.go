package cmd

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
	"io/ioutil"
	"os"
	"strings"

	x "slex/x"

	cli "github.com/urfave/cli/v2"
)

// GetCommandLineOptions return cli.App struct
func GetCommandLineOptions() cli.App {
	outputFilename := ""
	inputFilename := ""
	packageName := ""

	return cli.App{
		Name:    "Simple Lexer for goyacc",
		Usage:   "Simple lexer to generate file and code to use with goyacc",
		Version: "v1.0.0",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate go file from lexer file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "output filename",
						Destination: &outputFilename,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "input",
						Aliases:     []string{"i"},
						Usage:       "input filename",
						Destination: &inputFilename,
						Required:    true,
					},
					&cli.StringFlag{
						Name:        "package",
						Aliases:     []string{"p"},
						Usage:       "go package name to set",
						Destination: &packageName,
					},
				},
				Action: func(c *cli.Context) error {
					if packageName == "" {
						packageName = "main"
					}

					data := strings.Replace(slexTemplate, "package lexer", "package "+packageName, 1)

					f, errCreate := os.Create(outputFilename)
					defer f.Close()

					if errCreate != nil {
						return errCreate
					}

					_, errWrite := f.WriteString(data)

					if errWrite != nil {
						return errWrite
					}

					content, errInputfile := ioutil.ReadFile(inputFilename)

					if errInputfile != nil {
						return errInputfile
					}

					// Generate variable
					dataToWriteInFile, errParse := x.ParseParameters(string(content), packageName)

					if errParse != nil {
						return errParse
					}

					fmt.Printf("%s", dataToWriteInFile)

					return nil
				},
			},
			{
				Name:  "example",
				Usage: "Generate example file",
				Action: func(c *cli.Context) error {
					fmt.Println(slexExample)

					return nil
				},
			},
		},
	}
}

var slexExample = `// This file is an example generate by Simple Lexer cli.
// First field is identifier to link with %token (terminal) in goyacc. If indentifier start by '_' data is skip.
// Item are parse in order, becarefull.
//
// After type to lexer:
// == hard value
// ~= regex 
// => call a go function to find token
// 
// Regex
// -----
// Regex has extra parameters.
// By default, Regex return found value with identifier, but you can return another identifier.
// You can associate a go function (see below) to set sub-type or list of identifier with hard value to set sub type.
//
// Call a go function to find token
// --------------------------------
// The function must be type TokenCallback (see lexer.go file).
// Function return a token (or empty token) and if token found.
// If token not found, lexer stop.
// 
// Abstract
// --------
// <identifier> == hard value
// <identifier> ~= regex [ function call | ID=value ID=value ]
//                 regex and function call or value must be separate by tab
// <identifier> => function call
// If <indentifier> start by '_' data is skip
//
NUMBER     == 123
_SPACE     ~= (\s)
_NEWLINE   ~= (\n|\r|\r\n)
IDENTIFIER ~= ([a-z]+)	MODULE=module \
                       END=end
`

// Replace buy build script
var slexTemplate = `
@@@@data@@@@
`
