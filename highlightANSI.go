// Package syntaxhighlight provides syntax highlighting for code. It currently
// uses a language-independent lexer and performs decently on JavaScript, Java,
// Ruby, Python, Go, and C.
package syntaxhighlight

import (
	"bytes"
	"io"
	"text/template"

	"github.com/sourcegraph/annotate"
)

const (
	ResetAll = "\033[0m"

	Bold = "\033[1m"
	Dim = "\033[2m"
	Underlined = "\033[4m"
	Blink = "\033[5m"
	Reverse = "\033[7m"
	Hidden = "\033[8m"

	ResetBold = "\033[21m"
	ResetDim = "\033[22m"
	ResetUnderlined = "\033[24m"
	ResetBlink = "\033[25m"
	ResetReverse = "\033[27m"
	ResetHidden = "\033[28m"

	Default = "\033[39m"
	Black = "\033[30m"
	Red = "\033[31m"
	Green = "\033[32m"
	Yellow = "\033[33m"
	Blue = "\033[34m"
	Magenta = "\033[35m"
	Cyan = "\033[36m"
	LightGray = "\033[37m"
	DarkGray = "\033[90m"
	LightRed = "\033[91m"
	LightGreen = "\033[92m"
	LightYellow = "\033[93m"
	LightBlue = "\033[94m"
	LightMagenta = "\033[95m"
	LightCyan = "\033[96m"
	White = "\033[97m"

	BackgroundDefault = "\033[49m"
	BackgroundBlack = "\033[40m"
	BackgroundRed = "\033[41m"
	BackgroundGreen = "\033[42m"
	BackgroundYellow = "\033[43m"
	BackgroundBlue = "\033[44m"
	BackgroundMagenta = "\033[45m"
	BackgroundCyan = "\033[46m"
	BackgroundLightGray = "\033[47m"
	BackgroundDarkGray = "\033[100m"
	BackgroundLightRed = "\033[101m"
	BackgroundLightGreen = "\033[102m"
	BackgroundLightYellow = "\033[103m"
	BackgroundLightBlue = "\033[104m"
	BackgroundLightMagenta = "\033[105m"
	BackgroundLightCyan = "\033[106m"
	BackgroundWhite = "\033[107m"
)
// AnsiConfig holds the ANSI class configuration to be used by annotators when
// highlighting code.
type AnsiConfig struct {
	String        string
	Keyword       string
	Comment       string
	Type          string
	Literal       string
	Punctuation   string
	Plaintext     string
	Tag           string
	HTMLTag       string
	HTMLAttrName  string
	HTMLAttrValue string
	Decimal       string
	Whitespace    string
}

type AnsiPrinter AnsiConfig

// Class returns the set class for a given token Kind.
func (c AnsiConfig) Class(kind Kind) string {
	switch kind {
	case String:
		return c.String
	case Keyword:
		return c.Keyword
	case Comment:
		return c.Comment
	case Type:
		return c.Type
	case Literal:
		return c.Literal
	case Punctuation:
		return c.Punctuation
	case Plaintext:
		return c.Plaintext
	case Tag:
		return c.Tag
	case HTMLTag:
		return c.HTMLTag
	case HTMLAttrName:
		return c.HTMLAttrName
	case HTMLAttrValue:
		return c.HTMLAttrValue
	case Decimal:
		return c.Decimal
	}
	return ""
}

func (p AnsiPrinter) Print(w io.Writer, kind Kind, tokText string) error {
	class := ((HTMLConfig)(p)).Class(kind)
	if class != "" {

		_, err := io.WriteString(w, class)
		if err != nil {
			return err
		}

	}
	template.HTMLEscape(w, []byte(tokText))
	if class != "" {
		_, err := w.Write([]byte(ResetAll))
		if err != nil {
			return err
		}
	}
	return nil
}

type AnsiAnnotator  AnsiConfig

func (a AnsiAnnotator) Annotate(start int, kind Kind, tokText string) (*annotate.Annotation, error) {
	class := ((HTMLConfig)(a)).Class(kind)
	if class != "" {
		left := []byte(``)
		left = append(left, []byte(class)...)
		left = append(left, []byte(``)...)
		return &annotate.Annotation{
			Start: start, End: start + len(tokText),
			Left: left, Right: []byte(ResetAll),
		}, nil
	}
	return nil, nil
}

// DefaultHTMLConfig's class names match those of google-code-prettify
// (https://code.google.com/p/google-code-prettify/).
var DefaultAnsiConfig = AnsiConfig{
	String:        LightBlue,
	Keyword:       Bold+LightCyan,
	Comment:       Yellow,
	Type:          White,
	Literal:       LightYellow,
	Punctuation:   DarkGray,
	Plaintext:     DarkGray,
	Tag:           LightGray,
	HTMLTag:       LightGray,
	HTMLAttrName:  Blue,
	HTMLAttrValue: Blue,
	Decimal:       LightCyan,
	Whitespace:    "",
}

func AsANSI(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := Print(NewScanner(src), &buf, AnsiPrinter(DefaultAnsiConfig))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
