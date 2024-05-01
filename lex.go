package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Token struct {
	Type  string
	Value string
	Line  int
	Pos   int
}

type Lexer struct {
	input  string
	pos    int
	tokens []Token
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) NextToken() *Token {
	if l.pos >= len(l.input) {
		return nil
	}

	// Match regular expressions to identify tokens

	re := regexp.MustCompile(`([[:space:]]+)|([a-zA-Z]+)|([0-9]+)|([+\-*/%=])|(\(|\))|(\{|\})|(\n)`)
	if match := re.FindStringIndex(l.input[l.pos:]); match != nil {
		l.pos += match[1]
		tokenType := ""
		tokenValue := l.input[l.pos-match[1] : l.pos]

		// Update line number and position
		line := strings.Count(l.input[:l.pos], "\n") + 1
		pos := l.pos - strings.LastIndex(l.input[:l.pos], "\n")

		if strings.HasPrefix(tokenValue, "\r") || strings.HasPrefix(tokenValue, "\n") {
			// return &Token{Type: tokenType, Value: "\n", Line: line, Pos: pos}
			tokenValue="\\n"
		}
		index:=strings.LastIndex(tokenValue,"\"\r\n")
		if index != -1 {
			tokenValue =strings.TrimRight(tokenValue[:index+1], "\r\n")
			
		}
		switch tokenValue {
		case "package main":
			tokenType = "package"
		case "func":
			tokenType = "Func"
		case "main":
			tokenType = "main"
		case "(":
			tokenType = "OpenParen"
		case ")":
			tokenType = "CloseParen"
		case "{":
			tokenType = "OpenBrace"
		case "}":
			tokenType = "CloseBrace"
		case "+":
			tokenType = "Plus"
		case "-":
			tokenType = "Minus"
		case "*":
			tokenType = "Asterisk"
		case "/":
			tokenType = "Slash"
		case "=":
			tokenType = "Equals"
		case " ":
			tokenType = "blank space"
		case "\"":
			tokenType = "double quotes"
		case "fmt.Println":
			tokenType = "print function"
		case "\\n": 
			tokenType = "newline"
		default:
			tokenType = "identifier"
		}
		// fmt.Println()
		return &Token{Type: tokenType, Value: tokenValue, Line: line, Pos: pos}
	}

	// If no regular expression matches, return an error token

	return &Token{Type: "Error", Value: string(l.input[l.pos]), Line: 0, Pos: 0}
}

func LEX(input string) error {
	fmt.Println(input)
	outputFileName := "./analysis/lexer.txt" // Specify the output file name here

	file, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var res []Token

	l := NewLexer(input)
	t := true
	for token := l.NextToken(); token != nil; token = l.NextToken() {
		if token.Value == "(" {
			t = false
		} else if token.Value == ")" {
			t = true
		}

		if token.Value != " " && t == true {
			//fmt.Printf("%s: %s\n", token.Type, token.Value)
		} else if t == false {
			//fmt.Printf("%s: %s\n", token.Type, token.Value)
		}

		if token.Value == "\n" || token.Value == " " || token.Value == "\t" {

			continue
		}
		// if(token.Value==""){
		// 	fmt.Println("no value")
		// 	continue;
		// }
		// fmt.Print("Token Value: ",token.Value,"\n")
		res = append(res, Token{
			Type:  token.Type,
			Value: token.Value,
			Line:  token.Line,
			Pos:   token.Pos,
		})
	}
	// fmt.Println(res)
	// Write to the file with box-like structure
	fmt.Fprintln(file, "┌─────────────────┬────────────┬───────┬─────────┐")
	fmt.Fprintf(file, "│ %-15s │ %-10s │ %-5s │ %-7s│\n", "Type", "Value", "Line", "Position")
	fmt.Fprintln(file, "├─────────────────┼────────────┼───────┼─────────┤")
	for _, token := range res {
		fmt.Fprintf(file, "│ %-15s │ %-10s │ %-5d │ %-7d │\n", token.Type, token.Value, token.Line, token.Pos)
	}
	fmt.Fprintln(file, "└─────────────────┴────────────┴───────┴─────────┘")

	fmt.Println("Output written to", outputFileName)
	Ast(res)
	return nil
}
