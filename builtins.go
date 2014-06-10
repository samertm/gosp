package main

import (
	"github.com/samertm/gosp/env"
	"github.com/samertm/gosp/parse"

	"container/list"
	"errors"
	"strconv"
)

var builtins map[string]func(*list.Element, *env.Scope) (*parse.Atom, error)

func init() {
	builtins = map[string]func(*list.Element, *env.Scope) (*parse.Atom, error){
		"def":    defSetup,
		"lambda": lambdaSetup,
		"if":     ifSetup,
	}
}

func defSetup(e *list.Element, s *env.Scope) (*parse.Atom, error) {
	name := e.Next().Value.(*parse.Atom).Value.(string)
	val, err := eval(e.Next().Next().Value, s)
	if err != nil {
		return nil, err
	}
	return Def(name, val, s), nil
}

func lambdaSetup(e *list.Element, s *env.Scope) (*parse.Atom, error) {
	arglist := e.Next().Value.(*list.List)
	args := make([]string, 0)
	for a := arglist.Front(); a != nil; a = a.Next() {
		args = append(args, a.Value.(*parse.Atom).Value.(string))
	}
	body := make([]interface{}, 0)
	for b := e.Next().Next(); b != nil; b = b.Next() {
		body = append(body, b.Value)
	}
	// taking liberties with the name "Atom"
	return &parse.Atom{
		Value: Lambda(args, body, env.New(s)),
		Type:  "function",
	}, nil
}

func ifSetup(e *list.Element, s *env.Scope) (*parse.Atom, error) {
	test, err := eval(e.Next().Value, s)
	if err != nil {
		return nil, err
	}
	thenStmt := e.Next().Next().Value
	elseStmt := e.Next().Next().Next().Value
	return If(test, thenStmt, elseStmt, s)
}

func Def(name string, val *parse.Atom, s *env.Scope) *parse.Atom {
	switch val.Type {
	case "function":
		s.Current[name] = val.Value.(func([]*parse.Atom) (*parse.Atom, error))
	default:
		s.Current[name] = func([]*parse.Atom) (*parse.Atom, error) { return val, nil }
	}
	return val
}

func Lambda(args []string, body []interface{}, s *env.Scope) func([]*parse.Atom) (*parse.Atom, error) {
	return func(atoms []*parse.Atom) (*parse.Atom, error) {
		if len(args) != len(atoms) {
			expectedLen := strconv.Itoa(len(args))
			recievedLen := strconv.Itoa(len(atoms))
			return nil, errors.New("Mismatched arg lengths: expected " + expectedLen + ", recieved " + recievedLen + " args")
		}
		for i := 0; i < len(args); i++ {
			Def(args[i], atoms[i], s)
		}
		if len(body) == 0 {
			return nil, errors.New("No body for lambda")
		}
		var lastAtom *parse.Atom
		for _, b := range body {
			a, err := eval(b, s)
			if err != nil {
				return nil, err
			}
			// TODO make more efficient
			lastAtom = a
		}
		return lastAtom, nil
	}
}

// then-stmt and els-stmt passed in as interface{} to avoid evaluation
func If(test *parse.Atom, thenStmt, elseStmt interface{}, s *env.Scope) (*parse.Atom, error) {
	if test.Type != "nil" {
		return eval(thenStmt, s)
	}
	return eval(elseStmt, s)
}
