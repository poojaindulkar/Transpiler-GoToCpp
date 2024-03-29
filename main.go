package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	_ "fmt"
	"regexp"

	"github.com/poojaindulkar/Transpiler-GoToCpp/src/constants"
	"github.com/poojaindulkar/Transpiler-GoToCpp/src/replacements"
	"github.com/poojaindulkar/Transpiler-GoToCpp/src/utils"
)

/************************************************************ Lexer Start *******************************************************************/

type Token struct {
	Type  string
	Value string
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

	re := regexp.MustCompile(`([[:space:]]+)|([a-zA-Z]+)|([0-9]+)|([+\-*/%=])|(\(|\))|(\{|\})`)
	if match := re.FindStringIndex(l.input[l.pos:]); match != nil {
		l.pos += match[1]
		tokenType := ""
		tokenValue := l.input[l.pos-match[1] : l.pos]

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
		default:
			tokenType = "identifier"
		}

		return &Token{Type: tokenType, Value: tokenValue}
	}

	// If no regular expression matches, return an error token

	return &Token{Type: "Error", Value: string(l.input[l.pos])}
}

func LEX(input string) {

	var res []Token

	//	i:=0
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

		if token.Value == " " || token.Value == "\n" || token.Value == "\t" {
			continue
		}

		//		fmt.Printf("%s \n", token.Value)

		res = append(res, Token{
			Type:  token.Type,
			Value: token.Value,
		})
	}
	fmt.Println(res)

	println("__________")

	Ast(res)

}

/************************************************************ Lexer End *******************************************************************/

/************************************************************ Parser Start *******************************************************************/
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

func Ast(T []Token) {
	tokens := T

	root := ParseAST(tokens)

	// Print the AST tree.
	PrintAST(root, 0)
}

/************************************************************ Parser End *******************************************************************/

func go2cpp(source string) string {
	functionVarMap := map[string]string{} // variable names encountered in the function so far, and their corresponding smart names
	inMultilineString := false
	debugOutput := false
	lines := []string{}
	currentReturnType := ""
	currentFunctionName := ""
	inImport := false
	inVar := false
	inType := false
	inConst := false
	inHashMap := false
	hashKeyType := ""
	curlyCount := 0

	encounteredHashMaps := []string{}

	encounteredStructNames := []string{}
	inStruct := false
	// usePrettyPrint := false
	closingBracketNeedsASemicolon := false
	for _, line := range strings.Split(source, "\n") {

		if debugOutput {
			fmt.Fprintf(os.Stderr, "%s\n", line)
		}

		newLine := line

		trimmedLine := replacements.StripSingleLineComment(strings.TrimSpace(line))

		if strings.HasPrefix(trimmedLine, "//") {
			lines = append(lines, trimmedLine)
			continue
		}
		if strings.HasSuffix(trimmedLine, ";") {
			trimmedLine = trimmedLine[:len(trimmedLine)-1]
		}
		if len(trimmedLine) == 0 {
			lines = append(lines, newLine)
			continue
		}
		// Keep track of how deep we are into curly brackets
		curlyCount += (strings.Count(trimmedLine, "{") - strings.Count(trimmedLine, "}"))
		if inImport && strings.Contains(trimmedLine, ")") {
			inImport = false
			continue
		} else if inImport {
			continue
		} else if inVar && strings.Contains(trimmedLine, ")") {
			inVar = false
			continue
		} else if inType && strings.Contains(trimmedLine, ")") {
			inType = false
			continue
		} else if inConst && strings.Contains(trimmedLine, ")") {
			inConst = false
			continue
		} else if inHashMap && trimmedLine == "}" {
			inHashMap = false
			newLine = trimmedLine + ";"
		} else if inVar || (inStruct && trimmedLine != "}") {
			var varNames []string
			newLine, varNames = utils.VarDeclarations(trimmedLine)
			if strings.HasSuffix(newLine, "{") {
				closingBracketNeedsASemicolon = true
			}

			if inStruct {
				// Gathering variable names from this struct
				encounteredStructNames = append(encounteredStructNames, varNames...)
			}
		} else if inType {
			prevInStruct := inStruct
			newLine, inStruct = utils.TypeDeclaration(trimmedLine)
			if !prevInStruct && inStruct {
				// Entering struct, reset the slice that is used to gather variable names
				encounteredStructNames = []string{}
			}
		} else if inConst {
			newLine = utils.ConstDeclaration(line)
		} else if inHashMap && !inMultilineString {
			newLine = replacements.HashElements(trimmedLine, hashKeyType, false)
		} else if strings.HasPrefix(trimmedLine, "func") {
			functionVarMap = map[string]string{}
			newLine, currentReturnType, currentFunctionName = utils.FunctionSignature(trimmedLine)
		} else if strings.HasPrefix(trimmedLine, "for") {
			newLine = utils.ForLoop(line, encounteredHashMaps)
		} else if strings.HasPrefix(trimmedLine, "switch") {
			newLine = utils.Switch(line)
		} else if strings.HasPrefix(trimmedLine, "case") {
			newLine = utils.Case(line)
		} else if strings.HasPrefix(trimmedLine, "return") {
			if strings.HasPrefix(currentReturnType, constants.TupleType) {
				elems := strings.SplitN(newLine, "return ", 2)
				newLine = "return " + currentReturnType + "{" + elems[1] + "};"
				//} else {
				// Just use the standard tuple
			}
		} else if strings.HasPrefix(trimmedLine, "fmt.Print") || strings.HasPrefix(trimmedLine, "print") {
			// _ is if "pretty print" functionality may be needed, for non-literal strings and numbers
			// var pp bool
			newLine, _ = utils.PrintStatement(trimmedLine)
			// if pp {
			// 	usePrettyPrint = true
			// }
		} else if strings.Contains(trimmedLine, "=") && !strings.HasPrefix(trimmedLine, "var ") && !strings.HasPrefix(trimmedLine, "if ") && !strings.HasPrefix(trimmedLine, "const ") && !strings.HasPrefix(trimmedLine, "type ") {
			elem := strings.SplitN(trimmedLine, "=", 2)
			left := strings.TrimSpace(elem[0])
			declarationAssignment := false
			if strings.HasSuffix(left, ":") {
				declarationAssignment = true
				left = left[:len(left)-1]
			}
			right := strings.TrimSpace(elem[1])
			if strings.HasPrefix(right, "&") && strings.Contains(right, "{") && strings.Contains(right, "}") {
				right = "new " + right[1:]
			}
			if strings.Contains(left, ",") {
				if strings.Contains(left, "_") {
					varNames := strings.Split(left, ",")
					for _, name := range varNames {
						name = strings.TrimSpace(name)
						if value, found := functionVarMap[name]; found {
							// The key already exists, update the value
							if value == name {
								// Add a "0"
								functionVarMap[name] = value + "0"
							} else {
								// Increase the number in the current value by 1
								num := utils.TrailingNumber(value)
								num++
								functionVarMap[name] = name + strconv.Itoa(num)
							}
						} else {
							// The key does not exist, just add the name as it is
							functionVarMap[name] = name
						}
					}
					// The "varInFunction" map should now have been updated correctly, so use that

					useVarNames := []string{}
					for _, name := range varNames {
						name = strings.TrimSpace(name)
						useVarNames = append(useVarNames, functionVarMap[name])
					}
					newLine = "auto [" + strings.Join(useVarNames, ", ") + "] = " + right
					//fmt.Println("function var map:", functionVarMap)
					//panic(newLine)
				} else {
					newLine = "auto [" + left + "] = " + right
				}
			} else if declarationAssignment {
				if strings.HasPrefix(right, "[]") {
					if !strings.Contains(right, "{") {
						fmt.Fprintln(os.Stderr, "UNRECOGNIZED LINE: "+trimmedLine)
						//newLine = line

					}
					theType := replacements.TypeReplace(utils.LeftBetween(right, "]", "{"))
					fields := strings.SplitN(right, "{", 2)
					newLine = theType + " " + strings.TrimSpace(left) + "[] {" + fields[1]
				} else if strings.HasPrefix(right, "map[") {
					hashName := strings.TrimSpace(left)
					encounteredHashMaps = append(encounteredHashMaps, hashName)

					keyType := replacements.TypeReplace(utils.LeftBetween(right, "map[", "]"))
					valueType := replacements.TypeReplace(utils.LeftBetween(right, "]", "{"))

					closingBracket := strings.HasSuffix(strings.TrimSpace(right), "}")
					if !closingBracket {
						inHashMap = true
						hashKeyType = keyType
						newLine = "std::unordered_map<" + keyType + ", " + valueType + "> " + hashName + " {"
					} else {
						elements := utils.LeftBetween(right, "{", "}")
						newLine = "std::unordered_map<" + keyType + ", " + valueType + "> " + hashName + " " + replacements.HashElements(elements, keyType, false)
					}
				} else {
					varName := strings.TrimSpace(left)
					if value, found := functionVarMap[varName]; found {
						varName = value
					}

					for k, v := range functionVarMap {
						right = strings.Replace(right, k, v, -1)
					}

					newLine = "auto " + varName + " = " + strings.TrimSpace(right)
				}
			} else {
				newLine = left + " = " + right
			}
		} else if strings.HasPrefix(trimmedLine, "package ") {
			continue
		} else if strings.HasPrefix(trimmedLine, "import") {
			if strings.Contains(trimmedLine, "(") {
				inImport = true
			}
			if strings.Contains(trimmedLine, ")") {
				inImport = false
			}
			continue
			// } else if strings.HasPrefix(trimmedLine, "defer ") {
			// 	newLine = DeferCall(line)
		} else if strings.HasPrefix(trimmedLine, "if ") {
			newLine = utils.IfSentence(line)
			// TODO: Short variable names utils.Has the potential to ruin if expressions this way, do a smarter replacement
			// TODO: Also do this for for loops, switches and other cases where this makes sense
			for k, v := range functionVarMap {
				newLine = strings.Replace(newLine, k, v, -1)
			}
		} else if strings.HasPrefix(trimmedLine, "} else if ") {
			newLine = utils.ElseIfSentence(line)
		} else if trimmedLine == "var (" {
			inVar = true
			continue
		} else if trimmedLine == "type (" {
			inType = true
			continue
		} else if trimmedLine == "const (" {
			inConst = true
			continue
		} else if strings.HasPrefix(trimmedLine, "var ") {
			// Ignore variable name since it's not in a struct
			newLine, _ = utils.VarDeclarations(line)
			if strings.HasSuffix(newLine, "{") {
				closingBracketNeedsASemicolon = true
			}
		} else if strings.HasPrefix(trimmedLine, "type ") {
			newLine, inStruct = utils.TypeDeclaration(trimmedLine)
		} else if strings.HasPrefix(trimmedLine, "const ") {
			newLine = utils.ConstDeclaration(trimmedLine)
		} else if trimmedLine == "fallthrough" {
			newLine = "goto " + utils.LabelName() + "; // fallthrough"
			constants.SwitchLabel = utils.LabelName()
			constants.LabelCounter++
		} else if constants.UnfinishedDeferFunction && trimmedLine == "}()" {
			constants.UnfinishedDeferFunction = false
			newLine = "});"
		} else if trimmedLine == "default:" {
			newLine = "} else { // default case"
			if constants.SwitchLabel != "" {
				newLine += "\n" + constants.SwitchLabel + ":"
				constants.SwitchLabel = ""
			}
		}

		if constants.CppHasStdFormat {
			// Special case for fmt.Sprintf -> std::format
			if strings.Contains(newLine, "fmt.Sprintf(") && strings.Contains(newLine, "%v") {
				newLine = strings.Replace(strings.Replace(newLine, "%v", "{}", -1), "fmt.Sprintf(", "std::format(", -1)
			}
		}

		if currentFunctionName == "main" && trimmedLine == "}" && curlyCount == 0 { // curlyCount utils.Has already been decreased for this line
			newLine = strings.Replace(trimmedLine, "}", "return 0;\n}", 1)
		}

		if strings.HasSuffix(trimmedLine, "}") {
			// If the struct is being closed, add a semicolon
			if inStruct {
				// Create a _str() method for this struct
				newLine = utils.CreateStrMethod(encounteredStructNames) + newLine + ";"

				inStruct = false
			} else if closingBracketNeedsASemicolon {
				newLine += ";"
				closingBracketNeedsASemicolon = false
			}
			newLine += "\n"
		}
		if (!strings.HasSuffix(newLine, ";") && !utils.Has(constants.Endings, utils.Lastchar(trimmedLine)) || strings.Contains(trimmedLine, "=")) && !strings.HasPrefix(trimmedLine, "//") && (!utils.Has(constants.Endings, utils.Lastchar(newLine)) && !strings.Contains(newLine, "//")) {
			if !inMultilineString {
				newLine += ";"
			}
		}

		// multiline strings
		for strings.Contains(newLine, "`") {
			if !inMultilineString {
				if strings.HasSuffix(newLine, ";") {
					newLine = newLine[:len(newLine)-1]
				}
				newLine = strings.Replace(newLine, "`", "R\"(", 1)
				inMultilineString = true
			} else {
				newLine = strings.Replace(newLine, "`", ")\"", 1)
				//if !strings.HasSuffix(newLine, ",") {
				newLine += ";"
				//}
				inMultilineString = false
			}
		}

		lines = append(lines, newLine)
	}
	output := strings.Join(lines, "\n")

	// The order matters
	output = replacements.WholeProgramReplace(output)

	// The order matters
	// output = AddFunctions(output, usePrettyPrint, len(encounteredStructNames) > 0)
	// output = AddIncludes(output)

	return output
}

func main() {

	debug := false
	compile := true
	clangFormat := true

	inputFilename := "./go_test_code/test_first_1.txt"

	var sourceData []byte
	var err error
	if inputFilename != "" {
		sourceData, err = ioutil.ReadFile(inputFilename)
	} else {
		sourceData, err = ioutil.ReadAll(os.Stdin)
	}

	ip := (string(sourceData))
	LEX(ip)

	if err != nil {
		log.Fatal(err)
	}
	if debug {
		fmt.Println(go2cpp(string(sourceData)))
		return
	}

	cppSource := "output.txt"
	if clangFormat {
		cmd := exec.Command("clang-format", "-style={BasedOnStyle: Webkit, ColumnLimit: 99}")
		cmd.Stdin = strings.NewReader(go2cpp(string(sourceData)))
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			// log.Println("clang-format is not available, the output will look ugly!")
			cppSource = go2cpp(string(sourceData))
		} else {
			cppSource = out.String()
		}
	} else {
		cppSource = go2cpp(string(sourceData))
	}

	if !compile {
		fmt.Println(cppSource)
		return
	}

	tempFile, err := ioutil.TempFile("", "go2cpp*")
	if err != nil {
		log.Fatal(err)
	}
	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName)

	// Compile the string in cppSource
	cpp := "g++"
	if cppenv := os.Getenv("CXX"); cppenv != "" {
		cpp = cppenv
	}
	cmd2 := exec.Command(cpp, "-x", "c++", "-std=c++2a", "-O2", "-pipe", "-fPIC", "-Wfatal-errors", "-fpermissive", "-s", "-o", tempFileName, "-")
	cmd2.Stdin = strings.NewReader(cppSource)
	var compiled bytes.Buffer
	var errors bytes.Buffer
	cmd2.Stdout = &compiled
	cmd2.Stderr = &errors
	err = cmd2.Run()
	if err != nil {

		cppSource = `

template <typename T>
	void _format_output(std::ostream& out, const T& str) 
	{	
		out << str;
	}` + cppSource

		cppSource = `#include <bits/stdc++.h>` + cppSource
		cppSource = formatting(cppSource)
		//fmt.Println(cppSource)

		outputFileName := "cpp_output/" + strings.TrimSuffix(strings.Split(inputFilename, "/")[2], ".txt") + ".cpp"
		writeFile(outputFileName, cppSource)

		// fmt.Println("Errors:")
		// fmt.Println(errors.String())
		//log.Fatal(err)
	}
	compiledBytes, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		log.Fatal(err)
	}
	outputFilename := ""
	if len(os.Args) > 3 {
		outputFilename = os.Args[3]
	}
	if outputFilename != "" {
		err = ioutil.WriteFile(outputFilename, compiledBytes, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		//fmt.Println(cppSource)
	}

}

func writeFile(filename, data string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func formatting(data string) string {
	// temporary
	ans := data
	ans = strings.Replace(data, "Math.Sqrt", "sqrt", -1)
	ans = strings.Replace(ans, "float64", "double", -1)
	ans = strings.Replace(ans, "+ =", "+=", -1)
	ans = strings.Replace(ans, "- =", "-=", -1)
	ans = strings.Replace(ans, "* =", "*=", -1)
	ans = strings.Replace(ans, "/ =", "/=", -1)

	return ans

}
