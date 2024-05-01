package main

import (
	"fmt"
	"strings"
	"os"
)

type Tokenx struct {
	Type  string
	Value string
}

type Node struct {
	Type     string
	Value    string
	Children []Node
}

func NewNode(Type string, Value string) *Node {
	return &Node{Type: Type, Value: Value}
}

func AddChild(node *Node, child *Node) {
	node.Children = append(node.Children, *child)
}


func ParseAST(tokens []Token) *Node {
	root := NewNode("Root", "")

	for _, token := range tokens {
		switch token.Type {
		case "package":
			AddChild(root, NewNode("Idenpackagetifier", token.Value))
		case "StringLiteral":
			AddChild(root, NewNode("StringLiteral", token.Value))
		case "NumberLiteral":
			AddChild(root, NewNode("NumberLiteral", token.Value))
		case "func":
			AddChild(root, NewNode("func", token.Value))
		case "main":
			AddChild(root, NewNode("main", token.Value))

		case "OpenParen":
			AddChild(root, NewNode("OpenParen", token.Value))
		case "CloseParen":
			AddChild(root, NewNode("CloseParen", token.Value))
		case "OpenBrace":
			AddChild(root, NewNode("OpenBrace", token.Value))
		case "CloseBrace":
			AddChild(root, NewNode("CloseBrace", token.Value))
		case "Plus":
			AddChild(root, NewNode("Plus", token.Value))
		case "Minus":
			AddChild(root, NewNode("Minus", token.Value))
		case "Asterisk":
			AddChild(root, NewNode("Asterisk", token.Value))
		case "Slash":
			AddChild(root, NewNode("Slash", token.Value))
		case "Equals":
			AddChild(root, NewNode("Equals", token.Value))
		case "blank space":
			AddChild(root, NewNode("blank space", token.Value))
		case "double quotes":
			AddChild(root, NewNode("double quotes", token.Value))
		case "print function":
			AddChild(root, NewNode("print function", token.Value))
		case "newline":
			AddChild(root, NewNode("newline", token.Value))
		default:
			AddChild(root, NewNode("identifier", token.Value))
		}

	}

	return root
}

func PrintAST(node *Node, indent int) {
	fmt.Printf("%s%s: %s\n", strings.Repeat(" ", indent), node.Type, node.Value)

	for _, child := range node.Children {
		PrintAST(&child, indent+1)
	}
}
func PrintASTToFile(node *Node, file *os.File, indent int) {
	if node == nil {
		return
	}

	// Print the node
	fmt.Fprintf(file, "%s%s: %s\n", strings.Repeat(" ", indent), node.Type, node.Value)

	// Print lines connecting the node to its children
	for i, child := range node.Children {
		if i == len(node.Children)-1 {
			// Last child, no need for a connecting line
			fmt.Fprintf(file, "%s \n", strings.Repeat(" ", indent+1))
		} else {
			// Print a connecting line
			fmt.Fprintf(file, "%s├── ", strings.Repeat(" ", indent))
		}

		// Recursively print the child node
		PrintASTToFile(&child, file, indent+1)
	}
}

func AstToFile(T []Token) error {
	root := ParseAST(T)

	// Open the file for writing
	file, err := os.Create("analysis/parser.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	// Print the AST tree to the file
	PrintASTToFile(root, file, 0)

	fmt.Println("AST written to analysis/parser.txt")
	return nil
}
func Ast(T []Token) {
	tokens := T

	root := ParseAST(tokens)

	// Print the AST tree.
	PrintAST(root, 0)
	err := AstToFile(tokens)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
