package parser

import (
	"fmt"
	"strings"
)

var level = 0

func indent(msg string) string {
	return strings.Repeat("\t", level-1) + msg
}

func trace(msg string) string {
	level += 1
	fmt.Println(indent("BEGIN " + msg))
	return msg
}

func un(msg string) {
	fmt.Println(indent("END " + msg))
	level -= 1
}
