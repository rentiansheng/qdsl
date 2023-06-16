package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

/***************************
    @author: tiansheng.ren
    @date: 5/26/23
    @desc:

***************************/

const eof = 0
const (
	OpTypeNotFoundErr = 1
)

type lex struct {
	input      *bytes.Buffer
	result     []map[string]interface{}
	state      int
	stateStack []int
	err        error
}

const (
	lexStateField        = 0
	lexStateOp           = 1
	lexStateValue        = 2
	LexStateValueElement = 3
	lexStateJoint        = 4
)

func newLex(input []byte) *lex {
	return &lex{
		input: bytes.NewBuffer(input),
	}
}

func (l *lex) Lex(yylval *yySymType) (a int) {
	defer func() {
		//fmt.Println("lex state", l.state, "val", yylval.val, "op", yylval.op, "err")
	}()
	for {
		c, _, err := l.input.ReadRune()
		if err != nil {
			return eof
		}
		switch {
		case unicode.IsSpace(c):
			continue
		case c == '"':
			return l.scanString(yylval)
		case unicode.IsDigit(c) || c == '+' || c == '-':
			l.input.UnreadRune()
			return l.scanNum(yylval)
		case unicode.IsLetter(c):
			l.input.UnreadRune()
			return l.scanLiteral(yylval)
		default:
			if isOp, ok := symbolTypeRela[c]; ok {
				if isOp {
					l.state = lexStateValue
				}
				yylval.op = string(c)
				yylval.val = string(c)
				return int(c)
			}
			l.input.UnreadRune()
			c := l.scanLiteral(yylval)
			return c

		}

	}

}

func (l *lex) lexType(lval *yySymType) int {
	switch l.state {
	case lexStateField:
		l.state = lexStateOp
		return String
	case lexStateOp:
		str := fmt.Sprintf("%v", lval.val)
		if val, ok := opTypeRela[str]; !ok {
			return Literal
		} else {
			lval.op = val
		}
		l.state = lexStateValue
		return OpType
	case lexStateJoint:
		str := fmt.Sprintf("%v", lval.val)
		if jointType, ok := JointTypeRela[str]; !ok {
			lval.errMsg = fmt.Sprintf("joint %v not found", str)
			return eof
		} else {
			lval.val = jointType.String()
		}
		lval.op = str
		l.state = lexStateField
		return JointType
	}
	return Literal
}

func (l *lex) scanString(lval *yySymType) int {
	var buf bytes.Buffer
	for {
		c, _, err := l.input.ReadRune()
		if err != nil {
			break
		}
		if c == '"' {
			lval.val = buf.String()
			return l.lexType(lval)
		}
		buf.WriteRune(c)
	}
	return 0
}

func (l *lex) scanNum(lval *yySymType) int {
	var buf bytes.Buffer
	for {
		c, _, err := l.input.ReadRune()
		if err != nil {
			break
		}
		switch {
		case unicode.IsDigit(c):
			buf.WriteRune(c)
		case strings.IndexRune(".+-eE", c) != -1:
			buf.WriteRune(c)
		default:
			_ = l.input.UnreadRune()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return 0
			}
			lval.val = val
			return Number
		}
	}

	return 0
}

var literal = map[string]interface{}{
	"true":  true,
	"false": false,
	"null":  nil,
}

var opTypeRela = map[string]string{
	">":  ">",
	"<":  "<",
	"!=": "!=",
	">=": ">=",
	"<=": "<=",
	"in": "in",
}

var symbolTypeRela = map[rune]bool{
	'=': true,
	'(': false,
	')': false,
	',': false,
}

type JointTypeRelaType int

const (
	JointTypeRelaTypeAnd = 1
	JointTypeRelaTypeOr  = 2
)

func (j JointTypeRelaType) String() string {
	switch j {
	case JointTypeRelaTypeAnd:
		return "and"
	case JointTypeRelaTypeOr:
		return "or"
	}
	return "unkown"
}

var JointTypeRela = map[string]JointTypeRelaType{
	"and": JointTypeRelaTypeAnd,
	"or":  JointTypeRelaTypeOr,
	"AND": JointTypeRelaTypeAnd,
	"OR":  JointTypeRelaTypeOr,
	"And": JointTypeRelaTypeAnd,
	"Or":  JointTypeRelaTypeOr,
}

func (l *lex) scanRawString(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			break
		}
		switch {
		case unicode.IsLetter(r):
			buf.WriteRune(r)
		default:
			_ = l.input.UnreadRune()
			lval.val = buf.String()
			return l.lexType(lval)

		}
	}
	return 0
}

func (l *lex) scanLiteral(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			break
		}
		switch {
		case unicode.IsLetter(r) || unicode.IsUpper(r) || unicode.IsDigit(r) ||
			r == '_' || r == '.' || r == '>' || r == '<' || r == '!' || r == '=':
			buf.WriteRune(r)
		default:
			_ = l.input.UnreadRune()
			val, ok := literal[buf.String()]
			if ok {
				lval.val = val
				return l.lexType(lval)
			}
			lval.val = buf.String()
			return l.lexType(lval)
		}
	}
	return 0
}

func (l *lex) pushState() {
	l.stateStack = append(l.stateStack, l.state)
	l.state = 0
}

func (l *lex) popState() {
	if len(l.stateStack) == 0 {
		l.state = 0
		return
	}
	l.state = l.stateStack[len(l.stateStack)-1]
	l.stateStack = l.stateStack[:len(l.stateStack)-1]
}

// The parser calls this method on a parse error.
func (l *lex) Error(s string) {
	log.Printf("parse error: %s", s)
}

func main() {
	// goyacc  -o expr.y.go expr.y &&  go run *.go
	sqlArr := []string{
		"a = 1",
		"a = 1 and b = 2",
		"a > 1",
		"a >= 1",
		"a < 1",
		"a <= 1",
		"a != 1",
		`a = b and d in (1,2,3) and c = 2 `,
		`a = b and d in (1,2,3) and c > 2 and d = 3 `,
		`a = b and d in (1,2,3) and c > 2 and d != 3 `,
		`a = b and d in (1,2,3) and c > 2 and d < 3 `,
	}
	for _, sql := range sqlArr {
		l := newLex([]byte(sql + " "))
		yyParse(l)
		fmt.Println("sql: ", sql, "parese result: ", l.result)
	}

}
