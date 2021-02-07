First, generate `demo.go` file:
`$ goyacc -o demo.go -p Basic basic.y`

Then generate lexer and data :
`$ slex generate -i .basic.x -o basic.go`
That generate a file called `basic.go` (it's copy of `lexer/lexer.go` file) and print output below:
```
[]TokenEntry{
	NewHardValueToken("PRINT", "print", PRINT),
	NewRegexValueToken("IDENTIFIER", "([a-zA-Z]+)", IDENTIFIER),
	NewHardValueToken("ADD", "+", ADD),
	NewRegexValueToken("NUMBER", "([0-9]+)", NUMBER),
	NewHardValueToken("EQUAL", "=", EQUAL),
	NewRegexValueToken("_SPACE", "(\\s)", -1),
}
```
Replace `%%%TOKEN_LIST%%%` by previous data in `demo.go` file, the run `go run .`. You see `128`.
