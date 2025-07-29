package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	tkn_brace_open    = "brace_open"    // {
	tkn_brace_close   = "brace_close"   // }
	tkn_bracket_open  = "bracket_open"  // [
	tkn_bracket_close = "bracket_close" // ]
	tkn_colon         = "colon"         // :
	tkn_comma         = "comma"         // ,
	tkn_string        = "string"
	tkn_number        = "number"
	tkn_true          = "true"
	tkn_false         = "false"
	tkn_null          = "null"
)

type Token struct {
	Type, Value string
}

func Tokenize(jsonString string) ([]Token, error) {
	var (
		curr   = 0
		strLen = len(jsonString)
		tokens = []Token{}
	)
	for curr < strLen {
		char := jsonString[curr]
		if unicode.IsSpace(rune(char)) {
			curr++ // 跳过空白字符
			continue
		}
		switch char {
		case '{':
			tokens = append(tokens, Token{Type: tkn_brace_open, Value: "{"})
			curr++
		case '}':
			tokens = append(tokens, Token{Type: tkn_brace_close, Value: "}"})
			curr++

		case '[':
			tokens = append(tokens, Token{Type: tkn_bracket_open, Value: "["})
			curr++

		case ']':
			tokens = append(tokens, Token{Type: tkn_bracket_close, Value: "]"})
			curr++

		case ':':
			tokens = append(tokens, Token{Type: tkn_colon, Value: ":"})
			curr++

		case ',':
			tokens = append(tokens, Token{Type: tkn_comma, Value: ","})
			curr++

		case '"':
			curr++ // 跳过开始的引号
			start := curr
			for curr < strLen && jsonString[curr] != '"' {
				if jsonString[curr] == '\\' { // 处理转义字符
					curr++
				}
				curr++
			}
			str := jsonString[start:curr]
			if curr >= strLen {
				return nil, fmt.Errorf("Unterminated string literal: %s\n", str)
			}
			tokens = append(tokens, Token{Type: tkn_string, Value: str})
			curr++ // 跳过结束的引号
		default:
			res := jsonString[curr:]
			if strings.HasPrefix(res, "true") {
				tokens = append(tokens, Token{Type: tkn_true, Value: "true"})
				curr += 4
			} else if strings.HasPrefix(res, "false") {
				tokens = append(tokens, Token{Type: tkn_false, Value: "false"})
				curr += 5
			} else if strings.HasPrefix(res, "null") {
				tokens = append(tokens, Token{Type: tkn_null, Value: "null"})
				curr += 4
			} else if unicode.IsDigit(rune(char)) || char == '-' { // 处理数字
				start := curr
				curr++ // 跳过第一个数字字符

				hasDot := false
				hasExp := false
				expDigits := 0
				for curr < strLen {
					c := jsonString[curr]
					if unicode.IsDigit(rune(c)) {
						curr++ // 继续读取数字
						if hasExp {
							expDigits++
						}
					} else if c == '.' {
						if hasDot || hasExp {
							return nil, fmt.Errorf("invalid number: multiple dots or dot after exponent at position %d", curr)
						}
						hasDot = true
						curr++ // 跳过小数点
					} else if c == 'e' || c == 'E' {
						if hasExp {
							return nil, fmt.Errorf("invalid number: multiple exponents at position %d", curr)
						}
						hasExp = true
						curr++ // 跳过指数符号
						if curr < strLen && (jsonString[curr] == '+' || jsonString[curr] == '-') {
							curr++ // 跳过可选的正负号
						}
						expDigits = 0
					} else {
						break // 不是数字字符，结束读取
					}
				}
				number := jsonString[start:curr]
				if hasExp && expDigits == 0 {
					return nil, fmt.Errorf("invalid number: exponent without digits at position %d", curr)
				}
				if _, err := strconv.ParseFloat(number, 64); err != nil {
					return nil, fmt.Errorf("invalid number: %s at position %d", number, start)
				}
				tokens = append(tokens, Token{Type: tkn_number, Value: number})
			} else {
				return nil, fmt.Errorf("unexpected character '%c' at position %d", char, curr)
			}
		}
	}
	return tokens, nil
}

func printTokens(tokens []Token) {
	fmt.Printf("%-14s | %s\n", "Type", "Value")
	fmt.Println(strings.Repeat("-", 50))
	for _, token := range tokens {
		fmt.Printf("Type: %s, Value: %s\n", token.Type, token.Value)
	}
}

// https://dev.to/balapriya/abstract-syntax-tree-ast-explained-in-plain-english-1h38
type ASTNode interface {
	Type() string
}

type ObjectNode struct {
	Value map[string]ASTNode
}

func (o *ObjectNode) Type() string {
	return "Object"
}

type ArrayNode struct {
	Value []ASTNode
}

func (a *ArrayNode) Type() string {
	return "Array"
}

type StringNode struct {
	Value string
}

func (s *StringNode) Type() string {
	return "String"
}

type NumberNode struct {
	Value float64
}

func (n *NumberNode) Type() string {
	return "Number"
}

type BooleanNode struct {
	Value bool
}

func (b *BooleanNode) Type() string {
	return "Boolean"
}

type NullNode struct{}

func (n *NullNode) Type() string {
	return "Null"
}

func Parser(tokens []Token) (ASTNode, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}
	curr := 0
	return parseValue(tokens, &curr)
}

func parseValue(tokens []Token, curr *int) (ASTNode, error) {
	if *curr >= len(tokens) {
		return nil, fmt.Errorf("unexpected end of tokens")
	}
	token := tokens[*curr]
	switch token.Type {
	case tkn_string:
		*curr++
		return &StringNode{Value: token.Value}, nil
	case tkn_number:
		num, _ := strconv.ParseFloat(token.Value, 64)
		*curr++
		return &NumberNode{Value: num}, nil
	case tkn_true:
		*curr++
		return &BooleanNode{Value: true}, nil
	case tkn_false:
		*curr++
		return &BooleanNode{Value: false}, nil
	case tkn_null:
		*curr++
		return &NullNode{}, nil
	case tkn_brace_open:
		return parseObject(tokens, curr)
	case tkn_bracket_open:
		return parseArray(tokens, curr)
	default:
		return nil, fmt.Errorf("unexpected token type: %s", token.Type)
	}
}

func parseObject(tokens []Token, curr *int) (ASTNode, error) {
	node := &ObjectNode{Value: make(map[string]ASTNode)}
	*curr++ // 跳过 '{'
	for *curr < len(tokens) && tokens[*curr].Type != tkn_brace_close {
		curToken := tokens[*curr]
		if curToken.Type != tkn_string {
			return nil, fmt.Errorf("expected string key, got %s", curToken.Type)
		}
		key := curToken.Value
		*curr++ // 跳过键
		if *curr >= len(tokens) || tokens[*curr].Type != tkn_colon {
			return nil, fmt.Errorf("expected colon after key %s", key)
		}
		*curr++ // 跳过 ':'
		valueNode, err := parseValue(tokens, curr)
		if err != nil {
			return nil, fmt.Errorf("error parsing value for key %s: %v", key, err)
		}
		node.Value[key] = valueNode
		if *curr < len(tokens) && tokens[*curr].Type == tkn_comma {
			*curr++ // 跳过 ','
		}
	}
	if *curr >= len(tokens) || tokens[*curr].Type != tkn_brace_close {
		return nil, fmt.Errorf("expected '}' to close object")
	}
	*curr++ // 跳过 '}'
	return node, nil
}

func parseArray(tokens []Token, curr *int) (ASTNode, error) {
	node := &ArrayNode{Value: []ASTNode{}}
	*curr++ // 跳过 '['
	for *curr < len(tokens) && tokens[*curr].Type != tkn_bracket_close {
		valueNode, err := parseValue(tokens, curr)
		if err != nil {
			return nil, fmt.Errorf("error parsing array value: %v", err)
		}
		node.Value = append(node.Value, valueNode)
		if *curr < len(tokens) && tokens[*curr].Type == tkn_comma {
			*curr++ // 跳过 ','
		}
	}
	if *curr >= len(tokens) || tokens[*curr].Type != tkn_bracket_close {
		return nil, fmt.Errorf("expected ']' to close array")
	}
	*curr++ // 跳过 ']'
	return node, nil
}

func main() {
	tks, err := Tokenize(`{"name": "John", "age": 30, "is_student": false, "courses": ["Math", "Science"], "address": {"city": "New York", "zip": "10001"}, "graduated": null}`)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	printTokens(tks)
}
