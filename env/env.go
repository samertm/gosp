package env

import (
	"github.com/samertm/gosp/parse"

	"fmt"
	"errors"
)

type Env map[string]func([]*parse.Atom) (*parse.Atom, error)

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

func add(atoms []*parse.Atom) (*parse.Atom, error) {
	acc := 0
	for _, a := range atoms {
		switch a.Type {
		case "int":
			acc += a.Value.(int)
		default:
			return nil, errors.New("add: type mismatch. Got (" + a.Type + "), expected (int)")
		}
	}
	return &parse.Atom{Value: acc, Type: "int"}, nil
}

func New(s *Scope) *Scope {
	return &Scope{Current: Env{}, Parent: s}
}

func Find(s *Scope, name string) (func([]*parse.Atom) (*parse.Atom, error), bool) {
	if s == nil {
		return nil, false
	} else if f, ok := s.Current[name]; ok {
		return f, true
	}
	return Find(s.Parent, name)
}
