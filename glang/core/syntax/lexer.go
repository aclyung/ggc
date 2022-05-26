package syntax

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	source
	semi bool

	line, col int // position
	token     token
	lit       string   // valid if tok is _Name, _Literal, or _Semi ("semicolon", "newline", or "EOF"); may be malformed if bad is true
	bad       bool     // valid if tok is _Literal, true if a syntax error occurred, lit may be malformed
	kind      LitKind  // valid if tok is _Literal
	op        Operator // valid if tok is _Operator, _AssignOp, or _IncOp
	prec      int      // valid if tok is _Operator, _AssignOp, or _IncOp
}

const EOFCHAR = '\000'

func TokenizingTest(filename, str string) {
	testPrint("Testing Tokenization")
	println("Input " + filename + " { \n" + str + "\n}")
	r := strings.NewReader(str)
	var p parser
	start := time.Now()
	p.init(r, nil)
	p.next()
	testPrint("Results")

	for p.token != _EOF {
		str := p.token.String()
		switch p.token {
		case _Semi:
			println(str + " ")
			p.next()
			continue
		case _Operator:
			str = p.op.String()
		case _Name, _Literal:
			color := Yellow
			if p.token == _Literal {
				color = Green
			}
			str = color + str
			str += White + "(" + p.lit + ")"
		}
		color := Cyan
		if p.token.isKeyword() {
			color = Purple
		}
		if p.token == _Assign || p.token == _Operator {
			color = Cyan
		}
		print(color + str + " " + Reset)
		p.next()

	}
	testPrint("Tokenizing Test End")
	result(true, time.Since(start).Seconds())
}

func (l *lexer) init(r io.Reader, errh func(line, col int, msg string)) {
	l.source.init(r, errh)
	l.semi = false
}

func (l *lexer) errorf(format string, args ...interface{}) {
	l.error(fmt.Sprintf(format, args...))
}

// errorAtf reports an error at a byte column offset relative to the current token start.
func (l *lexer) errorAtf(offset int, format string, args ...interface{}) {
	l.errh(l.line, l.col+offset, fmt.Sprintf(format, args...))
}

func (l *lexer) ident() {
	for isLetter(l.ch) || isDecimal(l.ch) {
		l.nextch()
	}

	lit := l.segment()
	tok := keyword(lit)
	if tok.isKeyword() {
		l.token = tok
		return
	}
	l.semi = true
	l.lit = lit
	l.token = _Name
}

func (l *lexer) setLit(kind LitKind, ok bool) {
	l.semi = true
	l.token = _Literal
	l.lit = l.segment()
	l.bad = !ok
	l.kind = kind
}

func (l *lexer) next() {
	semi := l.semi
	l.semi = false

	//redo:
	//iLine, iCol := l.pos()

	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' && !semi || l.ch == '\r' {
		l.nextch()
	}
	l.source.start()

	l.line, l.col = l.pos()

	if isLetter(l.ch) || l.ch >= utf8.RuneSelf && l.atIdentChar(true) {
		l.ident()
		return
	}

	switch l.ch {
	case -1:
		if semi {
			l.lit = "EOF"
			l.token = _Semi
			break
		}
		l.token = _EOF

	case '\n':
		l.nextch()
		l.lit = "newline"
		l.token = _Semi

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.number(false)
	case '=':
		l.nextch()
		l.token = _Assign
	case '"':
		l.nextch()
	case '\'':
		l.nextch()
	case '(':
		l.nextch()
		l.token = _Lparen

	case '[':
		l.nextch()
		l.token = _Lbrack

	case '{':
		l.nextch()
		l.token = _Lbrace

	case ',':
		l.nextch()
		l.token = _Comma

	case ';':
		l.nextch()
		l.lit = "semicolon"
		l.token = _Semi

	case ')':
		l.nextch()
		l.semi = true
		l.token = _Rparen

	case ']':
		l.nextch()
		l.semi = true
		l.token = _Rbrack

	case '}':
		l.nextch()
		l.semi = true
		l.token = _Rbrace

	case ':':
		l.nextch()
		if l.ch == '=' {
			l.nextch()
			l.token = _Define
			break
		}
		l.token = _Colon

	case '.':
		l.nextch()
		if isDecimal(l.ch) {
			l.number(true)
			break
		}
		l.token = _Dot

	case '+':
		l.nextch()
		l.op, l.prec = Add, precAdd
		if l.ch != '+' {
			goto assignoper
		}
		l.nextch()
		l.semi = true
		l.token = _IncOp
	case '-':
		l.nextch()
		l.op, l.prec = Sub, precAdd
		if l.ch != '+' {
			goto assignoper
		}
		l.nextch()
		l.semi = true
		l.token = _IncOp
	case '*':
		l.nextch()
		l.op, l.prec = Mul, precMul
		if l.ch != '*' {
			goto assignoper
		}
		l.nextch()
		l.semi = true
		l.token = _IncOp
	case '/':
		l.nextch()
		l.op, l.prec = Div, precMul
		if l.ch != '+' {
			goto assignoper
		}
		l.nextch()
		l.semi = true
		l.token = _IncOp
	case '%':
		l.nextch()
		l.op, l.prec = Rem, precMul
		goto assignoper
	}

	return

assignoper:
	if l.ch == '=' {
		l.nextch()
		l.token = _AssignOp
		return
	}
	l.token = _Operator

}

func (l *lexer) number(afterDot bool) {
	kind := IntLit
	ok := true
	if !afterDot {
		for isDecimal(l.ch) {
			l.nextch()
			if l.ch == '.' {
				l.nextch()
				afterDot = true
				break
			}
		}
	}

	if afterDot {
		kind = FloatLit
		digitExist := false
		for isDecimal(l.ch) {
			l.nextch()
			digitExist = true
		}
		if !digitExist {
			ok = false
			l.errorf("No digit after '.'")
		}
	}
	l.setLit(kind, ok)
}

func (l *lexer) atIdentChar(first bool) bool {
	switch {
	case unicode.IsLetter(l.ch) || l.ch == '_':
		// ok
	case unicode.IsDigit(l.ch):
		if first {
			l.errorf("identifier cannot begin with digit %#U", l.ch)
		}
	case l.ch >= utf8.RuneSelf:
		l.errorf("invalid character %#U in identifier", l.ch)
	default:
		return false
	}
	return true
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isLetter(ch rune) bool  { return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' }
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }

//
//func (next *lexer) ReadNumberToken() {
//	next.token = INT
//	for isDecimal(next.current()) {
//		next.nextch()
//	}
//	if next.current() == '.' {
//		next.nextch()
//		next.token = FLOAT
//		for isDecimal(next.current()) {
//			next.nextch()
//		}
//	}
//	text := next.text[next.start:next.col]
//	if next.token == INT {
//		val, err := strconv.ParseInt(text, 10, 64)
//		if err != nil {
//			fmt.Println(next.text[next.start:next.col], " type not match error[Int64]")
//			//next.Diagnose("ERROR: the Number is not Valid int64.", general.ERROR)
//		}
//		next.value = val
//		return
//	}
//	val, err := strconv.ParseFloat(text, 8)
//	if err != nil {
//		fmt.Println(next.text[next.start:next.col], " type not match error[Float64]")
//		//next.Diagnose("ERROR: the Number is not Valid float64.", general.ERROR)
//	}
//	next.value = val
//}
