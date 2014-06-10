package main

import (
	"github.com/samertm/gosp/env"
	"github.com/samertm/gosp/parse"

	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
	"strings"
	"strconv"
)

func printList(t *list.Element) {
	fmt.Println("(")
	for ; t != nil; t = t.Next() {
		// my first go type switch!
		switch ty := t.Value.(type) {
		case *list.List:
			l := t.Value.(*list.List)
			printList(l.Front())
		case *parse.Atom:
			a := t.Value.(*parse.Atom)
			fmt.Print(a.Value)
		default:
			fmt.Println("error", ty)
		}
	}
	fmt.Println(")")
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

func eval(i interface{}, s *env.Scope) (*parse.Atom, error) {
	switch i.(type) {
	case *list.List:
		e := i.(*list.List).Front()
		t := e.Value.(*parse.Atom)
		if t.Type != "symbol" {
			return nil, errors.New("Expected symbol")
		}
		// built ins
		switch t.Value.(string) {
		case "def":
			name := e.Next().Value.(*parse.Atom).Value.(string)
			val, err := eval(e.Next().Next().Value, s)
			if err != nil {
				return nil, err
			}
			return Def(name, val, s), nil
		case "lambda":
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
		default:
			fun, ok := env.Find(s, t.Value.(string))
			if ok == false {
				return nil, errors.New("Symbol not in function table")
			}
			args := make([]*parse.Atom, 0)
			for e = e.Next(); e != nil; e = e.Next() {
				// eval step
				val, err := eval(e.Value, s)
				if err != nil {
					return nil, err
				}
				args = append(args, val)
			}
			// fun returns two values
			return fun(args)
		}
	case *parse.Atom:
		a := i.(*parse.Atom)
		switch a.Type {
		case "int":
			return a, nil
		case "symbol":
			val, ok := env.Find(s, a.Value.(string))
			if ok == false {
				return nil, errors.New("Symbol '" + a.Value.(string) +"' not found")
			}
			return val([]*parse.Atom{})
		}
	}
	return nil, errors.New("nope")
}

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("gosp> ")
		input, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		input = strings.TrimSpace(input)
		placeholder, err := parse.Parse(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		a, err := eval(placeholder, env.GlobalScope)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(a.Value)
	}
}
