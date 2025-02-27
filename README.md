# Simple lexer

Simple lexer is tool to transform a text into token to use with [goyacc](https://github.com/golang/tools/tree/master/cmd/goyacc).

To build, you need run `build.sh` cause `lexer/lexer.go` file is copied into `cmd/command_line_options.go` file.

## How to use it?

First, you need create a file with pattern of token, in X file (see `demo/basic.x`).

You can define some rules, like:
```
// This is comment
NUMBER     == 123
_SPACE     ~= (\s)
_NEWLINE   ~= (\n|\r|\r\n)
IDENTIFIER ~= ([a-z]+)	MODULE=module \
                       END=end
```

Each line contains 3 parts.

First part is the identifier use to link with goyacc. In goyacc file (`*.y`), this is `%token` identifier. If identifier start by `_`, token is ignored.

After, you have the type of token:
 * `==`: to check an hard value,
 * `=>`: to call a go function to generate a token (see the `skipComment` function in `lexer/lexer_test.go`),
 * `~=`: use a regex

The last part is value.

With regex, you can add sub-pattern. For example, you want check if identifier is keyword or something else.
In this case, after value, you can add function name or couple of identifier/value.

In case of regex, if regex contains space (don't sure that make sense), you can enclose with single-quote or double-quote.

If line is too long, you can split it by using `\` at end of line.
