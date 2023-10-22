package main

import (
	"fmt"
	"strings"
)

type Printer struct {
	indent        string
	builder       *strings.Builder
	atStartOfLine bool
}

func NewPrinter() *Printer {
	return &Printer{"", &strings.Builder{}, true}
}

func (p *Printer) String() string {
	return p.builder.String()
}

func (p *Printer) Print(values ...interface{}) {
	if p.atStartOfLine {
		p.writeIndent()
		p.atStartOfLine = false
	}
	for _, value := range values {
		p.builder.WriteString(fmt.Sprintf("%v", value))
	}
}

func (p *Printer) PrintLn(values ...interface{}) {
	p.Print(values...)
	p.NewLine()
}

func (p *Printer) NewLine() {
	p.Print("\n")
	p.atStartOfLine = true
}

func (p *Printer) writeIndent() {
	p.builder.WriteString(p.indent)
}

const INDENT = "    "

func (p *Printer) Start() {
	p.indent = fmt.Sprintf("%s%s", p.indent, INDENT)
}

func (p *Printer) End() {
	p.indent = p.indent[:len(p.indent)-len(INDENT)]
}
