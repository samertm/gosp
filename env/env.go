package env

import (
	"fmt"
	"github.com/samertm/gosp/parse"
)

var Keys map[string]func([]*parse.Atom) *parse.Atom

func init() {
	Keys = map[string]func([]*parse.Atom) *parse.Atom {
		"+": add,
	}
}

var _ = fmt.Print

func add(atoms []*parse.Atom) *parse.Atom {
	acc := 0
	for _, a := range atoms {
		acc += a.Value.(int)
	}
	return &parse.Atom{Value: acc, Type: "int"}
}
