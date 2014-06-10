package env

import (
	"fmt"
	"github.com/samertm/gosp/parse"
	_ "log"
)

type Env map[string]func([]*parse.Atom) *parse.Atom

type Scope struct {
	Current Env
	Parent *Scope
}

var GlobalScope *Scope

var GlobalKeys Env

func init() {
	GlobalKeys = Env{
		"+": add,
	}
	GlobalScope  = &Scope{Current: GlobalKeys, Parent: nil}
}

var _ = fmt.Print

func add(atoms []*parse.Atom) *parse.Atom {
	acc := 0
	for _, a := range atoms {
		acc += a.Value.(int)
	}
	return &parse.Atom{Value: acc, Type: "int"}
}

func New(s *Scope) *Scope {
	return &Scope{Current: Env{}, Parent: s}
}

func Find(s *Scope, name string) (func([]*parse.Atom) *parse.Atom, bool) {
	if s == nil {
		return nil, false
	} else if f, ok := s.Current[name]; ok {
		return f, true
	}
	return Find(s.Parent, name)
}
