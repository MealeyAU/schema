package printer

import "fmt"

type Printer struct {
	SeparatorCharacter string
	SeparatorFactor    int
}

const (
	SeparatorLong   = 4
	SeparatorMedium = 2
	SeparatorShort  = 1
)

// Stringf prints a string with a set of args
func (p *Printer) Stringf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Separator prints a separator
func (p *Printer) Separator(length int) {
	separatorLen := p.SeparatorFactor
	if separatorLen == 0 {
		separatorLen = 4
	}

	separatorChar := p.SeparatorCharacter
	if separatorChar == "" {
		separatorChar = "-"
	}

	str := ""
	for i := 0; i < separatorLen * length; i++ {
		str += separatorChar
	}
	fmt.Println(str)
}

