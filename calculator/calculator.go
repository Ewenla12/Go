package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

//  TOKENIZER

type TokenType string

const (
	NUMBER TokenType = "NUMBER"
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	MUL    TokenType = "*"
	DIV    TokenType = "/"
	MOD    TokenType = "%"
	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	EOF    TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(input string) ([]Token, error) {
	var tokens []Token
	i := 0
	runes := []rune(input)

	for i < len(runes) {
		ch := runes[i]

		if unicode.IsSpace(ch) {
			i++
			continue
		}

		if unicode.IsDigit(ch) || (ch == '.' && i+1 < len(runes) && unicode.IsDigit(runes[i+1])) {
			j := i
			hasDot := false
			for j < len(runes) && (unicode.IsDigit(runes[j]) || (runes[j] == '.' && !hasDot)) {
				if runes[j] == '.' {
					hasDot = true
				}
				j++
			}
			tokens = append(tokens, Token{NUMBER, string(runes[i:j])})
			i = j
			continue
		}

		switch ch {
		case '+':
			tokens = append(tokens, Token{PLUS, "+"})
		case '-':
			tokens = append(tokens, Token{MINUS, "-"})
		case '*', 'x', 'X':
			tokens = append(tokens, Token{MUL, "*"})
		case '/':
			tokens = append(tokens, Token{DIV, "/"})
		case '%':
			tokens = append(tokens, Token{MOD, "%"})
		case '(':
			tokens = append(tokens, Token{LPAREN, "("})
		case ')':
			tokens = append(tokens, Token{RPAREN, ")"})
		default:
			return nil, fmt.Errorf("unknown character: %q", ch)
		}
		i++
	}

	tokens = append(tokens, Token{EOF, ""})
	return tokens, nil
}

//  PARSER  (recursive descent — BODMAS order)
//  expression → term ((+ | -) term)*
//  term       → unary ((* | / | %) unary)*
//  unary      → - primary | primary
//  primary    → NUMBER | ( expression )

type Parser struct {
	tokens []Token
	pos    int
}

func (p *Parser) peek() Token    { return p.tokens[p.pos] }
func (p *Parser) consume() Token { t := p.tokens[p.pos]; p.pos++; return t }

func (p *Parser) expression() (float64, error) {
	left, err := p.term()
	if err != nil {
		return 0, err
	}
	for p.peek().Type == PLUS || p.peek().Type == MINUS {
		op := p.consume()
		right, err := p.term()
		if err != nil {
			return 0, err
		}
		if op.Type == PLUS {
			left += right
		} else {
			left -= right
		}
	}
	return left, nil
}

func (p *Parser) term() (float64, error) {
	left, err := p.unary()
	if err != nil {
		return 0, err
	}
	for p.peek().Type == MUL || p.peek().Type == DIV || p.peek().Type == MOD {
		op := p.consume()
		right, err := p.unary()
		if err != nil {
			return 0, err
		}
		switch op.Type {
		case MUL:
			left *= right
		case DIV:
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		case MOD:
			if right == 0 {
				return 0, fmt.Errorf("modulo by zero")
			}
			left = math.Mod(left, right)
		}
	}
	return left, nil
}

func (p *Parser) unary() (float64, error) {
	if p.peek().Type == MINUS {
		p.consume()
		val, err := p.primary()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}
	return p.primary()
}

func (p *Parser) primary() (float64, error) {
	tok := p.peek()
	if tok.Type == NUMBER {
		p.consume()
		return strconv.ParseFloat(tok.Value, 64)
	}
	if tok.Type == LPAREN {
		p.consume()
		val, err := p.expression()
		if err != nil {
			return 0, err
		}
		if p.peek().Type != RPAREN {
			return 0, fmt.Errorf("missing closing )")
		}
		p.consume()
		return val, nil
	}
	return 0, fmt.Errorf("unexpected token: %q", tok.Value)
}

//  EVALUATE

func evaluate(input string) (float64, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("empty expression")
	}
	tokens, err := tokenize(input)
	if err != nil {
		return 0, err
	}
	parser := &Parser{tokens: tokens}
	result, err := parser.expression()
	if err != nil {
		return 0, err
	}
	if parser.peek().Type != EOF {
		return 0, fmt.Errorf("unexpected token: %q", parser.peek().Value)
	}
	return result, nil
}

//  FORMAT — strips trailing zeros

func format(f float64) string {
	if f == math.Trunc(f) && !math.IsInf(f, 0) {
		return fmt.Sprintf("%.0f", f)
	}
	s := strconv.FormatFloat(f, 'f', 10, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

//  MAIN

func main() {
	reader := bufio.NewReader(os.Stdin)

	var lastResult float64
	hasLast := false

	for {
		fmt.Print("\n> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// if line starts with an operator, prepend ans automatically
		// e.g. *100 becomes ans*100
		if hasLast && len(line) > 0 {
			first := rune(line[0])
			if first == '+' || first == '-' || first == '*' || first == '/' || first == '%' {
				line = "ans" + line
			}
		}

		// "ans" lets you chain the last result into the next expression
		if hasLast {
			line = strings.ReplaceAll(line, "ans", format(lastResult))
		}

		result, err := evaluate(line)
		if err != nil {
			fmt.Println("  Error:", err)
			continue
		}

		lastResult = result
		hasLast = true
		fmt.Println("  =", format(result))
	}
}
