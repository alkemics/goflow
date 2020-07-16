package nodes

import (
	"context"
	"fmt"
)

func Printer(values []interface{}) {
	for _, v := range values {
		fmt.Println(v)
	}
}

// ConditionalPrinter prints values if print is true.
func ConditionalPrinter(print bool, values []interface{}) {
	if print {
		Printer(values)
	}
}

func Print1() { fmt.Println(1) }
func Print2() { fmt.Println(2) }

func PrinterCtx(ctx context.Context, values []interface{}) {
	Printer(append([]interface{}{ctx}, values...))
}
