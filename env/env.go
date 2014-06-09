package env

import (
	"fmt"
	"github.com/samertm/gosp/parse"
	_ "log"
)

var Keys map[string]func([]*parse.Atom) *parse.Atom

func init() {
	Keys = map[string]func([]*parse.Atom) *parse.Atom{
		"+":   add,
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

func Def(name string, val *parse.Atom) *parse.Atom {
	Keys[name] = func([]*parse.Atom) *parse.Atom {return val}
	return val
}
